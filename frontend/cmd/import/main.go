// Command import reads establecimientos_para_importar.csv and inserts each row
// as an establecimiento record through rcapi (used in-process as a library).
//
// It streams the CSV so the (large) file never needs to fit in memory, and
// fans rows out to a pool of workers that insert in parallel, each pausing a
// configurable amount between inserts to avoid saturating the DB.
//
// Usage:
//
//	CONFIG=../backend/dev.yaml go run ./cmd/import -csv ../establecimientos_para_importar.csv
//
// Flags:
//
//	-csv      path to the CSV file
//	-config   path to the rcapi yaml config (falls back to the CONFIG env var)
//	-workers  number of concurrent insert workers
//	-delay    pause each worker takes between its inserts (e.g. 500us, 5ms)
//	-start    skip this many data rows before inserting (for resuming)
//	-limit    stop after dispatching this many rows (0 = no limit)
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null/v6"

	rcapiconfig "github.com/mklfarha/radarcdmx/backend/rcapi/config"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
	establecimientotypes "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/contacto"
	mainentity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/ubicacion"
)

const (
	// dd/mm/yyyy, e.g. 01/04/2026
	fechaAltaLayout = "02/01/2006"
	// dd/mm/yyyy HH:MM, e.g. 06/06/2026 12:21
	timestampLayout = "02/01/2006 15:04"
)

func main() {
	var (
		csvPath    = flag.String("csv", "./establecimientos_para_importar.csv", "path to the CSV file to import")
		configPath = flag.String("config", os.Getenv("CONFIG"), "path to the rcapi yaml config (defaults to $CONFIG)")
		workers    = flag.Int("workers", 10, "number of concurrent insert workers")
		delay      = flag.Duration("delay", 500*time.Microsecond, "pause each worker takes between its inserts")
		start      = flag.Int64("start", 0, "skip this many data rows before inserting (for resuming)")
		limit      = flag.Int64("limit", 0, "stop after dispatching this many rows (0 = no limit)")
	)
	flag.Parse()

	if *workers < 1 {
		*workers = 1
	}

	if *configPath == "" {
		log.Fatal("no config provided: pass -config or set the CONFIG env var (e.g. ../backend/dev.yaml)")
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	provider, err := rcapiconfig.NewWithPathAndEnvironment(*configPath, os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	rcapiCore, err := core.New(core.Params{Provider: provider})
	if err != nil {
		log.Fatalf("init rcapi core: %v", err)
	}
	defer rcapiCore.Destroy()

	f, err := os.Open(*csvPath)
	if err != nil {
		log.Fatalf("open csv: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.ReuseRecord = true // we copy fields out immediately, so reuse is safe
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("read header: %v", err)
	}
	cols := indexHeader(header)

	ctx := context.Background()
	est := rcapiCore.Establecimiento()

	type job struct {
		row    int64
		record []string
	}

	var (
		row        int64
		dispatched int64
		inserted   int64
		skipped    int64
		failed     int64
		startTime  = time.Now()
	)

	// The establecimiento module is safe for concurrent use (each Insert opens
	// its own transaction on the shared *sql.DB pool), so workers can share it.
	jobs := make(chan job, *workers*2)
	var wg sync.WaitGroup

	for w := 0; w < *workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				e, err := buildEstablecimiento(j.record, cols)
				if err != nil {
					logger.Printf("row %d: parse error: %v", j.row, err)
					atomic.AddInt64(&failed, 1)
					continue
				}

				if _, err := est.Insert(ctx, establecimientotypes.UpsertRequest{Establecimiento: e}); err != nil {
					logger.Printf("row %d (uuid=%s): insert error: %v", j.row, e.UUID, err)
					atomic.AddInt64(&failed, 1)
					continue
				}

				if n := atomic.AddInt64(&inserted, 1); n%1000 == 0 {
					logger.Printf("progress: inserted=%d failed=%d elapsed=%s",
						n, atomic.LoadInt64(&failed), time.Since(startTime).Round(time.Second))
				}

				time.Sleep(*delay)
			}
		}()
	}

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			logger.Printf("row %d: read error: %v", row+1, err)
			atomic.AddInt64(&failed, 1)
			row++
			continue
		}
		row++

		if row <= *start {
			skipped++
			continue
		}

		// csv.Reader reuses the backing array between Read calls, so copy the
		// record before handing it to a worker goroutine.
		cp := make([]string, len(record))
		copy(cp, record)
		jobs <- job{row: row, record: cp}

		dispatched++
		if *limit > 0 && dispatched >= *limit {
			logger.Printf("reached limit of %d dispatched rows", *limit)
			break
		}
	}

	close(jobs)
	wg.Wait()

	logger.Printf("done: inserted=%d skipped=%d failed=%d dispatched=%d total_rows=%d elapsed=%s",
		atomic.LoadInt64(&inserted), skipped, atomic.LoadInt64(&failed), dispatched, row,
		time.Since(startTime).Round(time.Second))
}

// indexHeader maps each CSV column name to its position, stripping a leading
// UTF-8 BOM from the first column if present.
func indexHeader(header []string) map[string]int {
	cols := make(map[string]int, len(header))
	for i, name := range header {
		name = strings.TrimPrefix(name, "\ufeff")
		cols[strings.TrimSpace(name)] = i
	}
	return cols
}

func buildEstablecimiento(record []string, cols map[string]int) (mainentity.Establecimiento, error) {
	get := func(name string) string {
		if idx, ok := cols[name]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	id, err := uuid.FromString(get("uuid"))
	if err != nil {
		return mainentity.Establecimiento{}, fmt.Errorf("invalid uuid %q: %w", get("uuid"), err)
	}

	now := time.Now().UTC()
	createdAt := parseTimestamp(get("created_at"), now)
	updatedAt := parseTimestamp(get("updated_at"), now)

	return mainentity.Establecimiento{
		UUID:            id,
		IdDenue:         parseNullInt(get("id_denue")),
		Clee:            parseNullString(get("clee")),
		Nombre:          parseNullString(get("nombre")),
		RazonSocial:     parseNullString(get("razon_social")),
		PerOcu:          parseNullString(get("per_ocu")),
		CodigoActividad: parseNullInt(get("codigo_actividad")),
		NombreActividad: parseNullString(get("nombre_actividad")),
		UsoDeSuelo:      parseNullString(get("uso_de_suelo")),
		ClaveCatastral:  parseNullString(get("clave_catastral")),
		Contacto:        contacto.ContactoFromJSON(rawJSON(get("contacto"))),
		Ubicacion:       ubicacion.UbicacionFromJSON(rawJSON(get("ubicacion"))),
		FechaAlta:       parseNullDate(get("fecha_alta")),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}

func rawJSON(s string) []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

func parseNullString(s string) null.String {
	if s == "" {
		return null.String{}
	}
	return null.StringFrom(s)
}

func parseNullInt(s string) null.Int64 {
	if s == "" {
		return null.Int64{}
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return null.Int64{}
	}
	return null.IntFrom(v)
}

func parseNullDate(s string) null.Time {
	if s == "" {
		return null.Time{}
	}
	t, err := time.Parse(fechaAltaLayout, s)
	if err != nil {
		return null.Time{}
	}
	return null.TimeFrom(t)
}

func parseTimestamp(s string, fallback time.Time) time.Time {
	if s == "" {
		return fallback
	}
	t, err := time.Parse(timestampLayout, s)
	if err != nil {
		return fallback
	}
	return t
}

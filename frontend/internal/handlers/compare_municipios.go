package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
)

// compareMunicipiosQuery aggregates establecimientos for a single municipio,
// grouping by the 2-digit prefix of the activity code and the per_ocu range so
// the handler can derive totals, employee estimates and per-sector counts.
const compareMunicipiosQuery = `
SELECT
	LEFT(CAST(e.codigo_actividad AS CHAR), 2) AS prefix2,
	e.per_ocu,
	COUNT(*) AS total
FROM establecimiento e
WHERE LOWER(JSON_UNQUOTE(JSON_EXTRACT(e.ubicacion, '$.municipio'))) = LOWER(?)
GROUP BY prefix2, e.per_ocu`

var perOcuNumberRe = regexp.MustCompile(`\d+`)

// ListCompareMunicipiosHandler compares two municipios by establecimiento count,
// approximate employees (from per_ocu ranges) and counts per sector.
type ListCompareMunicipiosHandler struct {
	logger *log.Logger
	core   *core.Implementation
	// prefixToSector maps a 2-digit SCIAN prefix to its sector catalog entry.
	prefixToSector map[string]sectorItem
	// sectorOrder lists sector keys in catalog order for stable output.
	sectorCatalog []sectorItem
}

// NewListCompareMunicipiosHandler constructs the handler and builds the sector
// lookup from the embedded sectores.json catalog.
func NewListCompareMunicipiosHandler(logger *log.Logger, rcapiCore *core.Implementation) (*ListCompareMunicipiosHandler, error) {
	catalog := make([]sectorItem, 0)
	if err := json.Unmarshal(sectoresJSON, &catalog); err != nil {
		return nil, fmt.Errorf("invalid sectores.json: %w", err)
	}

	prefixToSector := make(map[string]sectorItem)
	for _, item := range catalog {
		for _, code := range item.CodigoSector {
			prefixToSector[strings.TrimSpace(code)] = item
		}
	}

	return &ListCompareMunicipiosHandler{
		logger:         logger,
		core:           rcapiCore,
		prefixToSector: prefixToSector,
		sectorCatalog:  catalog,
	}, nil
}

func (h *ListCompareMunicipiosHandler) Method() string  { return http.MethodGet }
func (h *ListCompareMunicipiosHandler) Pattern() string { return "/api/municipios/compare" }

type sectorCount struct {
	CodigoSector []string `json:"codigo_sector"`
	Sector       string   `json:"sector"`
	Total        int64    `json:"total"`
}

type municipioComparison struct {
	Municipio             string        `json:"municipio"`
	TotalEstablecimientos int64         `json:"total_establecimientos"`
	EmpleadosAproximados  int64         `json:"empleados_aproximados"`
	Sectores              []sectorCount `json:"sectores"`
}

type compareMunicipiosResponse struct {
	Municipios []municipioComparison `json:"municipios"`
}

func (h *ListCompareMunicipiosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	municipioA := strings.TrimSpace(query.Get("a"))
	municipioB := strings.TrimSpace(query.Get("b"))

	if municipioA == "" || municipioB == "" {
		http.Error(w, "both 'a' and 'b' municipio query params are required", http.StatusBadRequest)
		return
	}

	comparisonA, err := h.compareOne(r, municipioA)
	if err != nil {
		h.logger.Printf("compare municipio %q: %v", municipioA, err)
		http.Error(w, "failed to compare municipios", http.StatusInternalServerError)
		return
	}

	comparisonB, err := h.compareOne(r, municipioB)
	if err != nil {
		h.logger.Printf("compare municipio %q: %v", municipioB, err)
		http.Error(w, "failed to compare municipios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(compareMunicipiosResponse{
		Municipios: []municipioComparison{comparisonA, comparisonB},
	}); err != nil {
		h.logger.Printf("encode compare response: %v", err)
	}
}

// compareOne aggregates the metrics for a single municipio.
func (h *ListCompareMunicipiosHandler) compareOne(r *http.Request, municipio string) (municipioComparison, error) {
	rows, err := h.core.DB().QueryContext(r.Context(), compareMunicipiosQuery, municipio)
	if err != nil {
		return municipioComparison{}, err
	}
	defer rows.Close()

	var totalEstablecimientos int64
	var empleadosAproximados int64
	// sectorTotals is keyed by the sector name; "" is the unclassified bucket.
	sectorTotals := make(map[string]int64)

	for rows.Next() {
		var prefix sql.NullString
		var perOcu sql.NullString
		var total int64
		if err := rows.Scan(&prefix, &perOcu, &total); err != nil {
			return municipioComparison{}, err
		}

		totalEstablecimientos += total
		empleadosAproximados += upperPerOcu(perOcu.String) * total

		sectorKey := ""
		if prefix.Valid {
			if item, ok := h.prefixToSector[strings.TrimSpace(prefix.String)]; ok {
				sectorKey = item.Sector
			}
		}
		sectorTotals[sectorKey] += total
	}

	if err := rows.Err(); err != nil {
		return municipioComparison{}, err
	}

	return municipioComparison{
		Municipio:             municipio,
		TotalEstablecimientos: totalEstablecimientos,
		EmpleadosAproximados:  empleadosAproximados,
		Sectores:              h.buildSectorCounts(sectorTotals),
	}, nil
}

// buildSectorCounts produces a per-sector list ordered by the catalog, omitting
// sectors with no establecimientos, and appending an "Otros" bucket if needed.
func (h *ListCompareMunicipiosHandler) buildSectorCounts(sectorTotals map[string]int64) []sectorCount {
	result := make([]sectorCount, 0)

	for _, item := range h.sectorCatalog {
		total := sectorTotals[item.Sector]
		if total == 0 {
			continue
		}
		result = append(result, sectorCount{
			CodigoSector: item.CodigoSector,
			Sector:       item.Sector,
			Total:        total,
		})
	}

	// Keep a stable order by the first SCIAN code of each sector.
	sort.SliceStable(result, func(i, j int) bool {
		return firstCode(result[i].CodigoSector) < firstCode(result[j].CodigoSector)
	})

	if otros := sectorTotals[""]; otros > 0 {
		result = append(result, sectorCount{
			CodigoSector: []string{},
			Sector:       "Otros / No clasificado",
			Total:        otros,
		})
	}

	return result
}

func firstCode(codes []string) int {
	if len(codes) == 0 {
		return 1 << 30
	}
	n, err := strconv.Atoi(strings.TrimSpace(codes[0]))
	if err != nil {
		return 1 << 30
	}
	return n
}

// upperPerOcu extracts the upper bound of a per_ocu range string such as
// "11 a 30 personas" (-> 30) or "251 y más personas" (-> 251). It returns the
// largest number found, or 0 when the string carries no digits.
func upperPerOcu(perOcu string) int64 {
	matches := perOcuNumberRe.FindAllString(perOcu, -1)
	var upper int64
	for _, m := range matches {
		n, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			continue
		}
		if n > upper {
			upper = n
		}
	}
	return upper
}

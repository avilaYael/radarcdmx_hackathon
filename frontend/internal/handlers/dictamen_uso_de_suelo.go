package handlers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
	establecimientotypes "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
)

//go:embed sectores.json
var sectoresJSONDictamen []byte

type sectorDictamen struct {
	CodigoSector   []string `json:"codigo_sector"`
	Sector         string   `json:"sector"`
	UsosPermitidos []string `json:"usos_permitidos"`
}

type DictamenUsoDeSueloResponse struct {
	UUID            string   `json:"uuid"`
	Nombre          string   `json:"nombre"`
	UsoDeSuelo      string   `json:"uso_de_suelo"`
	CodigoActividad int64    `json:"codigo_actividad"`
	CodigoSector    string   `json:"codigo_sector"`
	Sector          string   `json:"sector"`
	UsosPermitidos  []string `json:"usos_permitidos"`
	Aprobado        bool     `json:"aprobado"`
	Razon           string   `json:"razon"`
}

// DictamenUsoDeSueloHandler evaluates whether an establecimiento's uso_de_suelo
// is permitted for its economic sector (derived from the first 2 digits of codigo_actividad).
type DictamenUsoDeSueloHandler struct {
	logger   *log.Logger
	core     *core.Implementation
	sectores []sectorDictamen
}

// NewDictamenUsoDeSueloHandler constructs DictamenUsoDeSueloHandler.
func NewDictamenUsoDeSueloHandler(logger *log.Logger, rcapiCore *core.Implementation) (*DictamenUsoDeSueloHandler, error) {
	var sectores []sectorDictamen
	if err := json.Unmarshal(sectoresJSONDictamen, &sectores); err != nil {
		return nil, fmt.Errorf("invalid sectores.json: %w", err)
	}
	return &DictamenUsoDeSueloHandler{logger: logger, core: rcapiCore, sectores: sectores}, nil
}

func (h *DictamenUsoDeSueloHandler) Method() string  { return http.MethodGet }
func (h *DictamenUsoDeSueloHandler) Pattern() string { return "/api/dictamen-uso-de-suelo" }

func (h *DictamenUsoDeSueloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawUUID := strings.TrimSpace(r.URL.Query().Get("uuid"))
	if rawUUID == "" {
		http.Error(w, "uuid is required", http.StatusBadRequest)
		return
	}

	parsedUUID, err := uuid.FromString(rawUUID)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	res, err := h.core.Establecimiento().FetchEstablecimientoByUuid(r.Context(), establecimientotypes.FetchEstablecimientoByUuidRequest{
		UUID: parsedUUID,
	})
	if err != nil {
		h.logger.Printf("fetch establecimiento by uuid: %v", err)
		http.Error(w, "failed to fetch establecimiento", http.StatusInternalServerError)
		return
	}

	if len(res.Results) == 0 {
		http.Error(w, "establecimiento not found", http.StatusNotFound)
		return
	}

	est := res.Results[0]

	if !est.UsoDeSuelo.Valid || strings.TrimSpace(est.UsoDeSuelo.String) == "" {
		h.writeJSON(w, DictamenUsoDeSueloResponse{
			UUID:     est.UUID.String(),
			Nombre:   est.Nombre.String,
			Aprobado: false,
			Razon:    "El establecimiento no tiene uso de suelo registrado",
		})
		return
	}

	if !est.CodigoActividad.Valid {
		h.writeJSON(w, DictamenUsoDeSueloResponse{
			UUID:       est.UUID.String(),
			Nombre:     est.Nombre.String,
			UsoDeSuelo: est.UsoDeSuelo.String,
			Aprobado:   false,
			Razon:      "El establecimiento no tiene código de actividad registrado",
		})
		return
	}

	codigoActividad := est.CodigoActividad.Int64
	codigoSector := fmt.Sprintf("%d", codigoActividad)
	if len(codigoSector) >= 2 {
		codigoSector = codigoSector[:2]
	}

	sector := h.findSector(codigoSector)
	if sector == nil {
		h.writeJSON(w, DictamenUsoDeSueloResponse{
			UUID:            est.UUID.String(),
			Nombre:          est.Nombre.String,
			UsoDeSuelo:      est.UsoDeSuelo.String,
			CodigoActividad: codigoActividad,
			CodigoSector:    codigoSector,
			Aprobado:        false,
			Razon:           fmt.Sprintf("No se encontró un sector para el código '%s'", codigoSector),
		})
		return
	}

	usoDeSuelo := strings.TrimSpace(est.UsoDeSuelo.String)
	aprobado := false
	for _, permitido := range sector.UsosPermitidos {
		if strings.EqualFold(usoDeSuelo, strings.TrimSpace(permitido)) {
			aprobado = true
			break
		}
	}

	razon := ""
	if aprobado {
		razon = fmt.Sprintf("El uso de suelo '%s' es compatible con el sector '%s'", usoDeSuelo, sector.Sector)
	} else {
		razon = fmt.Sprintf("El uso de suelo '%s' no está permitido para el sector '%s'", usoDeSuelo, sector.Sector)
	}

	h.writeJSON(w, DictamenUsoDeSueloResponse{
		UUID:            est.UUID.String(),
		Nombre:          est.Nombre.String,
		UsoDeSuelo:      usoDeSuelo,
		CodigoActividad: codigoActividad,
		CodigoSector:    codigoSector,
		Sector:          sector.Sector,
		UsosPermitidos:  sector.UsosPermitidos,
		Aprobado:        aprobado,
		Razon:           razon,
	})
}

func (h *DictamenUsoDeSueloHandler) findSector(codigoSector string) *sectorDictamen {
	for i := range h.sectores {
		for _, code := range h.sectores[i].CodigoSector {
			if code == codigoSector {
				return &h.sectores[i]
			}
		}
	}
	return nil
}

func (h *DictamenUsoDeSueloHandler) writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Printf("encode dictamen response: %v", err)
	}
}

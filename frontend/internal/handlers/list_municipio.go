package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
)

const listMunicipiosQuery = `
SELECT DISTINCT
	JSON_UNQUOTE(JSON_EXTRACT(e.ubicacion, '$.municipio')) AS municipio
FROM establecimiento e
WHERE JSON_EXTRACT(e.ubicacion, '$.municipio') IS NOT NULL
	AND TRIM(JSON_UNQUOTE(JSON_EXTRACT(e.ubicacion, '$.municipio'))) <> ''
ORDER BY municipio ASC`

// ListMunicipioHandler lists the distinct municipios used by establecimientos.
type ListMunicipioHandler struct {
	logger *log.Logger
	core   *core.Implementation
}

// NewListMunicipioHandler constructs ListMunicipioHandler.
func NewListMunicipioHandler(logger *log.Logger, rcapiCore *core.Implementation) *ListMunicipioHandler {
	return &ListMunicipioHandler{logger: logger, core: rcapiCore}
}

func (h *ListMunicipioHandler) Method() string  { return http.MethodGet }
func (h *ListMunicipioHandler) Pattern() string { return "/api/municipios" }

type listMunicipioResponse struct {
	Items []string `json:"items"`
}

func (h *ListMunicipioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := h.core.DB().QueryContext(r.Context(), listMunicipiosQuery)
	if err != nil {
		h.logger.Printf("list municipios query: %v", err)
		http.Error(w, "failed to list municipios", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := make([]string, 0)
	for rows.Next() {
		var municipio sql.NullString
		if err := rows.Scan(&municipio); err != nil {
			h.logger.Printf("scan municipios row: %v", err)
			http.Error(w, "failed to decode municipios", http.StatusInternalServerError)
			return
		}
		if !municipio.Valid || municipio.String == "" {
			continue
		}
		items = append(items, municipio.String)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("list municipios rows: %v", err)
		http.Error(w, "failed to list municipios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listMunicipioResponse{Items: items}); err != nil {
		h.logger.Printf("encode municipios response: %v", err)
	}
}

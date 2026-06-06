package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
)

const listUsoDeSueloQuery = `
SELECT DISTINCT
	e.uso_de_suelo
FROM establecimiento e
WHERE e.uso_de_suelo IS NOT NULL
	AND TRIM(e.uso_de_suelo) <> ''
ORDER BY e.uso_de_suelo ASC`

// ListUsoDeSueloHandler lists the distinct uso de suelo values used by establecimientos.
type ListUsoDeSueloHandler struct {
	logger *log.Logger
	core   *core.Implementation
}

// NewListUsoDeSueloHandler constructs ListUsoDeSueloHandler.
func NewListUsoDeSueloHandler(logger *log.Logger, rcapiCore *core.Implementation) *ListUsoDeSueloHandler {
	return &ListUsoDeSueloHandler{logger: logger, core: rcapiCore}
}

func (h *ListUsoDeSueloHandler) Method() string  { return http.MethodGet }
func (h *ListUsoDeSueloHandler) Pattern() string { return "/api/usos-de-suelo" }

type listUsoDeSueloResponse struct {
	Items []string `json:"items"`
}

func (h *ListUsoDeSueloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := h.core.DB().QueryContext(r.Context(), listUsoDeSueloQuery)
	if err != nil {
		h.logger.Printf("list usos de suelo query: %v", err)
		http.Error(w, "failed to list usos de suelo", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := make([]string, 0)
	for rows.Next() {
		var usoDeSuelo sql.NullString
		if err := rows.Scan(&usoDeSuelo); err != nil {
			h.logger.Printf("scan usos de suelo row: %v", err)
			http.Error(w, "failed to decode usos de suelo", http.StatusInternalServerError)
			return
		}
		if !usoDeSuelo.Valid || usoDeSuelo.String == "" {
			continue
		}
		items = append(items, usoDeSuelo.String)
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("list usos de suelo rows: %v", err)
		http.Error(w, "failed to list usos de suelo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listUsoDeSueloResponse{Items: items}); err != nil {
		h.logger.Printf("encode usos de suelo response: %v", err)
	}
}

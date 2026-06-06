package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
)

const listActividadesQueryBase = `
SELECT DISTINCT
	e.codigo_actividad,
	e.nombre_actividad
FROM establecimiento e
WHERE e.codigo_actividad IS NOT NULL
	AND e.nombre_actividad IS NOT NULL`

var twoDigitsRe = regexp.MustCompile(`^[0-9]{2}$`)

// ListActividadHandler lists unique activity codes and names.
type ListActividadHandler struct {
	logger *log.Logger
	core   *core.Implementation
}

// NewListActividadHandler constructs ListActividadHandler.
func NewListActividadHandler(logger *log.Logger, rcapiCore *core.Implementation) *ListActividadHandler {
	return &ListActividadHandler{logger: logger, core: rcapiCore}
}

func (h *ListActividadHandler) Method() string  { return http.MethodGet }
func (h *ListActividadHandler) Pattern() string { return "/api/actividades" }

type actividadItem struct {
	Codigo int64  `json:"codigo_actividad"`
	Nombre string `json:"nombre_actividad"`
}

type listActividadResponse struct {
	Items []actividadItem `json:"items"`
}

func (h *ListActividadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	prefixes, err := parsePrefixes(query["prefix2"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := int64(500)
	offset := int64(0)

	if raw := query.Get("limit"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || parsed <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	if raw := query.Get("offset"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || parsed < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		offset = parsed
	}

	whereClauses := []string{}
	args := make([]any, 0)

	if len(prefixes) > 0 {
		placeholders := strings.Repeat("?,", len(prefixes))
		placeholders = strings.TrimSuffix(placeholders, ",")
		whereClauses = append(whereClauses, "LEFT(CAST(e.codigo_actividad AS CHAR), 2) IN ("+placeholders+")")
		for _, prefix := range prefixes {
			args = append(args, prefix)
		}
	}

	finalQuery := listActividadesQueryBase
	if len(whereClauses) > 0 {
		finalQuery += "\n\tAND " + strings.Join(whereClauses, "\n\tAND ")
	}
	finalQuery += "\nORDER BY e.codigo_actividad ASC, e.nombre_actividad ASC\nLIMIT ?, ?"
	args = append(args, offset, limit)

	rows, err := h.core.DB().QueryContext(r.Context(), finalQuery, args...)
	if err != nil {
		h.logger.Printf("list actividades query: %v", err)
		http.Error(w, "failed to list actividades", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := make([]actividadItem, 0)
	for rows.Next() {
		var codigo sql.NullInt64
		var nombre sql.NullString
		if err := rows.Scan(&codigo, &nombre); err != nil {
			h.logger.Printf("scan actividades row: %v", err)
			http.Error(w, "failed to decode actividades", http.StatusInternalServerError)
			return
		}
		if !codigo.Valid || !nombre.Valid {
			continue
		}
		items = append(items, actividadItem{
			Codigo: codigo.Int64,
			Nombre: nombre.String,
		})
	}

	if err := rows.Err(); err != nil {
		h.logger.Printf("list actividades rows: %v", err)
		http.Error(w, "failed to list actividades", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listActividadResponse{Items: items}); err != nil {
		h.logger.Printf("encode actividades response: %v", err)
	}
}

func parsePrefixes(rawValues []string) ([]string, error) {
	if len(rawValues) == 0 {
		return nil, nil
	}

	unique := make(map[string]struct{})
	res := make([]string, 0)
	for _, raw := range rawValues {
		for _, token := range strings.Split(raw, ",") {
			prefix := strings.TrimSpace(token)
			if prefix == "" {
				continue
			}
			if !twoDigitsRe.MatchString(prefix) {
				return nil, errors.New("invalid prefix2, expected an array of 2-digit values")
			}
			if _, found := unique[prefix]; found {
				continue
			}
			unique[prefix] = struct{}{}
			res = append(res, prefix)
		}
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}

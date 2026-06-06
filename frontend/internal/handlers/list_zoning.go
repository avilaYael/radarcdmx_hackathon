package handlers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//go:embed zoning.json
var zoningJSON []byte

// ListZoningHandler serves the SEDUVI land-use zoning layer (GeoJSON
// FeatureCollection) from zoning.json.
type ListZoningHandler struct {
	logger *log.Logger
	body   []byte
}

// NewListZoningHandler constructs ListZoningHandler and validates JSON at startup.
func NewListZoningHandler(logger *log.Logger) (*ListZoningHandler, error) {
	if !json.Valid(zoningJSON) {
		return nil, fmt.Errorf("invalid zoning.json")
	}

	return &ListZoningHandler{logger: logger, body: zoningJSON}, nil
}

func (h *ListZoningHandler) Method() string  { return http.MethodGet }
func (h *ListZoningHandler) Pattern() string { return "/api/zoning" }

func (h *ListZoningHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(h.body); err != nil {
		h.logger.Printf("write zoning response: %v", err)
	}
}

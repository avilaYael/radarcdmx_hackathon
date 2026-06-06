package handlers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//go:embed mercados.json
var mercadosJSON []byte

// ListMercadosHandler serves the public markets layer (GeoJSON
// FeatureCollection) from mercados.json.
type ListMercadosHandler struct {
	logger *log.Logger
	body   []byte
}

// NewListMercadosHandler constructs ListMercadosHandler and validates JSON at startup.
func NewListMercadosHandler(logger *log.Logger) (*ListMercadosHandler, error) {
	if !json.Valid(mercadosJSON) {
		return nil, fmt.Errorf("invalid mercados.json")
	}

	return &ListMercadosHandler{logger: logger, body: mercadosJSON}, nil
}

func (h *ListMercadosHandler) Method() string  { return http.MethodGet }
func (h *ListMercadosHandler) Pattern() string { return "/api/mercados" }

func (h *ListMercadosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(h.body); err != nil {
		h.logger.Printf("write mercados response: %v", err)
	}
}

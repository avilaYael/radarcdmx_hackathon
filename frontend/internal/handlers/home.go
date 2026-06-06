package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// HomeHandler serves the landing page.
type HomeHandler struct {
	logger *log.Logger
}

// NewHomeHandler constructs a HomeHandler.
func NewHomeHandler(logger *log.Logger) *HomeHandler {
	return &HomeHandler{logger: logger}
}

func (h *HomeHandler) Method() string  { return http.MethodGet }
func (h *HomeHandler) Pattern() string { return "/" }

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join("RadarMX-main", "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		h.logger.Printf("home index file not found: %v", err)
		http.Error(w, "frontend asset not found", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, indexPath)
}

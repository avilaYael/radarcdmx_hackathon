package handlers

import (
	"log"
	"net/http"
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("radarcdmx — home"))
}

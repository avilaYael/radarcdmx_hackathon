package handlers

import (
	"net/http"
)

// HealthHandler reports service liveness.
type HealthHandler struct{}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Method() string  { return http.MethodGet }
func (h *HealthHandler) Pattern() string { return "/healthz" }

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

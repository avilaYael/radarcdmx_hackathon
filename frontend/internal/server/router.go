package server

import (
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"

	"radarcdmx-web/internal/route"
)

// RouterParams collects everything needed to build the router, including all
// routes registered into the "routes" group across the fx graph.
type RouterParams struct {
	fx.In

	Routes []route.Route `group:"routes"`
	Logger *log.Logger
}

// NewRouter wires up middleware and mounts every registered route.
func NewRouter(p RouterParams) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	for _, rt := range p.Routes {
		p.Logger.Printf("route: %-6s %s", rt.Method(), rt.Pattern())
		r.Method(rt.Method(), rt.Pattern(), rt)
	}

	return r
}

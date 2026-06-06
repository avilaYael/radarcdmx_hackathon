package server

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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

	uiFS := os.DirFS("RadarMX-main")
	uiFileServer := http.FileServer(http.FS(uiFS))

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			http.NotFound(w, req)
			return
		}

		if strings.HasPrefix(req.URL.Path, "/api/") || req.URL.Path == "/healthz" {
			http.NotFound(w, req)
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean(req.URL.Path), "/")
		if cleanPath == "." || cleanPath == "" {
			cleanPath = "index.html"
		}

		if _, err := fs.Stat(uiFS, cleanPath); err == nil {
			clone := req.Clone(req.Context())
			clone.URL.Path = "/" + cleanPath
			uiFileServer.ServeHTTP(w, clone)
			return
		}

		// SPA fallback for client-side routes that are not real files.
		if !strings.Contains(cleanPath, ".") {
			clone := req.Clone(req.Context())
			clone.URL.Path = "/index.html"
			uiFileServer.ServeHTTP(w, clone)
			return
		}

		http.NotFound(w, req)
	})

	return r
}

package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"radarcdmx-web/internal/config"
)

// NewLogger provides the application logger.
func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}

// NewHTTPServer constructs the HTTP server and binds its lifecycle to fx, so it
// starts when the app starts and shuts down gracefully when the app stops.
func NewHTTPServer(lc fx.Lifecycle, router *chi.Mux, cfg *config.Config, logger *log.Logger) *http.Server {
	srv := &http.Server{Addr: cfg.Addr, Handler: router}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Printf("radarcdmx listening on %s", srv.Addr)
			go func() {
				if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
					logger.Printf("http server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Println("shutting down http server")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

// Module wires the server, logger, and router into the fx graph. The invoke
// forces the server to be instantiated even though nothing depends on it.
var Module = fx.Module("server",
	fx.Provide(
		NewLogger,
		NewRouter,
		NewHTTPServer,
	),
	fx.Invoke(func(*http.Server) {}),
)

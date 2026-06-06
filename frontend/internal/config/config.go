package config

import (
	"os"

	"go.uber.org/fx"
)

// Config holds runtime configuration sourced from the environment.
type Config struct {
	// Addr is the TCP address the HTTP server listens on.
	Addr string
	// DevMode enables development conveniences (verbose logging, on-disk assets).
	DevMode bool
}

// New builds a Config from environment variables, applying sane defaults.
func New() *Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":4242"
	}
	return &Config{
		Addr:    addr,
		DevMode: os.Getenv("DEV") == "1",
	}
}

// Module exposes configuration to the fx graph.
var Module = fx.Module("config", fx.Provide(New))

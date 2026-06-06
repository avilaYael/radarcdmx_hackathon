package main

import (
	"go.uber.org/fx"

	"radarcdmx-web/internal/config"
	"radarcdmx-web/internal/handlers"
	"radarcdmx-web/internal/rcapi"
	"radarcdmx-web/internal/server"
)

func main() {
	fx.New(
		config.Module,
		rcapi.Module,
		server.Module,
		handlers.Module,
	).Run()
}

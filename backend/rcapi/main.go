package main

import (
	rcapiconfig "github.com/mklfarha/radarcdmx/backend/rcapi/config"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"

	pbserver "github.com/mklfarha/radarcdmx/backend/rcapi/idl/server"
)

func main() {
	fx.New(
		fx.Provide(
			zap.NewProduction,
			rcapiconfig.New,
			core.New,
		),
		fx.Invoke(httpServer),

		fx.Invoke(pbserver.New),
	).Run()
}

func httpServer(config config.Provider, logger *zap.Logger) {
	// http port from config
	httpPort := config.Get("ports.http").String()

	go http.ListenAndServe(":"+httpPort, nil)

	logger.Info(`Serving HTTP on PORT: %s`, zap.String("port", httpPort))
}

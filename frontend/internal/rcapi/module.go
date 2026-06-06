package rcapi

import (
	rcapiconfig "github.com/mklfarha/radarcdmx/backend/rcapi/config"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
	"go.uber.org/fx"
)

// Module exposes rcapi core as an in-process dependency for frontend handlers.
var Module = fx.Module("rcapi",
	fx.Provide(
		rcapiconfig.New,
		core.New,
	),
)

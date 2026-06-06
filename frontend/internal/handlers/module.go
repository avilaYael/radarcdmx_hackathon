package handlers

import (
	"go.uber.org/fx"

	"radarcdmx-web/internal/route"
)

// Module registers all HTTP handlers into the "routes" group so the router can
// mount them automatically. Add new handlers here.
var Module = fx.Module("handlers",
	fx.Provide(
		route.As(NewHomeHandler),
		route.As(NewHealthHandler),
		route.As(NewListSectoresHandler),
		route.As(NewListActividadHandler),
		route.As(NewListEstablecimientoHandler),
		route.As(NewListEstablecimientoNearbyHandler),
	),
)

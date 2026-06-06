package route

import (
	"net/http"

	"go.uber.org/fx"
)

// Route is an HTTP handler that knows where it should be mounted. Handlers
// implement this interface so the router can register them generically.
type Route interface {
	http.Handler
	// Method is the HTTP verb (e.g. http.MethodGet).
	Method() string
	// Pattern is the chi URL pattern (e.g. "/case/{slug}").
	Pattern() string
}

// As annotates a route constructor so its result is collected into the
// "routes" group consumed by the router.
func As(constructor any) any {
	return fx.Annotate(
		constructor,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

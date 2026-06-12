module radarcdmx-web

go 1.26.4

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/guregu/null/v6 v6.0.0
	github.com/mklfarha/radarcdmx/backend/rcapi v0.0.0
	go.uber.org/fx v1.24.0
)

replace github.com/mklfarha/radarcdmx/backend/rcapi => ../backend/rcapi

require (
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	go.einride.tech/aip v0.86.3 // indirect
	go.uber.org/config v1.4.1 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

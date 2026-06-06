package establecimiento

import (
	"context"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/repository"
	coretypes "github.com/mklfarha/radarcdmx/backend/rcapi/core/types"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type Module interface {
	FetchEstablecimientoByUuid(ctx context.Context, req types.FetchEstablecimientoByUuidRequest, opts ...Option) (types.FetchEstablecimientoByUuidResponse, error)

	List(ctx context.Context, req types.ListRequest, opts ...Option) (types.ListResponse, error)

	Upsert(ctx context.Context, req types.UpsertRequest, opts ...Option) (types.UpsertResponse, error)
	Insert(ctx context.Context, req types.UpsertRequest, opts ...Option) (types.UpsertResponse, error)
	Update(ctx context.Context, req types.UpsertRequest, opts ...Option) (types.UpsertResponse, error)
}

type module struct {
	mu         sync.Mutex
	sg         singleflight.Group
	cache      *gocache.Cache
	repository *repository.Implementation
}

func New(params coretypes.ModuleParams) Module {
	return &module{
		repository: params.Repository,
		cache:      gocache.New(30*time.Second, 5*time.Minute),
	}
}

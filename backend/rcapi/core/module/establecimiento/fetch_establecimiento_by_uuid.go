package establecimiento

import (
	"context"
	"fmt"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
)

func (m *module) FetchEstablecimientoByUuid(
	ctx context.Context,
	req types.FetchEstablecimientoByUuidRequest,
	opts ...Option,
) (types.FetchEstablecimientoByUuidResponse, error) {

	resolvedOpts := applyAllOptions(opts)
	cacheKey := fmt.Sprintf("FetchEstablecimientoByUuid:%v", req)
	if !resolvedOpts.SkipCache {
		if cached, found := m.cache.Get(cacheKey); found {
			return cached.(types.FetchEstablecimientoByUuidResponse), nil
		}
	}
	v, err, _ := m.sg.Do(cacheKey, func() (any, error) {
		models, err := m.repository.Queries.FetchEstablecimientoByUuid(
			ctx,
			req.UUID.String(),
		)
		if err != nil {

			return types.FetchEstablecimientoByUuidResponse{}, err
		}
		return types.FetchEstablecimientoByUuidResponse{
			Results: mapModelsToEntities(models),
		}, nil
	})
	if err != nil {
		return types.FetchEstablecimientoByUuidResponse{}, err
	}
	result := v.(types.FetchEstablecimientoByUuidResponse)
	if !resolvedOpts.SkipCache {
		m.cache.Set(cacheKey, result, 0)
	}
	return result, nil

}

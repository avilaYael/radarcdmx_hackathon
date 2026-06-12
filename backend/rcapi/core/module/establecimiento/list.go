package establecimiento

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	repogen "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"

	"go.uber.org/zap"
	"slices"
)

func (m *module) List(ctx context.Context,
	request types.ListRequest,
	opts ...Option) (types.ListResponse, error) {

	reqPlusOne := request
	reqPlusOne.PageSize = request.PageSize + 1

	query, err := m.repository.BuildListEntityQuery(
		ctx,
		reqPlusOne,
		main_entity.Establecimiento{},
		false)
	if err != nil {

		return types.ListResponse{}, err
	}

	resolvedOpts := applyAllOptions(opts)
	cacheKey := fmt.Sprintf("ListEstablecimiento:%v", request)
	if !resolvedOpts.SkipCache {
		if cached, found := m.cache.Get(cacheKey); found {
			return cached.(types.ListResponse), nil
		}
	}

	v, listErr, _ := m.sg.Do(cacheKey, func() (any, error) {
		var rows *sql.Rows
		rows, err = m.repository.DB.QueryContext(ctx, query)
		if err != nil {

			m.logger.Error("error in executing query for ListEstablecimiento", zap.String("query", query), zap.Error(err))
			return types.ListResponse{}, err
		}
		defer rows.Close()
		var scanGetters []func(*repogen.Establecimiento) any
		if len(request.GetIncludeFields()) > 0 {
			for _, f := range listFields {
				if slices.Contains(request.GetIncludeFields(), f) {
					scanGetters = append(scanGetters, listFieldRegistry[f])
				}
			}
		} else if len(request.GetExcludeFields()) > 0 {
			for _, f := range listFields {
				if !slices.Contains(request.GetExcludeFields(), f) {
					scanGetters = append(scanGetters, listFieldRegistry[f])
				}
			}
		} else {
			for _, f := range listFields {
				scanGetters = append(scanGetters, listFieldRegistry[f])
			}
		}

		var items []repogen.Establecimiento
		for rows.Next() {
			var i repogen.Establecimiento
			fields := make([]any, 0, len(scanGetters))
			for _, getter := range scanGetters {
				fields = append(fields, getter(&i))
			}
			if err := rows.Scan(fields...); err != nil {

				return types.ListResponse{}, err
			}
			items = append(items, i)
		}
		if err := rows.Close(); err != nil {

			return types.ListResponse{}, err
		}
		if err := rows.Err(); err != nil {

			return types.ListResponse{}, err
		}

		hasNextPage := false
		if len(items) > int(request.PageSize) {
			hasNextPage = true
			items = items[:request.PageSize]
		}

		return types.ListResponse{
			Establecimiento: mapModelsToEntities(items),
			HasNextPage:     hasNextPage,
		}, nil
	})
	if listErr != nil {
		return types.ListResponse{}, listErr
	}
	listResponse := v.(types.ListResponse)
	if !resolvedOpts.SkipCache {
		m.cache.Set(cacheKey, listResponse, 0)
	}
	return listResponse, nil
}

var listFields = []string{

	"uuid",

	"id_denue",

	"clee",

	"nombre",

	"razon_social",

	"per_ocu",

	"codigo_actividad",

	"nombre_actividad",

	"uso_de_suelo",

	"clave_catastral",

	"contacto",

	"ubicacion",

	"fecha_alta",

	"created_at",

	"updated_at",
}

var listFieldRegistry = map[string]func(*repogen.Establecimiento) any{

	"uuid": func(i *repogen.Establecimiento) any { return &i.UUID },

	"id_denue": func(i *repogen.Establecimiento) any { return &i.IdDenue },

	"clee": func(i *repogen.Establecimiento) any { return &i.Clee },

	"nombre": func(i *repogen.Establecimiento) any { return &i.Nombre },

	"razon_social": func(i *repogen.Establecimiento) any { return &i.RazonSocial },

	"per_ocu": func(i *repogen.Establecimiento) any { return &i.PerOcu },

	"codigo_actividad": func(i *repogen.Establecimiento) any { return &i.CodigoActividad },

	"nombre_actividad": func(i *repogen.Establecimiento) any { return &i.NombreActividad },

	"uso_de_suelo": func(i *repogen.Establecimiento) any { return &i.UsoDeSuelo },

	"clave_catastral": func(i *repogen.Establecimiento) any { return &i.ClaveCatastral },

	"contacto": func(i *repogen.Establecimiento) any { return &i.Contacto },

	"ubicacion": func(i *repogen.Establecimiento) any { return &i.Ubicacion },

	"fecha_alta": func(i *repogen.Establecimiento) any { return &i.FechaAlta },

	"created_at": func(i *repogen.Establecimiento) any { return &i.CreatedAt },

	"updated_at": func(i *repogen.Establecimiento) any { return &i.UpdatedAt },
}

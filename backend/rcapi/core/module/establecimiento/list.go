package establecimiento

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	repogen "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"

	"slices"
)

func (m *module) List(ctx context.Context,
	request types.ListRequest,
	opts ...Option) (types.ListResponse, error) {
	query, err := m.repository.BuildListEntityQuery(
		ctx,
		request,
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
		if len(request.GetOrderBy().Fields) > 0 {
			// pin to a single connection so SET sort_buffer_size applies to the same session
			conn, connErr := m.repository.DB.Conn(ctx)
			if connErr != nil {
				return types.ListResponse{}, connErr
			}
			defer conn.Close()

			// increase sort buffer size
			// TODO make this configurable
			bufferRows, bufferErr := conn.QueryContext(ctx, "SET sort_buffer_size=5000000")
			if bufferErr != nil {

				return types.ListResponse{}, bufferErr
			}
			bufferRows.Close()
			rows, err = conn.QueryContext(ctx, query)
			if err != nil {

				fmt.Printf("error in executing query for ListEstablecimiento: %v \n %v\n", query, err)
				return types.ListResponse{}, err
			}
		} else {
			rows, err = m.repository.DB.QueryContext(ctx, query)
			if err != nil {

				fmt.Printf("error in executing query for ListEstablecimiento: %v \n %v\n", query, err)
				return types.ListResponse{}, err
			}
		}
		defer rows.Close()
		var items []repogen.Establecimiento
		for rows.Next() {
			var i repogen.Establecimiento
			fields := []any{}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "uuid") {
					fields = append(fields, &i.UUID)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "uuid") {
					fields = append(fields, &i.UUID)
				}
			} else {
				fields = append(fields, &i.UUID)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "id_denue") {
					fields = append(fields, &i.IdDenue)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "id_denue") {
					fields = append(fields, &i.IdDenue)
				}
			} else {
				fields = append(fields, &i.IdDenue)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "clee") {
					fields = append(fields, &i.Clee)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "clee") {
					fields = append(fields, &i.Clee)
				}
			} else {
				fields = append(fields, &i.Clee)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "nombre") {
					fields = append(fields, &i.Nombre)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "nombre") {
					fields = append(fields, &i.Nombre)
				}
			} else {
				fields = append(fields, &i.Nombre)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "razon_social") {
					fields = append(fields, &i.RazonSocial)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "razon_social") {
					fields = append(fields, &i.RazonSocial)
				}
			} else {
				fields = append(fields, &i.RazonSocial)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "per_ocu") {
					fields = append(fields, &i.PerOcu)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "per_ocu") {
					fields = append(fields, &i.PerOcu)
				}
			} else {
				fields = append(fields, &i.PerOcu)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "codigo_actividad") {
					fields = append(fields, &i.CodigoActividad)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "codigo_actividad") {
					fields = append(fields, &i.CodigoActividad)
				}
			} else {
				fields = append(fields, &i.CodigoActividad)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "nombre_actividad") {
					fields = append(fields, &i.NombreActividad)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "nombre_actividad") {
					fields = append(fields, &i.NombreActividad)
				}
			} else {
				fields = append(fields, &i.NombreActividad)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "uso_de_suelo") {
					fields = append(fields, &i.UsoDeSuelo)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "uso_de_suelo") {
					fields = append(fields, &i.UsoDeSuelo)
				}
			} else {
				fields = append(fields, &i.UsoDeSuelo)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "clave_catastral") {
					fields = append(fields, &i.ClaveCatastral)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "clave_catastral") {
					fields = append(fields, &i.ClaveCatastral)
				}
			} else {
				fields = append(fields, &i.ClaveCatastral)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "contacto") {
					fields = append(fields, &i.Contacto)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "contacto") {
					fields = append(fields, &i.Contacto)
				}
			} else {
				fields = append(fields, &i.Contacto)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "ubicacion") {
					fields = append(fields, &i.Ubicacion)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "ubicacion") {
					fields = append(fields, &i.Ubicacion)
				}
			} else {
				fields = append(fields, &i.Ubicacion)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "fecha_alta") {
					fields = append(fields, &i.FechaAlta)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "fecha_alta") {
					fields = append(fields, &i.FechaAlta)
				}
			} else {
				fields = append(fields, &i.FechaAlta)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "created_at") {
					fields = append(fields, &i.CreatedAt)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "created_at") {
					fields = append(fields, &i.CreatedAt)
				}
			} else {
				fields = append(fields, &i.CreatedAt)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "updated_at") {
					fields = append(fields, &i.UpdatedAt)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "updated_at") {
					fields = append(fields, &i.UpdatedAt)
				}
			} else {
				fields = append(fields, &i.UpdatedAt)
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

		hasNextPage, err := m.listHasNextPage(ctx, request, opts...)
		if err != nil {

			return types.ListResponse{}, err
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

func (m *module) ListCount(ctx context.Context,
	request types.ListRequest,
	opts ...Option) (int64, error) {
	query, err := m.repository.BuildListEntityQuery(ctx, request, main_entity.Establecimiento{}, true)
	if err != nil {
		return -1, err
	}

	//fmt.Printf("query count: %s \n", query)

	var count int64
	err = m.repository.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (m *module) listHasNextPage(ctx context.Context,
	request types.ListRequest,
	opts ...Option) (bool, error) {
	count, err := m.ListCount(ctx, request, opts...)
	if err != nil {
		return false, err
	}

	if request.GetOffset()+int64(request.GetPageSize()) < count {
		return true, nil
	}
	return false, nil
}

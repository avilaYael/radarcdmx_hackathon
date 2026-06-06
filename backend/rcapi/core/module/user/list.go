package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user/types"
	repogen "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/user"

	"slices"
)

func (m *module) List(ctx context.Context,
	request types.ListRequest,
	opts ...Option) (types.ListResponse, error) {
	query, err := m.repository.BuildListEntityQuery(
		ctx,
		request,
		main_entity.User{},
		false)
	if err != nil {

		return types.ListResponse{}, err
	}

	resolvedOpts := applyAllOptions(opts)
	cacheKey := fmt.Sprintf("ListUser:%v", request)
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

				fmt.Printf("error in executing query for ListUser: %v \n %v\n", query, err)
				return types.ListResponse{}, err
			}
		} else {
			rows, err = m.repository.DB.QueryContext(ctx, query)
			if err != nil {

				fmt.Printf("error in executing query for ListUser: %v \n %v\n", query, err)
				return types.ListResponse{}, err
			}
		}
		defer rows.Close()
		var items []repogen.User
		for rows.Next() {
			var i repogen.User
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
				if slices.Contains(request.GetIncludeFields(), "name") {
					fields = append(fields, &i.Name)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "name") {
					fields = append(fields, &i.Name)
				}
			} else {
				fields = append(fields, &i.Name)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "lastname") {
					fields = append(fields, &i.Lastname)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "lastname") {
					fields = append(fields, &i.Lastname)
				}
			} else {
				fields = append(fields, &i.Lastname)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "email") {
					fields = append(fields, &i.Email)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "email") {
					fields = append(fields, &i.Email)
				}
			} else {
				fields = append(fields, &i.Email)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "password") {
					fields = append(fields, &i.Password)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "password") {
					fields = append(fields, &i.Password)
				}
			} else {
				fields = append(fields, &i.Password)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "status") {
					fields = append(fields, &i.Status)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "status") {
					fields = append(fields, &i.Status)
				}
			} else {
				fields = append(fields, &i.Status)
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

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "created_by") {
					fields = append(fields, &i.CreatedBy)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "created_by") {
					fields = append(fields, &i.CreatedBy)
				}
			} else {
				fields = append(fields, &i.CreatedBy)
			}

			if len(request.GetIncludeFields()) > 0 {
				if slices.Contains(request.GetIncludeFields(), "updated_by") {
					fields = append(fields, &i.UpdatedBy)
				}
			} else if len(request.GetExcludeFields()) > 0 {
				if !slices.Contains(request.GetExcludeFields(), "updated_by") {
					fields = append(fields, &i.UpdatedBy)
				}
			} else {
				fields = append(fields, &i.UpdatedBy)
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
			User:        mapModelsToEntities(items),
			HasNextPage: hasNextPage,
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
	query, err := m.repository.BuildListEntityQuery(ctx, request, main_entity.User{}, true)
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

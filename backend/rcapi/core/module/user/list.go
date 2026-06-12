package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user/types"
	repogen "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/user"

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
		rows, err = m.repository.DB.QueryContext(ctx, query)
		if err != nil {

			m.logger.Error("error in executing query for ListUser", zap.String("query", query), zap.Error(err))
			return types.ListResponse{}, err
		}
		defer rows.Close()
		var scanGetters []func(*repogen.User) any
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

		var items []repogen.User
		for rows.Next() {
			var i repogen.User
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

var listFields = []string{

	"uuid",

	"name",

	"lastname",

	"email",

	"password",

	"status",

	"updated_at",

	"created_by",

	"updated_by",

	"created_at",
}

var listFieldRegistry = map[string]func(*repogen.User) any{

	"uuid": func(i *repogen.User) any { return &i.UUID },

	"name": func(i *repogen.User) any { return &i.Name },

	"lastname": func(i *repogen.User) any { return &i.Lastname },

	"email": func(i *repogen.User) any { return &i.Email },

	"password": func(i *repogen.User) any { return &i.Password },

	"status": func(i *repogen.User) any { return &i.Status },

	"updated_at": func(i *repogen.User) any { return &i.UpdatedAt },

	"created_by": func(i *repogen.User) any { return &i.CreatedBy },

	"updated_by": func(i *repogen.User) any { return &i.UpdatedBy },

	"created_at": func(i *repogen.User) any { return &i.CreatedAt },
}

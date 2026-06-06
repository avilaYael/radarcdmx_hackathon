package repository

import (
	"context"

	"database/sql"

	rcapidb "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/list"
)

type Implementation struct {
	Queries *rcapidb.Queries
	DB      *sql.DB
	List    *list.Implementation
}

func New(db *sql.DB) *Implementation {
	queries := rcapidb.New(db)
	return &Implementation{
		Queries: queries,
		DB:      db,
		List:    list.New(),
	}
}

func (i *Implementation) BuildListEntityQuery(ctx context.Context, request list.ListRequest, entity list.ListEntity, onlyCount bool) (string, error) {
	return i.List.BuildListEntityQuery(ctx, request, entity, onlyCount)
}

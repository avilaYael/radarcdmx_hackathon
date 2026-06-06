package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	rcapiconfig "github.com/mklfarha/radarcdmx/backend/rcapi/config"
	coretypes "github.com/mklfarha/radarcdmx/backend/rcapi/core/types"
	"go.uber.org/config"
	"go.uber.org/fx"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/repository"
)

type Implementation struct {
	db         *sql.DB
	repository *repository.Implementation

	establecimiento establecimiento.Module

	user user.Module
}

type Params struct {
	fx.In
	Provider  config.Provider
	Lifecycle fx.Lifecycle
}

func New(params Params) (*Implementation, error) {

	var dbs rcapiconfig.DBs
	if err := params.Provider.Get("db").Populate(&dbs); err != nil {
		return nil, err
	}

	if len(dbs) == 0 {
		return nil, errors.New("db configuration not found")
	}

	dbconfig := dbs[0]
	db, err := sql.Open(dbconfig.Driver, dbconfig.Path())
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(0)
	repository := repository.New(db)

	if params.Lifecycle != nil {
		params.Lifecycle.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				db.Close()
				return nil
			},
		})
	}

	return &Implementation{
		db:         db,
		repository: repository,
	}, nil
}

func (i *Implementation) Destroy() {
	i.db.Close()
}

func (i *Implementation) DB() *sql.DB {
	return i.db
}

func (i Implementation) Establecimiento() establecimiento.Module {
	if i.establecimiento == nil {
		i.establecimiento = establecimiento.New(coretypes.ModuleParams{
			Repository: i.repository,
		})
	}
	return i.establecimiento
}

func (i Implementation) User() user.Module {
	if i.user == nil {
		i.user = user.New(coretypes.ModuleParams{
			Repository: i.repository,
		})
	}
	return i.user
}

package user

import (
	"context"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user/types"
	rcapidb "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
)

func (m *module) Insert(
	ctx context.Context,
	req types.UpsertRequest,
	opts ...Option,
) (types.UpsertResponse, error) {

	optConfig := applyAllOptions(opts)

	tx := optConfig.SQLTx
	createdTx := false
	if tx == nil {
		ntx, err := m.repository.DB.Begin()
		if err != nil {
			return types.UpsertResponse{}, err
		}
		tx = ntx
		defer tx.Rollback()
		createdTx = true
	}

	qtx := m.repository.Queries.WithTx(tx)
	params := mapUpsertRequestToInsertParams(req)

	_, err := qtx.InsertUser(
		ctx,
		params,
	)
	if err != nil {

		return types.UpsertResponse{}, err
	}

	if createdTx {
		err := tx.Commit()
		if err != nil {

			return types.UpsertResponse{}, err
		}
	}

	return buildInsertResponse(req), nil
}

func buildInsertResponse(req types.UpsertRequest) types.UpsertResponse {
	return types.UpsertResponse{

		UUID: req.User.UUID,
	}
}

func mapUpsertRequestToInsertParams(req types.UpsertRequest) rcapidb.InsertUserParams {
	return rcapidb.InsertUserParams{
		UUID: req.User.UUID.String(),

		Name: req.User.Name,

		Lastname: req.User.Lastname,

		Email: req.User.Email,

		Password: req.User.Password,

		Status: req.User.Status.ToInt64(),

		UpdatedAt: req.User.UpdatedAt,

		CreatedBy: req.User.CreatedBy.String(),

		UpdatedBy: req.User.UpdatedBy.String(),

		CreatedAt: req.User.CreatedAt,
	}
}

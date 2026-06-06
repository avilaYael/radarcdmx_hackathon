package user

import (
	"context"
	"errors"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user/types"
	rcapidb "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
)

func (m *module) Update(
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
	existing, err := qtx.FetchUserByUuidForUpdate(ctx,
		req.User.UUID.String(),
	)
	if err != nil {

		return types.UpsertResponse{}, err
	}

	if len(existing) == 0 {
		err := errors.New("entity not found")

		return types.UpsertResponse{}, err
	}

	params := mapUpsertRequestToUpdateParams(req)
	err = qtx.UpdateUser(
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

	return buildUpdateResponse(req), nil
}

func buildUpdateResponse(req types.UpsertRequest) types.UpsertResponse {
	return types.UpsertResponse{

		UUID: req.User.UUID,
	}
}

func mapUpsertRequestToUpdateParams(req types.UpsertRequest) rcapidb.UpdateUserParams {
	return rcapidb.UpdateUserParams{
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

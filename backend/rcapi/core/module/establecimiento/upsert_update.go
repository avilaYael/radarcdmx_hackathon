package establecimiento

import (
	"context"
	"errors"

	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	rcapidb "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/contacto"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/ubicacion"
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
	existing, err := qtx.FetchEstablecimientoByUuidForUpdate(ctx,
		req.Establecimiento.UUID.String(),
	)
	if err != nil {

		return types.UpsertResponse{}, err
	}

	if len(existing) == 0 {
		err := errors.New("entity not found")

		return types.UpsertResponse{}, err
	}

	params := mapUpsertRequestToUpdateParams(req)
	err = qtx.UpdateEstablecimiento(
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

		UUID: req.Establecimiento.UUID,
	}
}

func mapUpsertRequestToUpdateParams(req types.UpsertRequest) rcapidb.UpdateEstablecimientoParams {
	return rcapidb.UpdateEstablecimientoParams{
		UUID: req.Establecimiento.UUID.String(),

		IdDenue: req.Establecimiento.IdDenue,

		Clee: req.Establecimiento.Clee,

		Nombre: req.Establecimiento.Nombre,

		RazonSocial: req.Establecimiento.RazonSocial,

		PerOcu: req.Establecimiento.PerOcu,

		CodigoActividad: req.Establecimiento.CodigoActividad,

		NombreActividad: req.Establecimiento.NombreActividad,

		UsoDeSuelo: req.Establecimiento.UsoDeSuelo,

		ClaveCatastral: req.Establecimiento.ClaveCatastral,

		Contacto: contacto.ContactoToJSON(req.Establecimiento.Contacto),

		Ubicacion: ubicacion.UbicacionToJSON(req.Establecimiento.Ubicacion),

		FechaAlta: req.Establecimiento.FechaAlta,

		CreatedAt: req.Establecimiento.CreatedAt,

		UpdatedAt: req.Establecimiento.UpdatedAt,
	}
}

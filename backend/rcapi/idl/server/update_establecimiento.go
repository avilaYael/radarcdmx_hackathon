package server

import (
	"context"
	"errors"
	establecimientomodule "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"
	pbmapper "github.com/mklfarha/radarcdmx/backend/rcapi/idl/mapper"

	"go.einride.tech/aip/fieldmask"
	"strings"
)

func (s *server) UpdateEstablecimiento(ctx context.Context, req *pb.UpdateEstablecimientoRequest) (*pb.Establecimiento, error) {

	if req.Establecimiento.Uuid == "" {
		return nil, errors.New("please provide a valid UUID to update")
	}

	err := fieldmask.Validate(req.UpdateMask, req.GetEstablecimiento())
	if err != nil {

		return nil, err
	}

	isFull := fieldmask.IsFullReplacement(req.UpdateMask)

	if !isFull && req.UpdateMask != nil {

		if !strings.Contains(req.UpdateMask.String(), "uuid") {
			req.UpdateMask.Append(req.GetEstablecimiento(), "uuid")
		}

		pkEntity := pbmapper.EstablecimientoFromProto(req.GetEstablecimiento())
		existingRes, err := s.core.Establecimiento().FetchEstablecimientoByUuid(ctx,
			types.FetchEstablecimientoByUuidRequest{
				UUID: pkEntity.UUID,
			},
			establecimientomodule.WithSkipCache(),
		)
		if err != nil {

			return nil, err
		}
		if len(existingRes.Results) == 0 {
			return nil, errors.New("entity not found")
		}

		merged := pbmapper.EstablecimientoToProto(existingRes.Results[0])
		fieldmask.Update(req.UpdateMask, merged, req.GetEstablecimiento())
		req = &pb.UpdateEstablecimientoRequest{Establecimiento: merged, UpdateMask: req.UpdateMask}
	}

	res, err := s.core.Establecimiento().Update(ctx, types.UpsertRequest{
		Establecimiento: pbmapper.EstablecimientoFromProto(req.GetEstablecimiento()),
	})
	if err != nil {

		return nil, err
	}

	fetchRes, err := s.core.Establecimiento().FetchEstablecimientoByUuid(ctx, types.FetchEstablecimientoByUuidRequest(res), establecimientomodule.WithSkipCache())
	if err != nil {

		return nil, err
	}

	if len(fetchRes.Results) == 0 {
		err := errors.New("error fetching entity")

		return nil, err
	}

	return pbmapper.EstablecimientoToProto(fetchRes.Results[0]), nil
}

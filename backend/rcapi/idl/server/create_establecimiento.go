package server

import (
	"context"
	"errors"
	establecimientomodule "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"
	pbmapper "github.com/mklfarha/radarcdmx/backend/rcapi/idl/mapper"
)

func (s *server) CreateEstablecimiento(ctx context.Context, req *pb.CreateEstablecimientoRequest) (*pb.Establecimiento, error) {
	res, err := s.core.Establecimiento().Insert(ctx, types.UpsertRequest{
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

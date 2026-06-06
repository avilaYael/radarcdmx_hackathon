package server

import (
	"context"
	"errors"
	usermodule "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user/types"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"
	pbmapper "github.com/mklfarha/radarcdmx/backend/rcapi/idl/mapper"
)

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	res, err := s.core.User().Insert(ctx, types.UpsertRequest{
		User: pbmapper.UserFromProto(req.GetUser()),
	})
	if err != nil {

		return nil, err
	}

	fetchRes, err := s.core.User().FetchUserByUuid(ctx, types.FetchUserByUuidRequest(res), usermodule.WithSkipCache())
	if err != nil {

		return nil, err
	}

	if len(fetchRes.Results) == 0 {
		err := errors.New("error fetching entity")

		return nil, err
	}

	return pbmapper.UserToProto(fetchRes.Results[0]), nil
}

package mapper

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/user"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"

	"github.com/guregu/null/v6"

	"github.com/mklfarha/radarcdmx/backend/rcapi/enum"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(e main_entity.User) *pb.User {
	return &pb.User{
		Uuid:      e.UUID.String(),
		Name:      e.Name.ValueOrZero(),
		Lastname:  e.Lastname.ValueOrZero(),
		Email:     e.Email,
		Password:  e.Password,
		Status:    pb.UserStatus(e.Status),
		UpdatedAt: timestamppb.New(e.UpdatedAt),
		CreatedBy: e.CreatedBy.String(),
		UpdatedBy: e.UpdatedBy.String(),
		CreatedAt: timestamppb.New(e.CreatedAt),
	}
}

func UserSliceToProto(es []main_entity.User) []*pb.User {
	res := []*pb.User{}
	for _, e := range es {
		res = append(res, UserToProto(e))
	}
	return res
}

func UserFromProto(m *pb.User) main_entity.User {
	if m == nil {
		return main_entity.User{}
	}
	return main_entity.User{
		UUID:      StringToUUID(m.GetUuid()),
		Name:      null.StringFrom(m.Name),
		Lastname:  null.StringFrom(m.Lastname),
		Email:     m.GetEmail(),
		Password:  m.GetPassword(),
		Status:    enum.UserStatus(m.GetStatus()),
		UpdatedAt: m.GetUpdatedAt().AsTime(),
		CreatedBy: StringToUUID(m.GetCreatedBy()),
		UpdatedBy: StringToUUID(m.GetUpdatedBy()),
		CreatedAt: m.GetCreatedAt().AsTime(),
	}
}

func UserSliceFromProto(es []*pb.User) []main_entity.User {
	if es == nil {
		return []main_entity.User{}
	}
	res := []main_entity.User{}
	for _, e := range es {
		res = append(res, UserFromProto(e))
	}
	return res
}

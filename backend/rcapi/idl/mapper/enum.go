package mapper

import (
	"github.com/mklfarha/radarcdmx/backend/rcapi/enum"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"
)

func EstablecimientoStatusSliceToProto(s []enum.EstablecimientoStatus) []pb.EstablecimientoStatus {
	res := []pb.EstablecimientoStatus{}
	for _, e := range s {
		res = append(res, pb.EstablecimientoStatus(e))
	}
	return res
}
func EstablecimientoStatusSliceFromProto(s []pb.EstablecimientoStatus) []enum.EstablecimientoStatus {
	res := []enum.EstablecimientoStatus{}
	for _, e := range s {
		res = append(res, enum.EstablecimientoStatus(e))
	}
	return res
}

func UserStatusSliceToProto(s []enum.UserStatus) []pb.UserStatus {
	res := []pb.UserStatus{}
	for _, e := range s {
		res = append(res, pb.UserStatus(e))
	}
	return res
}
func UserStatusSliceFromProto(s []pb.UserStatus) []enum.UserStatus {
	res := []enum.UserStatus{}
	for _, e := range s {
		res = append(res, enum.UserStatus(e))
	}
	return res
}

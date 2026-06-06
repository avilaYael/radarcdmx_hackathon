package mapper

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"

	"github.com/guregu/null/v6"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func EstablecimientoToProto(e main_entity.Establecimiento) *pb.Establecimiento {
	return &pb.Establecimiento{
		Uuid:            e.UUID.String(),
		IdDenue:         e.IdDenue.ValueOrZero(),
		Clee:            e.Clee.ValueOrZero(),
		Nombre:          e.Nombre.ValueOrZero(),
		RazonSocial:     e.RazonSocial.ValueOrZero(),
		PerOcu:          e.PerOcu.ValueOrZero(),
		CodigoActividad: e.CodigoActividad.ValueOrZero(),
		NombreActividad: e.NombreActividad.ValueOrZero(),
		UsoDeSuelo:      e.UsoDeSuelo.ValueOrZero(),
		ClaveCatastral:  e.ClaveCatastral.ValueOrZero(),
		Contacto:        ContactoToProto(e.Contacto),
		Ubicacion:       UbicacionToProto(e.Ubicacion),
		FechaAlta:       timestamppb.New(e.FechaAlta.ValueOrZero()),
		CreatedAt:       timestamppb.New(e.CreatedAt),
		UpdatedAt:       timestamppb.New(e.UpdatedAt),
	}
}

func EstablecimientoSliceToProto(es []main_entity.Establecimiento) []*pb.Establecimiento {
	res := []*pb.Establecimiento{}
	for _, e := range es {
		res = append(res, EstablecimientoToProto(e))
	}
	return res
}

func EstablecimientoFromProto(m *pb.Establecimiento) main_entity.Establecimiento {
	if m == nil {
		return main_entity.Establecimiento{}
	}
	return main_entity.Establecimiento{
		UUID:            StringToUUID(m.GetUuid()),
		IdDenue:         null.IntFrom(m.GetIdDenue()),
		Clee:            null.StringFrom(m.Clee),
		Nombre:          null.StringFrom(m.Nombre),
		RazonSocial:     null.StringFrom(m.RazonSocial),
		PerOcu:          null.StringFrom(m.PerOcu),
		CodigoActividad: null.IntFrom(m.GetCodigoActividad()),
		NombreActividad: null.StringFrom(m.NombreActividad),
		UsoDeSuelo:      null.StringFrom(m.UsoDeSuelo),
		ClaveCatastral:  null.StringFrom(m.ClaveCatastral),
		Contacto:        ContactoFromProto(m.GetContacto()),
		Ubicacion:       UbicacionFromProto(m.GetUbicacion()),
		FechaAlta:       null.TimeFrom(m.GetFechaAlta().AsTime()),
		CreatedAt:       m.GetCreatedAt().AsTime(),
		UpdatedAt:       m.GetUpdatedAt().AsTime(),
	}
}

func EstablecimientoSliceFromProto(es []*pb.Establecimiento) []main_entity.Establecimiento {
	if es == nil {
		return []main_entity.Establecimiento{}
	}
	res := []main_entity.Establecimiento{}
	for _, e := range es {
		res = append(res, EstablecimientoFromProto(e))
	}
	return res
}

package mapper

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/contacto"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"

	"github.com/guregu/null/v6"
)

func ContactoToProto(e main_entity.Contacto) *pb.Contacto {
	return &pb.Contacto{
		Telefono: e.Telefono.ValueOrZero(),
		Correo:   e.Correo.ValueOrZero(),
		SitioWeb: e.SitioWeb.ValueOrZero(),
	}
}

func ContactoSliceToProto(es []main_entity.Contacto) []*pb.Contacto {
	res := []*pb.Contacto{}
	for _, e := range es {
		res = append(res, ContactoToProto(e))
	}
	return res
}

func ContactoFromProto(m *pb.Contacto) main_entity.Contacto {
	if m == nil {
		return main_entity.Contacto{}
	}
	return main_entity.Contacto{
		Telefono: null.StringFrom(m.Telefono),
		Correo:   null.StringFrom(m.Correo),
		SitioWeb: null.StringFrom(m.SitioWeb),
	}
}

func ContactoSliceFromProto(es []*pb.Contacto) []main_entity.Contacto {
	if es == nil {
		return []main_entity.Contacto{}
	}
	res := []main_entity.Contacto{}
	for _, e := range es {
		res = append(res, ContactoFromProto(e))
	}
	return res
}

package mapper

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/ubicacion"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"

	"github.com/guregu/null/v6"
)

func UbicacionToProto(e main_entity.Ubicacion) *pb.Ubicacion {
	return &pb.Ubicacion{
		Entidad:      e.Entidad.ValueOrZero(),
		Municipio:    e.Municipio.ValueOrZero(),
		Localidad:    e.Localidad.ValueOrZero(),
		Manzana:      e.Manzana.ValueOrZero(),
		CodigoPostal: e.CodigoPostal.ValueOrZero(),
		Calle:        e.Calle.ValueOrZero(),
		NumExt:       e.NumExt.ValueOrZero(),
		NumInt:       e.NumInt.ValueOrZero(),
		Latitud:      e.Latitud.ValueOrZero(),
		Longitud:     e.Longitud.ValueOrZero(),
	}
}

func UbicacionSliceToProto(es []main_entity.Ubicacion) []*pb.Ubicacion {
	res := []*pb.Ubicacion{}
	for _, e := range es {
		res = append(res, UbicacionToProto(e))
	}
	return res
}

func UbicacionFromProto(m *pb.Ubicacion) main_entity.Ubicacion {
	if m == nil {
		return main_entity.Ubicacion{}
	}
	return main_entity.Ubicacion{
		Entidad:      null.StringFrom(m.Entidad),
		Municipio:    null.StringFrom(m.Municipio),
		Localidad:    null.StringFrom(m.Localidad),
		Manzana:      null.IntFrom(m.GetManzana()),
		CodigoPostal: null.StringFrom(m.CodigoPostal),
		Calle:        null.StringFrom(m.Calle),
		NumExt:       null.StringFrom(m.NumExt),
		NumInt:       null.StringFrom(m.NumInt),
		Latitud:      null.FloatFrom(m.GetLatitud()),
		Longitud:     null.FloatFrom(m.GetLongitud()),
	}
}

func UbicacionSliceFromProto(es []*pb.Ubicacion) []main_entity.Ubicacion {
	if es == nil {
		return []main_entity.Ubicacion{}
	}
	res := []main_entity.Ubicacion{}
	for _, e := range es {
		res = append(res, UbicacionFromProto(e))
	}
	return res
}

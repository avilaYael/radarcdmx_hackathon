package establecimiento

import (
	rcapidb "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/mapper"

	"github.com/guregu/null/v6"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/contacto"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/ubicacion"
)

func mapModelsToEntities(models []rcapidb.Establecimiento) []main_entity.Establecimiento {
	result := []main_entity.Establecimiento{}
	for _, p := range models {
		result = append(result, mapModelToEntity(p))
	}
	return result
}

func mapModelToEntity(m rcapidb.Establecimiento) main_entity.Establecimiento {
	return main_entity.Establecimiento{
		UUID:            mapper.StringToUUID(m.UUID),
		IdDenue:         null.NewInt(m.IdDenue.Int64, m.IdDenue.Valid),
		Clee:            null.NewString(m.Clee.String, m.Clee.Valid),
		Nombre:          null.NewString(m.Nombre.String, m.Nombre.Valid),
		RazonSocial:     null.NewString(m.RazonSocial.String, m.RazonSocial.Valid),
		PerOcu:          null.NewString(m.PerOcu.String, m.PerOcu.Valid),
		CodigoActividad: null.NewInt(m.CodigoActividad.Int64, m.CodigoActividad.Valid),
		NombreActividad: null.NewString(m.NombreActividad.String, m.NombreActividad.Valid),
		UsoDeSuelo:      null.NewString(m.UsoDeSuelo.String, m.UsoDeSuelo.Valid),
		ClaveCatastral:  null.NewString(m.ClaveCatastral.String, m.ClaveCatastral.Valid),
		Contacto:        contacto.ContactoFromJSON(m.Contacto),
		Ubicacion:       ubicacion.UbicacionFromJSON(m.Ubicacion),
		FechaAlta:       null.NewTime(m.FechaAlta.Time, m.FechaAlta.Valid),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

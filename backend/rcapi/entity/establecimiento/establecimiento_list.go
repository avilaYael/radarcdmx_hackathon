package establecimiento

import (
	entitytypes "github.com/mklfarha/radarcdmx/backend/rcapi/entity/types"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/contacto"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/ubicacion"
)

func (e Establecimiento) FieldIdentifierToTypeMap() map[string]entitytypes.FieldType {
	return map[string]entitytypes.FieldType{
		"uuid":             entitytypes.StringFieldType,
		"id_denue":         entitytypes.IntFieldType,
		"clee":             entitytypes.StringFieldType,
		"nombre":           entitytypes.StringFieldType,
		"razon_social":     entitytypes.StringFieldType,
		"per_ocu":          entitytypes.StringFieldType,
		"codigo_actividad": entitytypes.IntFieldType,
		"nombre_actividad": entitytypes.StringFieldType,
		"uso_de_suelo":     entitytypes.StringFieldType,
		"clave_catastral":  entitytypes.StringFieldType,
		"contacto":         entitytypes.MultiDependantEntityFieldType,
		"ubicacion":        entitytypes.MultiDependantEntityFieldType,
		"fecha_alta":       entitytypes.TimestampFieldType,
		"created_at":       entitytypes.TimestampFieldType,
		"updated_at":       entitytypes.TimestampFieldType,
	}
}

func (e Establecimiento) OrderedFieldIdentifiers() []string {
	res := []string{}
	res = append(res, "uuid")
	res = append(res, "id_denue")
	res = append(res, "clee")
	res = append(res, "nombre")
	res = append(res, "razon_social")
	res = append(res, "per_ocu")
	res = append(res, "codigo_actividad")
	res = append(res, "nombre_actividad")
	res = append(res, "uso_de_suelo")
	res = append(res, "clave_catastral")
	res = append(res, "contacto")
	res = append(res, "ubicacion")
	res = append(res, "fecha_alta")
	res = append(res, "created_at")
	res = append(res, "updated_at")

	return res
}

func (e Establecimiento) DependantFieldIdentifierToTypeMap() map[string]map[string]entitytypes.FieldType {
	res := make(map[string]map[string]entitytypes.FieldType)

	res["contacto"] = contacto.Contacto{}.FieldIdentifierToTypeMap()
	res["ubicacion"] = ubicacion.Ubicacion{}.FieldIdentifierToTypeMap()
	return res
}

func (e Establecimiento) EntityIdentifier() string {
	return "establecimiento"
}

func (e Establecimiento) PrimaryKeyIdentifiers() []string {
	return []string{
		"uuid",
	}
}

func (e Establecimiento) ArrayFieldIdentifierToType() map[string]entitytypes.FieldType {
	res := make(map[string]entitytypes.FieldType)

	return res
}

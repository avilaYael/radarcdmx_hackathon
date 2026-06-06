package ubicacion

import (
	entitytypes "github.com/mklfarha/radarcdmx/backend/rcapi/entity/types"
)

func (e Ubicacion) FieldIdentifierToTypeMap() map[string]entitytypes.FieldType {
	return map[string]entitytypes.FieldType{
		"entidad":       entitytypes.StringFieldType,
		"municipio":     entitytypes.StringFieldType,
		"localidad":     entitytypes.StringFieldType,
		"manzana":       entitytypes.IntFieldType,
		"codigo_postal": entitytypes.StringFieldType,
		"calle":         entitytypes.StringFieldType,
		"num_ext":       entitytypes.StringFieldType,
		"num_int":       entitytypes.StringFieldType,
		"latitud":       entitytypes.FloatFieldType,
		"longitud":      entitytypes.FloatFieldType,
	}
}

func (e Ubicacion) OrderedFieldIdentifiers() []string {
	res := []string{}
	res = append(res, "entidad")
	res = append(res, "municipio")
	res = append(res, "localidad")
	res = append(res, "manzana")
	res = append(res, "codigo_postal")
	res = append(res, "calle")
	res = append(res, "num_ext")
	res = append(res, "num_int")
	res = append(res, "latitud")
	res = append(res, "longitud")

	return res
}

func (e Ubicacion) DependantFieldIdentifierToTypeMap() map[string]map[string]entitytypes.FieldType {
	res := make(map[string]map[string]entitytypes.FieldType)

	return res
}

func (e Ubicacion) EntityIdentifier() string {
	return "ubicacion"
}

func (e Ubicacion) PrimaryKeyIdentifiers() []string {
	return []string{}
}

func (e Ubicacion) ArrayFieldIdentifierToType() map[string]entitytypes.FieldType {
	res := make(map[string]entitytypes.FieldType)

	return res
}

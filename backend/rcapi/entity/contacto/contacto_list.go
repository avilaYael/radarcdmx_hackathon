package contacto

import (
	entitytypes "github.com/mklfarha/radarcdmx/backend/rcapi/entity/types"
)

func (e Contacto) FieldIdentifierToTypeMap() map[string]entitytypes.FieldType {
	return map[string]entitytypes.FieldType{
		"telefono":  entitytypes.StringFieldType,
		"correo":    entitytypes.StringFieldType,
		"sitio_web": entitytypes.StringFieldType,
	}
}

func (e Contacto) OrderedFieldIdentifiers() []string {
	res := []string{}
	res = append(res, "telefono")
	res = append(res, "correo")
	res = append(res, "sitio_web")

	return res
}

func (e Contacto) DependantFieldIdentifierToTypeMap() map[string]map[string]entitytypes.FieldType {
	res := make(map[string]map[string]entitytypes.FieldType)

	return res
}

func (e Contacto) EntityIdentifier() string {
	return "contacto"
}

func (e Contacto) PrimaryKeyIdentifiers() []string {
	return []string{}
}

func (e Contacto) ArrayFieldIdentifierToType() map[string]entitytypes.FieldType {
	res := make(map[string]entitytypes.FieldType)

	return res
}

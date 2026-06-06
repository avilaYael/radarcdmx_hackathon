package user

import (
	entitytypes "github.com/mklfarha/radarcdmx/backend/rcapi/entity/types"
)

func (e User) FieldIdentifierToTypeMap() map[string]entitytypes.FieldType {
	return map[string]entitytypes.FieldType{
		"uuid":       entitytypes.StringFieldType,
		"name":       entitytypes.StringFieldType,
		"lastname":   entitytypes.StringFieldType,
		"email":      entitytypes.StringFieldType,
		"password":   entitytypes.StringFieldType,
		"status":     entitytypes.SingleEnumFieldType,
		"updated_at": entitytypes.TimestampFieldType,
		"created_by": entitytypes.StringFieldType,
		"updated_by": entitytypes.StringFieldType,
		"created_at": entitytypes.TimestampFieldType,
	}
}

func (e User) OrderedFieldIdentifiers() []string {
	res := []string{}
	res = append(res, "uuid")
	res = append(res, "name")
	res = append(res, "lastname")
	res = append(res, "email")
	res = append(res, "password")
	res = append(res, "status")
	res = append(res, "updated_at")
	res = append(res, "created_by")
	res = append(res, "updated_by")
	res = append(res, "created_at")

	return res
}

func (e User) DependantFieldIdentifierToTypeMap() map[string]map[string]entitytypes.FieldType {
	res := make(map[string]map[string]entitytypes.FieldType)

	return res
}

func (e User) EntityIdentifier() string {
	return "user"
}

func (e User) PrimaryKeyIdentifiers() []string {
	return []string{
		"uuid",
	}
}

func (e User) ArrayFieldIdentifierToType() map[string]entitytypes.FieldType {
	res := make(map[string]entitytypes.FieldType)

	return res
}

package user

import (
	rcapidb "github.com/mklfarha/radarcdmx/backend/rcapi/core/repository/gen"
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/user"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/mapper"

	"github.com/guregu/null/v6"
	"github.com/mklfarha/radarcdmx/backend/rcapi/enum"
)

func mapModelsToEntities(models []rcapidb.User) []main_entity.User {
	result := []main_entity.User{}
	for _, p := range models {
		result = append(result, mapModelToEntity(p))
	}
	return result
}

func mapModelToEntity(m rcapidb.User) main_entity.User {
	return main_entity.User{
		UUID:      mapper.StringToUUID(m.UUID),
		Name:      null.NewString(m.Name.String, m.Name.Valid),
		Lastname:  null.NewString(m.Lastname.String, m.Lastname.Valid),
		Email:     m.Email,
		Password:  m.Password,
		Status:    enum.UserStatus(m.Status),
		UpdatedAt: m.UpdatedAt,
		CreatedBy: mapper.StringToUUID(m.CreatedBy),
		UpdatedBy: mapper.StringToUUID(m.UpdatedBy),
		CreatedAt: m.CreatedAt,
	}
}

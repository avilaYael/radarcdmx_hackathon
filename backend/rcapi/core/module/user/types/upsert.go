package types

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/user"

	"github.com/gofrs/uuid"
)

type UpsertRequest struct {
	User main_entity.User
}

type UpsertResponse struct {
	UUID uuid.UUID
}

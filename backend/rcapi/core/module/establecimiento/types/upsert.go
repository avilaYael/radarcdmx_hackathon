package types

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"

	"github.com/gofrs/uuid"
)

type UpsertRequest struct {
	Establecimiento main_entity.Establecimiento
}

type UpsertResponse struct {
	UUID uuid.UUID
}

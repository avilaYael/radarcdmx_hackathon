package types

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/user"

	"github.com/gofrs/uuid"
	"go.uber.org/zap/zapcore"
)

type FetchUserByUuidRequest struct {
	UUID uuid.UUID
}

func (r FetchUserByUuidRequest) MarshalLogObject(e zapcore.ObjectEncoder) error {

	e.AddString("uuid", r.UUID.String())

	return nil
}

type FetchUserByUuidResponse struct {
	Results []main_entity.User
}

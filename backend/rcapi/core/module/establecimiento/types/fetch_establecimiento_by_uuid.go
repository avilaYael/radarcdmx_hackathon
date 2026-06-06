package types

import (
	main_entity "github.com/mklfarha/radarcdmx/backend/rcapi/entity/establecimiento"

	"github.com/gofrs/uuid"
	"go.uber.org/zap/zapcore"
)

type FetchEstablecimientoByUuidRequest struct {
	UUID uuid.UUID
}

func (r FetchEstablecimientoByUuidRequest) MarshalLogObject(e zapcore.ObjectEncoder) error {

	e.AddString("uuid", r.UUID.String())

	return nil
}

type FetchEstablecimientoByUuidResponse struct {
	Results []main_entity.Establecimiento
}

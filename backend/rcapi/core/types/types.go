package types

import (
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/repository"
	"go.uber.org/zap"
)

type ModuleParams struct {
	Repository *repository.Implementation
	Logger     *zap.Logger
}

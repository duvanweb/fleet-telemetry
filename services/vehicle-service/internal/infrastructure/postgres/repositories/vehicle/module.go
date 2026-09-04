package vehicle

import (
	"go.uber.org/fx"

	"fleet/vehicle-service/internal/core/ports/repositories"
	"fleet/shared/pkg/logger"
)

func newRepository(log logger.Logger, db repositories.Databaser) repositories.VehicleRepository {
	return NewRepository(log, Dependencies{DB: db})
}

// Module provides the vehicle repository to the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"vehicle_repository",
		fx.Provide(newRepository),
	)
}

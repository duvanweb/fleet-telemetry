package vehicle

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/logger"
	"fleet/vehicle-service/internal/core/ports/repositories"
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

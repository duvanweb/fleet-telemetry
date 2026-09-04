package main

import (
	"go.uber.org/fx"

	"fleet/vehicle-service/internal/core/health"
	"fleet/vehicle-service/internal/core/vehicle"
	"fleet/vehicle-service/internal/infrastructure/api/router"
	"fleet/vehicle-service/internal/infrastructure/postgres"
	vehiclerepo "fleet/vehicle-service/internal/infrastructure/postgres/repositories/vehicle"
	"fleet/shared/pkg/logger"
)

// Module assembles all FX modules for the vehicle-service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() (logger.Logger, error) {
			return logger.NewLogger()
		}),
		health.Module,
		postgres.Module(),
		vehiclerepo.Module(),
		vehicle.Module(),
		router.Module(),
	)
}

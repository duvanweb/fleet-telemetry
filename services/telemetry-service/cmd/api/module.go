package main

import (
	"go.uber.org/fx"

	"fleet/telemetry-service/internal/core/health"
	telemetrycore "fleet/telemetry-service/internal/core/telemetry"
	"fleet/telemetry-service/internal/infrastructure/api/router"
	vehicle_client "fleet/telemetry-service/internal/infrastructure/http/vehicle_client"
	"fleet/telemetry-service/internal/infrastructure/postgres"
	telemetryrepo "fleet/telemetry-service/internal/infrastructure/postgres/repositories/telemetry"
	"fleet/shared/pkg/logger"
)

// Module assembles all FX modules for the telemetry-service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() (logger.Logger, error) {
			return logger.NewLogger()
		}),
		health.Module,
		postgres.Module(),
		telemetryrepo.Module(),
		vehicle_client.Module(),
		telemetrycore.Module(),
		router.Module(),
	)
}

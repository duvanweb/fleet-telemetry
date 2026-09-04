package main

import (
	"go.uber.org/fx"

	"fleet/telemetry-service/internal/core/health"
	telemetrycore "fleet/telemetry-service/internal/core/telemetry"
	"fleet/telemetry-service/internal/infrastructure/api/router"
	vehicle_client "fleet/telemetry-service/internal/infrastructure/http/vehicle_client"
	"fleet/telemetry-service/internal/infrastructure/postgres"
	outboxrepo "fleet/telemetry-service/internal/infrastructure/postgres/repositories/outbox"
	telemetryrepo "fleet/telemetry-service/internal/infrastructure/postgres/repositories/telemetry"
	"fleet/telemetry-service/internal/infrastructure/rabbitmq"
	redisinfra "fleet/telemetry-service/internal/infrastructure/redis"
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
		outboxrepo.Module(),
		vehicle_client.Module(),
		redisinfra.Module(),
		telemetrycore.Module(),
		rabbitmq.Module(),
		router.Module(),
	)
}

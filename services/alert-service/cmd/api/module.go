package main

import (
	"go.uber.org/fx"

	"fleet/alert-service/internal/core/health"
	stopped_detection "fleet/alert-service/internal/core/stopped_detection"
	"fleet/alert-service/internal/infrastructure/api/router"
	"fleet/alert-service/internal/infrastructure/rabbitmq"
	"fleet/shared/pkg/logger"
)

// Module assembles all FX modules for the alert-service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() (logger.Logger, error) {
			return logger.NewLogger()
		}),
		health.Module,
		stopped_detection.Module(),
		rabbitmq.Module(),
		router.Module(),
	)
}

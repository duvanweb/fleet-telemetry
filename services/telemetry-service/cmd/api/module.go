package main

import (
	"go.uber.org/fx"

	"fleet/telemetry-service/internal/core/health"
	"fleet/telemetry-service/internal/infrastructure/api/router"
	"fleet/shared/pkg/logger"
)

// Module assembles all FX modules for the telemetry-service.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() (logger.Logger, error) {
			return logger.NewLogger()
		}),
		health.Module,
		router.Module(),
	)
}

package main

import (
	"go.uber.org/fx"

	"fleet/simulator/internal/core/health"
	"fleet/simulator/internal/core/simulator"
	"fleet/simulator/internal/infrastructure/api/router"
	telemetryclient "fleet/simulator/internal/infrastructure/http/telemetry_client"
	"fleet/shared/pkg/logger"
)

// Module assembles all FX modules for the simulator.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() (logger.Logger, error) {
			return logger.NewLogger()
		}),
		health.Module,
		telemetryclient.Module(),
		simulator.Module,
		router.Module(),
	)
}

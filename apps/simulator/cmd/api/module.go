package main

import (
	"go.uber.org/fx"

	"fleet/simulator/internal/core/health"
	"fleet/simulator/internal/infrastructure/api/router"
	"fleet/shared/pkg/logger"
)

// Module assembles all FX modules for the simulator.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(func() (logger.Logger, error) {
			return logger.NewLogger()
		}),
		health.Module,
		router.Module(),
	)
}

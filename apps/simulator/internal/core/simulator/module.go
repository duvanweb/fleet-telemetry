package simulator

import (
	"go.uber.org/fx"

	"fleet/simulator/internal/core/ports/services"
)

// Module registers the simulator service with FX.
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewService,
			fx.As(new(services.SimulatorService)),
		),
	),
)

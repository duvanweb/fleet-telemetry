package alerts

import (
	"go.uber.org/fx"

	"fleet/alert-service/internal/core/ports/services"
)

// Module wires the alerts domain module into the FX dependency graph.
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewService,
			fx.As(new(services.AlertService)),
		),
	),
)

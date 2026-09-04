package health

import (
	"go.uber.org/fx"

	"fleet/vehicle-service/internal/core/ports/services"
)

// Module wires the health domain into the FX dependency graph.
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewService, fx.As(new(services.HealthService))),
	),
)

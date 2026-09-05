package vehicle

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/logger"
	"fleet/vehicle-service/internal/core/ports/repositories"
	"fleet/vehicle-service/internal/core/ports/services"
)

func newService(log logger.Logger, repo repositories.VehicleRepository) *Service {
	return NewService(log, Repositories{Vehicle: repo})
}

// Module wires the vehicle service into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"vehicle",
		fx.Provide(
			fx.Annotate(newService, fx.As(new(services.VehicleService))),
		),
	)
}

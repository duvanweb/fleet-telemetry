package telemetry

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/logger"
	"fleet/telemetry-service/internal/core/ports/repositories"
	"fleet/telemetry-service/internal/core/ports/resources"
	"fleet/telemetry-service/internal/core/ports/services"
)

func newService(log logger.Logger, repo repositories.TelemetryRepository, checker resources.VehicleChecker, cache resources.TelemetryCache) *Service {
	return NewService(log, Repositories{Telemetry: repo}, Resources{VehicleChecker: checker, Cache: cache})
}

// Module wires the telemetry domain into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"telemetry",
		fx.Provide(fx.Annotate(newService, fx.As(new(services.TelemetryService)))),
	)
}

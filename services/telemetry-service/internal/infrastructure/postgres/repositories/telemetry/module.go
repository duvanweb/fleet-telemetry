package telemetry

import (
	"go.uber.org/fx"

	"fleet/telemetry-service/internal/core/ports/repositories"
	"fleet/shared/pkg/logger"
)

func newRepository(log logger.Logger, db repositories.Databaser) repositories.TelemetryRepository {
	return NewRepository(log, Dependencies{DB: db})
}

// Module wires the telemetry repository into the FX dependency graph.
func Module() fx.Option {
	return fx.Module("telemetry_repository", fx.Provide(newRepository))
}

package vehicle_client

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/env"
	"fleet/telemetry-service/internal/core/ports/resources"
)

func newClient(cfg *Configuration) resources.VehicleChecker {
	return NewClient(cfg)
}

// Module wires the vehicle HTTP client into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"vehicle_client",
		fx.Provide(
			env.LoadEnvConfiguration[Configuration],
			newClient,
		),
	)
}

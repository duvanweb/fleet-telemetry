package telemetry_client

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/env"
	"fleet/simulator/internal/core/ports/resources"
)

// Module registers the telemetry HTTP client with FX.
func Module() fx.Option {
	return fx.Module(
		"telemetry_client",
		fx.Provide(
			env.LoadEnvConfiguration[Configuration],
			fx.Annotate(
				NewClient,
				fx.As(new(resources.TelemetrySender)),
			),
		),
	)
}

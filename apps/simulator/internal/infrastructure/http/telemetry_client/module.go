package telemetry_client

import (
	"go.uber.org/fx"

	"fleet/simulator/internal/core/ports/resources"
)

// Module registers the telemetry HTTP client with FX.
func Module() fx.Option {
	return fx.Module(
		"telemetry_client",
		fx.Provide(
			func() *Configuration { return &Configuration{} },
			fx.Annotate(
				NewClient,
				fx.As(new(resources.TelemetrySender)),
			),
		),
	)
}

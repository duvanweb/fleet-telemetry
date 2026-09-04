package alert

import (
	"go.uber.org/fx"

	"fleet/alert-service/internal/core/ports/repositories"
)

// Module wires the alert repository into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"alert_repository",
		fx.Provide(
			fx.Annotate(
				NewRepository,
				fx.As(new(repositories.AlertRepository)),
			),
		),
	)
}

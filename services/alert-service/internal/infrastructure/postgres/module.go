package postgres

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/env"
)

// Module provides the PostgreSQL connection to the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"postgres",
		fx.Provide(
			env.LoadEnvConfiguration[Configuration],
			NewConnection,
		),
	)
}

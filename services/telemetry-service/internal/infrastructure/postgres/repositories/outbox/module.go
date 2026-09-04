package outbox

import (
	"go.uber.org/fx"

	"fleet/telemetry-service/internal/core/ports/repositories"
	"fleet/shared/pkg/logger"
)

func newRepository(log logger.Logger, db repositories.Databaser) repositories.OutboxRepository {
	return NewRepository(log, Dependencies{DB: db})
}

// Module wires the outbox repository into the FX dependency graph.
func Module() fx.Option {
	return fx.Module("outbox_repository", fx.Provide(newRepository))
}

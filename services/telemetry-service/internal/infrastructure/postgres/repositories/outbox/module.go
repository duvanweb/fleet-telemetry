package outbox

import (
	"go.uber.org/fx"

	"fleet/shared/pkg/logger"
	"fleet/telemetry-service/internal/core/ports/repositories"
)

func newRepository(log logger.Logger, db repositories.Databaser) repositories.OutboxRepository {
	return NewRepository(log, Dependencies{DB: db})
}

// Module wires the outbox repository into the FX dependency graph.
func Module() fx.Option {
	return fx.Module("outbox_repository", fx.Provide(newRepository))
}

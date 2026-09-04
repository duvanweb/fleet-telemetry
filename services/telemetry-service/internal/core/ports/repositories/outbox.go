package repositories

import (
	"context"

	"fleet/telemetry-service/internal/core/domain"
)

//go:generate mockery --name OutboxRepository --dir=. --output=./mocks

// OutboxRepository defines persistence operations for transactional outbox events.
type OutboxRepository interface {
	// GetPending returns up to limit pending outbox events, locking rows for this worker only.
	GetPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error)

	// MarkPublished marks the given event IDs as published with the current timestamp.
	MarkPublished(ctx context.Context, ids []string) error
}

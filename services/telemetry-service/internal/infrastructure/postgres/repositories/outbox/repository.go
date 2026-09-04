package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"

	"fleet/shared/pkg/logger"
	"fleet/telemetry-service/internal/core/domain"
	"fleet/telemetry-service/internal/core/ports/repositories"
	sqlqueries "fleet/telemetry-service/internal/infrastructure/postgres/repositories/outbox/sql"
)

// Dependencies holds the repository's injected dependencies.
type Dependencies struct {
	DB repositories.Databaser
}

// Repository implements OutboxRepository using PostgreSQL.
type Repository struct {
	logger logger.Logger
	db     repositories.Databaser
}

// NewRepository creates and returns a new outbox Repository.
func NewRepository(log logger.Logger, deps Dependencies) *Repository {
	return &Repository{logger: log, db: deps.DB}
}

// GetPending returns up to limit pending outbox events, locking rows for this worker only.
func (r *Repository) GetPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	rows, err := r.db.QueryContext(ctx, sqlqueries.SelectPending, limit)
	if err != nil {
		return nil, fmt.Errorf("querying pending outbox events: %w", err)
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var event domain.OutboxEvent
		var publishedAt *time.Time
		if err := rows.Scan(
			&event.ID,
			&event.EventType,
			&event.Payload,
			&event.Status,
			&event.CreatedAt,
			&publishedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning outbox row: %w", err)
		}
		event.PublishedAt = publishedAt
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating outbox rows: %w", err)
	}

	return events, nil
}

// MarkPublished marks the given event IDs as published with the current timestamp.
func (r *Repository) MarkPublished(ctx context.Context, ids []string) error {
	_, err := r.db.ExecContext(ctx, sqlqueries.UpdateMarkPublished, pq.Array(ids))
	if err != nil {
		return fmt.Errorf("marking outbox events as published: %w", err)
	}

	return nil
}

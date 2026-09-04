package rabbitmq

import (
	"context"
	"time"

	"fleet/telemetry-service/internal/core/ports/repositories"
	"fleet/shared/pkg/logger"
)

const (
	workerInterval = 500 * time.Millisecond
	workerBatch    = 10
)

// OutboxWorker polls the outbox table and publishes pending events to RabbitMQ.
type OutboxWorker struct {
	logger    logger.Logger
	outbox    repositories.OutboxRepository
	publisher *Publisher
}

// newOutboxWorker creates and returns a new OutboxWorker.
func newOutboxWorker(log logger.Logger, outbox repositories.OutboxRepository, pub *Publisher) *OutboxWorker {
	return &OutboxWorker{logger: log, outbox: outbox, publisher: pub}
}

func (w *OutboxWorker) start(ctx context.Context) {
	ticker := time.NewTicker(workerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *OutboxWorker) process(ctx context.Context) {
	events, err := w.outbox.GetPending(ctx, workerBatch)
	if err != nil {
		w.logger.Warnw(ctx, "failed to get pending outbox events", "error", err)
		return
	}

	if len(events) == 0 {
		return
	}

	var published []string
	for _, event := range events {
		if err := w.publisher.Publish(ctx, event.ID, event.EventType, event.Payload); err != nil {
			w.logger.Warnw(ctx, "failed to publish outbox event", "event_id", event.ID, "event_type", event.EventType, "error", err)
			continue
		}
		published = append(published, event.ID)
	}

	if len(published) == 0 {
		return
	}

	if err := w.outbox.MarkPublished(ctx, published); err != nil {
		w.logger.Warnw(ctx, "failed to mark outbox events as published", "count", len(published), "error", err)
	}
}

package domain

import "time"

const (
	OutboxStatusPending   = "pending"
	OutboxStatusPublished = "published"
)

// OutboxEvent represents a transactional outbox message pending delivery to the event bus.
type OutboxEvent struct {
	ID          string
	EventType   string
	Payload     []byte
	Status      string
	CreatedAt   time.Time
	PublishedAt *time.Time
}

package data

import (
	"time"

	"fleet/telemetry-service/internal/core/domain"
)

// GetTestOutboxEvent returns a test OutboxEvent fixture.
func GetTestOutboxEvent() domain.OutboxEvent {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return domain.OutboxEvent{
		ID:          "01HQ0000000000000000000001",
		EventType:   "telemetry.received",
		Payload:     []byte(`{"telemetry_id":"01HQ0000000000000000000000","vehicle_id":"vehicle-1","latitude":4.711,"longitude":-74.0721,"device_timestamp":"2024-01-01T00:00:00Z","received_at":"2024-01-01T00:00:00Z"}`),
		Status:      domain.OutboxStatusPending,
		CreatedAt:   now,
		PublishedAt: nil,
	}
}

package repositories

import (
	"context"

	"fleet/alert-service/internal/core/domain"
)

//go:generate mockery --name AlertRepository --dir=. --output=./mocks

// AlertRepository defines the persistence contract for alerts.
type AlertRepository interface {
	// Create persists a new alert and returns the saved record.
	Create(ctx context.Context, alert domain.Alert) (domain.Alert, error)
	// GetAll returns all alerts.
	GetAll(ctx context.Context) ([]domain.Alert, error)
	// GetOpenByVehicle returns the open alert for a vehicle, if any.
	GetOpenByVehicle(ctx context.Context, vehicleID string) (domain.Alert, bool, error)
	// MarkResolved sets status=RESOLVED and resolved_at=now for the given alert ID.
	MarkResolved(ctx context.Context, id string) error
}

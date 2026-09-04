package services

import (
	"context"
	"time"

	"fleet/alert-service/internal/core/domain"
)

//go:generate mockery --name AlertService --dir=. --output=./mocks

// AlertService defines the business operations for the alert lifecycle.
type AlertService interface {
	// CreateAlert creates a new open alert for a stopped vehicle.
	CreateAlert(ctx context.Context, vehicleID, alertType string, startedAt time.Time) (domain.Alert, error)

	// GetAll returns all alerts.
	GetAll(ctx context.Context) ([]domain.Alert, error)

	// GetOpenByVehicle returns the open alert for a vehicle, if any.
	GetOpenByVehicle(ctx context.Context, vehicleID string) (domain.Alert, bool, error)

	// ResolveAlert resolves the open alert for the given vehicle.
	ResolveAlert(ctx context.Context, vehicleID string) error
}

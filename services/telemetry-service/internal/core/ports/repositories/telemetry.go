package repositories

import (
	"context"

	"fleet/telemetry-service/internal/core/domain"
)

//go:generate mockery --name TelemetryRepository --dir=. --output=./mocks

// TelemetryRepository defines persistence operations for telemetry points.
type TelemetryRepository interface {
	// Create persists a new telemetry point and returns it with its assigned ID.
	Create(ctx context.Context, point domain.TelemetryPoint) (domain.TelemetryPoint, error)

	// GetByVehicleID returns all telemetry points for the given vehicle, ordered by device_timestamp ASC.
	GetByVehicleID(ctx context.Context, vehicleID string) ([]domain.TelemetryPoint, error)
}

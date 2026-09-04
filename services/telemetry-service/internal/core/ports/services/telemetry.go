package services

import (
	"context"

	"fleet/telemetry-service/internal/core/domain"
)

//go:generate mockery --name TelemetryService --dir=. --output=./mocks

// TelemetryService defines the use-cases for telemetry ingestion and retrieval.
type TelemetryService interface {
	// GetByVehicleID returns all telemetry points for the given vehicle.
	GetByVehicleID(ctx context.Context, vehicleID string) ([]domain.TelemetryPoint, error)

	// IngestTelemetry validates and persists a telemetry point.
	IngestTelemetry(ctx context.Context, point domain.TelemetryPoint) (domain.TelemetryPoint, error)
}

package services

import (
	"context"

	"fleet/alert-service/internal/core/domain"
)

//go:generate mockery --name StoppedDetectionService --dir=. --output=./mocks

// StoppedDetectionService evaluates telemetry events to detect stopped vehicles.
type StoppedDetectionService interface {
	// Evaluate processes a telemetry event and creates or resolves stopped alerts as needed.
	Evaluate(ctx context.Context, event domain.TelemetryEvent) error
}

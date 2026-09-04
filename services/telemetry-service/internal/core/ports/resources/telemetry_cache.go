package resources

import (
	"context"
	"time"

	"fleet/telemetry-service/internal/core/domain"
)

//go:generate mockery --name TelemetryCache --dir=. --output=./mocks

// TelemetryCache defines the caching operations for telemetry deduplication and last-position tracking.
type TelemetryCache interface {
	// CheckDedup returns true if the deduplication key already exists in the cache.
	CheckDedup(ctx context.Context, key string) (bool, error)

	// SetDedup persists the deduplication key with the given TTL.
	SetDedup(ctx context.Context, key string, ttl time.Duration) error

	// GetLastPosition retrieves the last known position for the given vehicle.
	// Returns the point, a found flag, and any error.
	GetLastPosition(ctx context.Context, vehicleID string) (domain.TelemetryPoint, bool, error)

	// SetLastPosition stores the given point as the last known position for the vehicle.
	SetLastPosition(ctx context.Context, vehicleID string, point domain.TelemetryPoint, ttl time.Duration) error
}

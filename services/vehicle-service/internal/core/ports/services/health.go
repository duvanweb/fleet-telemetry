package services

import (
	"context"

	"fleet/vehicle-service/internal/core/domain"
)

//go:generate mockery --name HealthService --dir=. --output=./mocks

// HealthService defines health check operations.
type HealthService interface {
	GetHealth(ctx context.Context) (domain.Health, error)
	GetReady(ctx context.Context) (domain.Health, error)
}

package services

import (
	"context"

	"fleet/telemetry-service/internal/core/domain"
)

//go:generate mockery --name HealthService --dir=. --output=./mocks
type HealthService interface {
	GetHealth(ctx context.Context) (domain.Health, error)
	GetReady(ctx context.Context) (domain.Health, error)
}

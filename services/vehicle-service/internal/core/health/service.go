package health

import (
	"context"

	"fleet/vehicle-service/internal/core/domain"
	"fleet/shared/pkg/logger"
)

// Service implements health check logic.
type Service struct {
	logger logger.Logger
}

// GetHealth returns the current liveness status of the service.
func (s *Service) GetHealth(ctx context.Context) (domain.Health, error) {
	return domain.Health{Status: "ok"}, nil
}

// GetReady returns the current readiness status of the service.
func (s *Service) GetReady(ctx context.Context) (domain.Health, error) {
	return domain.Health{Status: "ready"}, nil
}

// NewService creates and returns a new health Service.
func NewService(log logger.Logger) *Service {
	return &Service{logger: log}
}

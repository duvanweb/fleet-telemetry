package health

import (
	"context"

	"fleet/shared/pkg/logger"
	"fleet/telemetry-service/internal/core/domain"
)

// Service implements services.HealthService.
type Service struct {
	logger logger.Logger
}

// GetHealth returns the liveness status.
func (s *Service) GetHealth(_ context.Context) (domain.Health, error) {
	return domain.Health{Status: "ok"}, nil
}

// GetReady returns the readiness status.
func (s *Service) GetReady(_ context.Context) (domain.Health, error) {
	return domain.Health{Status: "ready"}, nil
}

// NewService creates and returns a new health Service.
func NewService(log logger.Logger) *Service {
	return &Service{logger: log}
}

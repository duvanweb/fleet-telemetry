package alerts

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/fx"

	"fleet/alert-service/internal/core/domain"
	"fleet/alert-service/internal/core/ports/repositories"
	"fleet/shared/pkg/logger"
)

// Dependencies holds injected dependencies for the alerts Service.
type Dependencies struct {
	fx.In

	Repository repositories.AlertRepository
}

// Service implements the AlertService business logic.
type Service struct {
	logger     logger.Logger
	repository repositories.AlertRepository
}

// NewService creates and returns a new alerts Service.
func NewService(log logger.Logger, deps Dependencies) *Service {
	return &Service{
		logger:     log,
		repository: deps.Repository,
	}
}

// CreateAlert persists a new open alert for the given vehicle. If an open alert already
// exists for the vehicle it returns the existing alert without creating a duplicate.
func (s *Service) CreateAlert(ctx context.Context, vehicleID, alertType string, startedAt time.Time) (domain.Alert, error) {
	existing, found, err := s.repository.GetOpenByVehicle(ctx, vehicleID)
	if err != nil {
		s.logger.Errorw(ctx, "failed to check for existing open alert", "vehicle_id", vehicleID, "error", err)
		return domain.Alert{}, err
	}

	if found {
		return existing, nil
	}

	now := time.Now().UTC()
	alert := domain.Alert{
		ID:        ulid.Make().String(),
		VehicleID: vehicleID,
		Type:      alertType,
		Status:    domain.AlertStatusOpen,
		StartedAt: startedAt,
		CreatedAt: now,
		UpdatedAt: now,
	}

	saved, err := s.repository.Create(ctx, alert)
	if err != nil {
		s.logger.Errorw(ctx, "failed to create alert", "vehicle_id", vehicleID, "alert_type", alertType, "error", err)
		return domain.Alert{}, err
	}

	return saved, nil
}

// GetAll returns all alerts regardless of status.
func (s *Service) GetAll(ctx context.Context) ([]domain.Alert, error) {
	alerts, err := s.repository.GetAll(ctx)
	if err != nil {
		s.logger.Errorw(ctx, "failed to get all alerts", "error", err)
		return nil, err
	}

	return alerts, nil
}

// GetOpenByVehicle returns the open alert for a vehicle, if one exists.
func (s *Service) GetOpenByVehicle(ctx context.Context, vehicleID string) (domain.Alert, bool, error) {
	alert, found, err := s.repository.GetOpenByVehicle(ctx, vehicleID)
	if err != nil {
		s.logger.Errorw(ctx, "failed to get open alert by vehicle", "vehicle_id", vehicleID, "error", err)
		return domain.Alert{}, false, err
	}

	return alert, found, nil
}

// ResolveAlert marks the open alert for the given vehicle as resolved.
// If no open alert exists, it returns nil (no-op).
func (s *Service) ResolveAlert(ctx context.Context, vehicleID string) error {
	existing, found, err := s.repository.GetOpenByVehicle(ctx, vehicleID)
	if err != nil {
		s.logger.Errorw(ctx, "failed to find open alert for resolution", "vehicle_id", vehicleID, "error", err)
		return err
	}

	if !found {
		return nil
	}

	if err := s.repository.MarkResolved(ctx, existing.ID); err != nil {
		s.logger.Errorw(ctx, "failed to mark alert as resolved", "alert_id", existing.ID, "vehicle_id", vehicleID, "error", err)
		return err
	}

	return nil
}

package stopped_detection

import (
	"context"
	"sync"
	"time"

	"fleet/alert-service/internal/core/domain"
	"fleet/alert-service/internal/core/ports/services"
	"fleet/shared/pkg/logger"
)

// Service evaluates telemetry events to detect stopped vehicles.
type Service struct {
	logger    logger.Logger
	mu        sync.Mutex
	states    map[string]*domain.VehicleState
	alerts    services.AlertService
	threshold time.Duration
	now       func() time.Time
}

// NewService creates and returns a new stopped detection Service.
func NewService(log logger.Logger, alerts services.AlertService, threshold time.Duration, now func() time.Time) *Service {
	return &Service{
		logger:    log,
		states:    make(map[string]*domain.VehicleState),
		alerts:    alerts,
		threshold: threshold,
		now:       now,
	}
}

// Evaluate processes a telemetry event and creates or resolves stopped alerts as needed.
func (s *Service) Evaluate(ctx context.Context, event domain.TelemetryEvent) error {
	now := s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	state, ok := s.states[event.VehicleID]
	if !ok {
		state = &domain.VehicleState{}
		s.states[event.VehicleID] = state
	}

	samePos := event.Latitude == state.LastLatitude && event.Longitude == state.LastLongitude

	if samePos {
		if !state.HasOpenAlert {
			if state.FirstSamePositionAt.IsZero() {
				state.FirstSamePositionAt = now
			} else if now.Sub(state.FirstSamePositionAt) >= s.threshold {
				_, err := s.alerts.CreateAlert(ctx, event.VehicleID, domain.AlertTypeVehicleStopped, state.FirstSamePositionAt)
				if err != nil {
					s.logger.Errorw(ctx, "failed to create stopped alert", "vehicle_id", event.VehicleID, "error", err)
					return err
				}
				state.HasOpenAlert = true
			}
		}
	} else {
		if state.HasOpenAlert {
			if err := s.alerts.ResolveAlert(ctx, event.VehicleID); err != nil {
				s.logger.Errorw(ctx, "failed to resolve stopped alert", "vehicle_id", event.VehicleID, "error", err)
				return err
			}
			state.HasOpenAlert = false
		}
		state.FirstSamePositionAt = time.Time{}
		state.LastLatitude = event.Latitude
		state.LastLongitude = event.Longitude
	}

	return nil
}

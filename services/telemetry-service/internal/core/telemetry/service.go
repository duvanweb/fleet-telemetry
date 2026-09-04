package telemetry

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"fleet/telemetry-service/internal/core/domain"
	"fleet/telemetry-service/internal/core/ports/repositories"
	"fleet/telemetry-service/internal/core/ports/resources"
	"fleet/shared/pkg/logger"
)

const (
	maxFutureDuration   = 5 * time.Minute
	dedupTTL            = 60 * time.Second
	lastPositionTTL     = 10 * time.Minute
)

// Repositories holds repository dependencies for the telemetry service.
type Repositories struct {
	Telemetry repositories.TelemetryRepository
}

// Resources holds external resource dependencies for the telemetry service.
type Resources struct {
	VehicleChecker resources.VehicleChecker
	Cache          resources.TelemetryCache
}

// Service implements TelemetryService.
type Service struct {
	logger       logger.Logger
	repositories Repositories
	resources    Resources
}

// NewService creates and returns a new telemetry Service.
func NewService(log logger.Logger, repos Repositories, res Resources) *Service {
	return &Service{logger: log, repositories: repos, resources: res}
}

// GetByVehicleID returns all telemetry points for the given vehicle.
func (s *Service) GetByVehicleID(ctx context.Context, vehicleID string) ([]domain.TelemetryPoint, error) {
	points, err := s.repositories.Telemetry.GetByVehicleID(ctx, vehicleID)
	if err != nil {
		s.logger.Errorw(ctx, "failed to get telemetry by vehicle id", "vehicle_id", vehicleID, "error", err)
		return nil, err
	}

	return points, nil
}

// IngestTelemetry validates and persists a GPS telemetry point.
func (s *Service) IngestTelemetry(ctx context.Context, point domain.TelemetryPoint) (domain.TelemetryPoint, error) {
	if err := validateCoordinates(point.Latitude, point.Longitude); err != nil {
		return domain.TelemetryPoint{}, err
	}

	if err := validateTimestamp(point.DeviceTimestamp); err != nil {
		return domain.TelemetryPoint{}, err
	}

	active, err := s.resources.VehicleChecker.ExistsAndActive(ctx, point.VehicleID)
	if err != nil {
		s.logger.Errorw(ctx, "failed to check vehicle existence", "vehicle_id", point.VehicleID, "error", err)
		return domain.TelemetryPoint{}, err
	}

	if !active {
		return domain.TelemetryPoint{}, domain.ErrVehicleNotFound
	}

	point.DeduplicationKey = computeDeduplicationKey(point)

	isDup, err := s.resources.Cache.CheckDedup(ctx, point.DeduplicationKey)
	if err != nil {
		s.logger.Errorw(ctx, "failed to check dedup key", "vehicle_id", point.VehicleID, "error", err)
		return domain.TelemetryPoint{}, err
	}

	if isDup {
		return domain.TelemetryPoint{}, domain.ErrDuplicateTelemetry
	}

	if err := s.resources.Cache.SetDedup(ctx, point.DeduplicationKey, dedupTTL); err != nil {
		s.logger.Errorw(ctx, "failed to set dedup key", "vehicle_id", point.VehicleID, "error", err)
		return domain.TelemetryPoint{}, err
	}

	point.ReceivedAt = time.Now().UTC()

	saved, err := s.repositories.Telemetry.Create(ctx, point)
	if err != nil {
		s.logger.Errorw(ctx, "failed to persist telemetry point", "vehicle_id", point.VehicleID, "error", err)
		return domain.TelemetryPoint{}, err
	}

	if err := s.updateLastPosition(ctx, saved); err != nil {
		s.logger.Warnw(ctx, "failed to update last position cache", "vehicle_id", saved.VehicleID, "error", err)
	}

	return saved, nil
}

func (s *Service) updateLastPosition(ctx context.Context, point domain.TelemetryPoint) error {
	last, found, err := s.resources.Cache.GetLastPosition(ctx, point.VehicleID)
	if err != nil {
		return err
	}

	if found && point.DeviceTimestamp.Before(last.DeviceTimestamp) {
		return nil
	}

	return s.resources.Cache.SetLastPosition(ctx, point.VehicleID, point, lastPositionTTL)
}

func validateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return domain.ErrInvalidLatitude
	}

	if lon < -180 || lon > 180 {
		return domain.ErrInvalidLongitude
	}

	return nil
}

func validateTimestamp(ts time.Time) error {
	if ts.IsZero() {
		return domain.ErrInvalidTimestamp
	}

	if time.Since(ts) < -maxFutureDuration {
		return domain.ErrFutureTelemetry
	}

	return nil
}

func computeDeduplicationKey(point domain.TelemetryPoint) string {
	raw := point.VehicleID +
		fmt.Sprintf("%.7f", point.Latitude) +
		fmt.Sprintf("%.7f", point.Longitude) +
		point.DeviceTimestamp.UTC().Format(time.RFC3339Nano)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

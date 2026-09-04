package vehicle

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"

	"fleet/vehicle-service/internal/core/domain"
	"fleet/vehicle-service/internal/core/ports/repositories"
	"fleet/shared/pkg/logger"
)

// Repositories holds the repository dependencies for the vehicle service.
type Repositories struct {
	Vehicle repositories.VehicleRepository
}

// Service implements services.VehicleService.
type Service struct {
	logger       logger.Logger
	repositories Repositories
}

// Create validates uniqueness and persists a new vehicle.
func (s *Service) Create(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
	exists, err := s.repositories.Vehicle.ExistsByPlate(ctx, v.Plate)
	if err != nil {
		s.logger.Errorw(ctx, "failed to check plate existence", "plate", v.Plate, "error", err)
		return domain.Vehicle{}, fmt.Errorf("checking plate existence: %w", err)
	}

	if exists {
		return domain.Vehicle{}, domain.ErrDuplicatePlate
	}

	exists, err = s.repositories.Vehicle.ExistsByExternalID(ctx, v.ExternalID)
	if err != nil {
		s.logger.Errorw(ctx, "failed to check external id existence", "external_id", v.ExternalID, "error", err)
		return domain.Vehicle{}, fmt.Errorf("checking external id existence: %w", err)
	}

	if exists {
		return domain.Vehicle{}, domain.ErrDuplicatePlate
	}

	now := time.Now().UTC()
	v.ID = ulid.Make().String()
	v.CreatedAt = now
	v.UpdatedAt = now

	created, err := s.repositories.Vehicle.Create(ctx, v)
	if err != nil {
		s.logger.Errorw(ctx, "failed to create vehicle", "plate", v.Plate, "error", err)
		return domain.Vehicle{}, fmt.Errorf("creating vehicle: %w", err)
	}

	return created, nil
}

// Delete soft-deletes the vehicle with the given ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repositories.Vehicle.SoftDelete(ctx, id); err != nil {
		s.logger.Errorw(ctx, "failed to delete vehicle", "vehicle_id", id, "error", err)
		return err
	}

	return nil
}

// GetAll returns all non-deleted vehicles.
func (s *Service) GetAll(ctx context.Context) ([]domain.Vehicle, error) {
	vehicles, err := s.repositories.Vehicle.GetAll(ctx)
	if err != nil {
		s.logger.Errorw(ctx, "failed to get all vehicles", "error", err)
		return nil, fmt.Errorf("getting all vehicles: %w", err)
	}

	return vehicles, nil
}

// GetByID returns the vehicle with the given ID, rejecting deleted vehicles.
func (s *Service) GetByID(ctx context.Context, id string) (domain.Vehicle, error) {
	v, err := s.repositories.Vehicle.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw(ctx, "failed to get vehicle by id", "vehicle_id", id, "error", err)
		return domain.Vehicle{}, err
	}

	if v.DeletedAt != nil {
		return domain.Vehicle{}, domain.ErrVehicleDeleted
	}

	return v, nil
}

// Update applies plate and name changes to the given vehicle.
func (s *Service) Update(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
	v.UpdatedAt = time.Now().UTC()

	updated, err := s.repositories.Vehicle.Update(ctx, v)
	if err != nil {
		s.logger.Errorw(ctx, "failed to update vehicle", "vehicle_id", v.ID, "error", err)
		return domain.Vehicle{}, err
	}

	return updated, nil
}

// NewService creates and returns a new vehicle Service.
func NewService(log logger.Logger, repos Repositories) *Service {
	return &Service{logger: log, repositories: repos}
}

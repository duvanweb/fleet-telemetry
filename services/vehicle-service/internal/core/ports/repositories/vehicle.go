package repositories

import (
	"context"

	"fleet/vehicle-service/internal/core/domain"
)

// VehicleRepository defines persistence operations for vehicles.
//
//go:generate mockery --name VehicleRepository --dir=. --output=./mocks
type VehicleRepository interface {
	// Create persists a new vehicle and returns it with its assigned ID.
	Create(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error)

	// ExistsByExternalID reports whether a non-deleted vehicle with the given external ID exists.
	ExistsByExternalID(ctx context.Context, externalID string) (bool, error)

	// ExistsByPlate reports whether a non-deleted vehicle with the given plate exists.
	ExistsByPlate(ctx context.Context, plate string) (bool, error)

	// GetAll returns all non-deleted vehicles.
	GetAll(ctx context.Context) ([]domain.Vehicle, error)

	// GetByID returns the vehicle with the given ID regardless of deletion status.
	GetByID(ctx context.Context, id string) (domain.Vehicle, error)

	// SoftDelete marks the vehicle with the given ID as deleted.
	SoftDelete(ctx context.Context, id string) error

	// Update persists changes to an existing vehicle.
	Update(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error)
}

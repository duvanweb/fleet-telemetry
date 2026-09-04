package services

import (
	"context"

	"fleet/vehicle-service/internal/core/domain"
)

//go:generate mockery --name VehicleService --dir=. --output=./mocks
// VehicleService defines the business operations for managing vehicles.
type VehicleService interface {
	// Create validates and persists a new vehicle.
	Create(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error)

	// Delete soft-deletes the vehicle with the given ID.
	Delete(ctx context.Context, id string) error

	// GetAll returns all non-deleted vehicles.
	GetAll(ctx context.Context) ([]domain.Vehicle, error)

	// GetByID returns the vehicle with the given ID.
	GetByID(ctx context.Context, id string) (domain.Vehicle, error)

	// Update applies plate and name changes to the given vehicle.
	Update(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error)
}

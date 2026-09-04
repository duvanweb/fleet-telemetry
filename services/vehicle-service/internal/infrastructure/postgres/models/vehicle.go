package models

import (
	"time"

	"fleet/vehicle-service/internal/core/domain"
)

// VehicleRow maps a vehicles table row to a Go struct.
type VehicleRow struct {
	ID         string
	ExternalID string
	Plate      string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

// ToDomain converts a VehicleRow to the domain Vehicle entity.
func (r VehicleRow) ToDomain() domain.Vehicle {
	return domain.Vehicle{
		ID:         r.ID,
		ExternalID: r.ExternalID,
		Plate:      r.Plate,
		Name:       r.Name,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
		DeletedAt:  r.DeletedAt,
	}
}

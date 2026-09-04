package data

import (
	"time"

	"fleet/vehicle-service/internal/core/domain"
)

// GetTestVehicle returns a fixed Vehicle fixture for use in tests.
func GetTestVehicle() domain.Vehicle {
	createdAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return domain.Vehicle{
		ID:         "01HQZK9M0XVNP3F7BRDT4YEWCX",
		ExternalID: "EXT-001",
		Plate:      "ABC-123",
		Name:       "Test Vehicle",
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
		DeletedAt:  nil,
	}
}

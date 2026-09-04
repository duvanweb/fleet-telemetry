package resources

import "context"

//go:generate mockery --name VehicleChecker --dir=. --output=./mocks

// VehicleChecker verifies whether a vehicle exists and is active (not deleted).
type VehicleChecker interface {
	// ExistsAndActive returns true if the vehicle with the given ID exists and has not been deleted.
	ExistsAndActive(ctx context.Context, vehicleID string) (bool, error)
}

package domain

import "time"

const (
	AlertTypeVehicleStopped = "VEHICLE_STOPPED"
	AlertStatusOpen         = "OPEN"
	AlertStatusResolved     = "RESOLVED"
)

// Alert represents a vehicle alert event with its lifecycle.
type Alert struct {
	ID         string
	VehicleID  string
	Type       string
	Status     string
	StartedAt  time.Time
	ResolvedAt *time.Time
	Metadata   map[string]string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

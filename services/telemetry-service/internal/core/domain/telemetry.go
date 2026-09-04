package domain

import "time"

// TelemetryPoint represents a single GPS reading from a vehicle.
type TelemetryPoint struct {
	ID               string
	VehicleID        string
	Latitude         float64
	Longitude        float64
	DeviceTimestamp  time.Time
	ReceivedAt       time.Time
	DeduplicationKey string
}

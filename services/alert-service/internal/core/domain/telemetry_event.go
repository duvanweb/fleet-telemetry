package domain

import "time"

// TelemetryEvent represents a GPS reading received from the event bus.
type TelemetryEvent struct {
	TelemetryID     string    `json:"telemetry_id"`
	VehicleID       string    `json:"vehicle_id"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	DeviceTimestamp time.Time `json:"device_timestamp"`
	ReceivedAt      time.Time `json:"received_at"`
}

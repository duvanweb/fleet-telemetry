package resources

import (
	"context"
	"time"
)

//go:generate mockery --name TelemetrySender --dir=. --output=./mocks

// TelemetrySendRequest holds a single telemetry point to send.
type TelemetrySendRequest struct {
	VehicleID       string
	Latitude        float64
	Longitude       float64
	DeviceTimestamp time.Time
}

// TelemetrySender sends telemetry points to the telemetry service.
type TelemetrySender interface {
	Send(ctx context.Context, req TelemetrySendRequest) error
}

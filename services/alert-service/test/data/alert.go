package data

import (
	"time"

	"fleet/alert-service/internal/core/domain"
)

// GetTestAlert returns a test Alert fixture.
func GetTestAlert() domain.Alert {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return domain.Alert{
		ID:        "01HQ0000000000000000000001",
		VehicleID: "vehicle-1",
		Type:      domain.AlertTypeVehicleStopped,
		Status:    domain.AlertStatusOpen,
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// GetTestResolvedAlert returns a test Alert fixture with status RESOLVED.
func GetTestResolvedAlert() domain.Alert {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	resolved := now.Add(time.Hour)
	return domain.Alert{
		ID:         "01HQ0000000000000000000001",
		VehicleID:  "vehicle-1",
		Type:       domain.AlertTypeVehicleStopped,
		Status:     domain.AlertStatusResolved,
		StartedAt:  now,
		ResolvedAt: &resolved,
		CreatedAt:  now,
		UpdatedAt:  resolved,
	}
}

// GetTestTelemetryEvent returns a test TelemetryEvent fixture.
func GetTestTelemetryEvent() domain.TelemetryEvent {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return domain.TelemetryEvent{
		TelemetryID:     "01HQ0000000000000000000000",
		VehicleID:       "vehicle-1",
		Latitude:        4.711,
		Longitude:       -74.0721,
		DeviceTimestamp: now,
		ReceivedAt:      now,
	}
}

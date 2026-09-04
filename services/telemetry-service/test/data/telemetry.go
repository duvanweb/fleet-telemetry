package data

import (
	"time"

	"fleet/telemetry-service/internal/core/domain"
)

// GetTestTelemetryPoint returns a fixed TelemetryPoint fixture for use in tests.
func GetTestTelemetryPoint() domain.TelemetryPoint {
	deviceTS := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	receivedAt := time.Date(2024, 1, 15, 10, 0, 1, 0, time.UTC)
	return domain.TelemetryPoint{
		ID:               "01HQZK9M0XVNP3F7BRDT4YEWCA",
		VehicleID:        "01HQZK9M0XVNP3F7BRDT4YEWCX",
		Latitude:         7.1193,
		Longitude:        -73.1227,
		DeviceTimestamp:  deviceTS,
		ReceivedAt:       receivedAt,
		DeduplicationKey: "abc123dedup",
	}
}

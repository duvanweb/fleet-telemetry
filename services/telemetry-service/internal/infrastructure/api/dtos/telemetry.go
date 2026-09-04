package dtos

// IngestTelemetryRequest is the request body for ingesting a telemetry point.
type IngestTelemetryRequest struct {
	VehicleID       string  `json:"vehicle_id"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	DeviceTimestamp string  `json:"device_timestamp"`
}

// TelemetryPointResponse is the JSON representation of a telemetry point.
type TelemetryPointResponse struct {
	ID               string  `json:"id"`
	VehicleID        string  `json:"vehicle_id"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	DeviceTimestamp  string  `json:"device_timestamp"`
	ReceivedAt       string  `json:"received_at"`
	DeduplicationKey string  `json:"deduplication_key"`
}

// VehicleTelemetryResponse is the response for listing a vehicle's telemetry points.
type VehicleTelemetryResponse struct {
	VehicleID string                   `json:"vehicle_id"`
	Points    []TelemetryPointResponse `json:"points"`
}

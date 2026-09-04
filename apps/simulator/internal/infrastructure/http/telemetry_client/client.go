package telemetry_client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"fleet/shared/pkg/logger"
	"fleet/simulator/internal/core/ports/resources"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Configuration holds the telemetry service connection settings.
type Configuration struct {
	TelemetryServiceURL string `env:"TELEMETRY_SERVICE_URL" envDefault:"http://localhost:8082"`
}

type ingestRequest struct {
	VehicleID       string  `json:"vehicle_id"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	DeviceTimestamp string  `json:"device_timestamp"`
}

// Client sends telemetry points to the telemetry-service via HTTP.
type Client struct {
	logger     logger.Logger
	baseURL    string
	httpClient *http.Client
}

// NewClient creates and returns a new telemetry HTTP Client.
func NewClient(log logger.Logger, cfg *Configuration) *Client {
	return &Client{
		logger:  log,
		baseURL: cfg.TelemetryServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send sends a single telemetry point to the telemetry service.
func (c *Client) Send(ctx context.Context, req resources.TelemetrySendRequest) error {
	body := ingestRequest{
		VehicleID:       req.VehicleID,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		DeviceTimestamp: req.DeviceTimestamp.UTC().Format(time.RFC3339Nano),
	}

	data, err := json.Marshal(body)
	if err != nil {
		c.logger.Errorw(ctx, "failed to marshal telemetry request", "error", err)
		return fmt.Errorf("marshaling telemetry: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/telemetry", bytes.NewReader(data))
	if err != nil {
		c.logger.Errorw(ctx, "failed to build telemetry request", "error", err)
		return fmt.Errorf("building telemetry request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Errorw(ctx, "failed to send telemetry", "vehicle_id", req.VehicleID, "error", err)
		return fmt.Errorf("sending telemetry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("telemetry service error: status %d", resp.StatusCode)
	}

	return nil
}

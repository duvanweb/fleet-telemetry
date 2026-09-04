package vehicle_client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Client implements the VehicleChecker resource using the vehicle-service HTTP API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates and returns a new vehicle service HTTP Client.
func NewClient(cfg *Configuration) *Client {
	return &Client{
		baseURL:    cfg.VehicleServiceURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// ExistsAndActive returns true if the vehicle exists and has not been deleted.
// It treats a 200 response as active, 404 as not found/deleted, and any other
// status as an error.
func (c *Client) ExistsAndActive(ctx context.Context, vehicleID string) (bool, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/api/vehicles/%s", c.baseURL, vehicleID)
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("building vehicle check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("calling vehicle service: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected vehicle service response: %d", resp.StatusCode)
	}
}

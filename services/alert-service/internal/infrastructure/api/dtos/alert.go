package dtos

import "time"

// AlertResponse is the API response for a single alert.
type AlertResponse struct {
	ID         string     `json:"id"`
	VehicleID  string     `json:"vehicle_id"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// AlertListResponse is the API response for a list of alerts.
type AlertListResponse struct {
	Alerts []AlertResponse `json:"alerts"`
}

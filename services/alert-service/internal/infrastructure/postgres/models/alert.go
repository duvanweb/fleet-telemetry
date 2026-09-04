package models

import (
	"database/sql"
	"time"

	"fleet/alert-service/internal/core/domain"
)

// AlertRow maps a database row to the Alert domain entity.
type AlertRow struct {
	ID         string
	VehicleID  string
	Type       string
	Status     string
	StartedAt  time.Time
	ResolvedAt sql.NullTime
	Metadata   []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ToDomain converts the AlertRow to a domain.Alert.
func (r AlertRow) ToDomain() domain.Alert {
	alert := domain.Alert{
		ID:        r.ID,
		VehicleID: r.VehicleID,
		Type:      r.Type,
		Status:    r.Status,
		StartedAt: r.StartedAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if r.ResolvedAt.Valid {
		alert.ResolvedAt = &r.ResolvedAt.Time
	}

	return alert
}

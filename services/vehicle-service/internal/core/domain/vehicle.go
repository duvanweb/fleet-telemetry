package domain

import "time"

// Vehicle represents a fleet vehicle entity.
type Vehicle struct {
	ID         string
	ExternalID string
	Plate      string
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

package domain

import "time"

// VehicleState tracks in-memory position and alert state for stopped-vehicle detection.
type VehicleState struct {
	LastLatitude        float64
	LastLongitude       float64
	FirstSamePositionAt time.Time
	HasOpenAlert        bool
}

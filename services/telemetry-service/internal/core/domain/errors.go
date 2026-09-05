package domain

import "errors"

var (
	ErrInvalidLatitude    = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude   = errors.New("longitude must be between -180 and 180")
	ErrInvalidTimestamp   = errors.New("device timestamp is invalid")
	ErrFutureTelemetry    = errors.New("device timestamp is too far in the future")
	ErrVehicleNotFound    = errors.New("vehicle not found or has been deleted")
	ErrDuplicateTelemetry = errors.New("duplicate telemetry point")
)

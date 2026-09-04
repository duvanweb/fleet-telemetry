package domain

import "errors"

// ErrDuplicatePlate is returned when a vehicle with the same plate already exists.
var ErrDuplicatePlate = errors.New("vehicle with this plate already exists")

// ErrVehicleDeleted is returned when an operation is attempted on a soft-deleted vehicle.
var ErrVehicleDeleted = errors.New("vehicle has been deleted")

// ErrVehicleNotFound is returned when a vehicle cannot be found.
var ErrVehicleNotFound = errors.New("vehicle not found")

package domain

import "time"

// Coordinate represents a GPS coordinate pair.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// Scenario defines a named simulation path with coordinate waypoints.
type Scenario struct {
	Name        string
	Description string
	Points      []Coordinate
}

// SimulationStatus describes the current state of the simulator.
type SimulationStatus struct {
	Running       bool
	Scenario      string
	VehicleCount  int
	IntervalMs    int
	DuplicateRate float64
	InvalidRate   float64
	StartedAt     *time.Time
}

// StartRequest holds parameters for starting a simulation.
type StartRequest struct {
	VehicleIDs    []string
	ScenarioName  string
	IntervalMs    int
	DuplicateRate float64
	InvalidRate   float64
}

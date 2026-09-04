package dtos

// StartSimulationRequest is the request body for starting a simulation.
type StartSimulationRequest struct {
	VehicleIDs    []string `json:"vehicle_ids"`
	ScenarioName  string   `json:"scenario"`
	IntervalMs    int      `json:"interval_ms"`
	DuplicateRate float64  `json:"duplicate_rate"`
	InvalidRate   float64  `json:"invalid_rate"`
}

// ScenarioResponse represents a simulation scenario.
type ScenarioResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PointCount  int    `json:"point_count"`
}

// SimulationStatusResponse represents the current simulator status.
type SimulationStatusResponse struct {
	Running       bool    `json:"running"`
	Scenario      string  `json:"scenario,omitempty"`
	VehicleCount  int     `json:"vehicle_count,omitempty"`
	IntervalMs    int     `json:"interval_ms,omitempty"`
	DuplicateRate float64 `json:"duplicate_rate,omitempty"`
	InvalidRate   float64 `json:"invalid_rate,omitempty"`
	StartedAt     string  `json:"started_at,omitempty"`
}

// ScenariosResponse wraps a list of scenarios.
type ScenariosResponse struct {
	Scenarios []ScenarioResponse `json:"scenarios"`
}

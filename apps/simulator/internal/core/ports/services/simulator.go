package services

import (
	"context"

	"fleet/simulator/internal/core/domain"
)

//go:generate mockery --name SimulatorService --dir=. --output=./mocks

// SimulatorService defines the simulator control operations.
type SimulatorService interface {
	// GetScenarios returns all available simulation scenarios.
	GetScenarios() []domain.Scenario
	// Start begins a simulation run with the given parameters.
	Start(ctx context.Context, req domain.StartRequest) error
	// Status returns the current simulation status.
	Status() domain.SimulationStatus
	// Stop halts the current simulation run.
	Stop()
}

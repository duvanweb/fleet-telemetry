package simulator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"

	"fleet/shared/pkg/logger"
	"fleet/simulator/internal/core/domain"
	"fleet/simulator/internal/core/ports/resources"
	"fleet/simulator/internal/core/scenario"
)

// Service implements the simulator control logic.
type Service struct {
	logger logger.Logger
	sender resources.TelemetrySender

	mu     sync.Mutex
	status domain.SimulationStatus
	cancel context.CancelFunc
}

// Dependencies holds injected dependencies.
type Dependencies struct {
	fx.In

	Sender resources.TelemetrySender
}

// NewService creates and returns a new simulator Service.
func NewService(log logger.Logger, deps Dependencies) *Service {
	return &Service{
		logger: log,
		sender: deps.Sender,
	}
}

// GetScenarios returns all available simulation scenarios.
func (s *Service) GetScenarios() []domain.Scenario {
	return scenario.All()
}

// Start begins a simulation run. Returns an error if one is already running or the scenario is unknown.
func (s *Service) Start(ctx context.Context, req domain.StartRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.Running {
		return fmt.Errorf("simulation already running")
	}

	sc, found := scenario.ByName(req.ScenarioName)
	if !found {
		return fmt.Errorf("unknown scenario: %s", req.ScenarioName)
	}

	if req.IntervalMs <= 0 {
		req.IntervalMs = 5000
	}

	// Use context.Background() so the simulation goroutine outlives the HTTP request lifecycle.
	// r.Context() is cancelled by Go's net/http when ServeHTTP returns, which would kill the
	// goroutine immediately. The caller's ctx is only used for the synchronous setup above.
	runCtx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	now := time.Now().UTC()
	s.status = domain.SimulationStatus{
		Running:       true,
		Scenario:      req.ScenarioName,
		VehicleCount:  len(req.VehicleIDs),
		IntervalMs:    req.IntervalMs,
		DuplicateRate: req.DuplicateRate,
		InvalidRate:   req.InvalidRate,
		StartedAt:     &now,
	}

	go s.run(runCtx, req, sc)

	return nil
}

// Status returns the current simulation status.
func (s *Service) Status() domain.SimulationStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

// Stop halts the current simulation run.
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.status = domain.SimulationStatus{Running: false}
}

// run starts goroutines per vehicle using errgroup.
func (s *Service) run(ctx context.Context, req domain.StartRequest, sc domain.Scenario) {
	g, gCtx := errgroup.WithContext(ctx)

	for _, vehicleID := range req.VehicleIDs {
		vid := vehicleID
		g.Go(func() error {
			return s.runVehicle(gCtx, vid, req, sc)
		})
	}

	if err := g.Wait(); err != nil && ctx.Err() == nil {
		s.logger.Errorw(ctx, "simulation ended with error", "error", err)
	}

	s.mu.Lock()
	s.status = domain.SimulationStatus{Running: false}
	s.mu.Unlock()
}

// runVehicle cycles through scenario waypoints, sending telemetry for each.
func (s *Service) runVehicle(ctx context.Context, vehicleID string, req domain.StartRequest, sc domain.Scenario) error {
	interval := time.Duration(req.IntervalMs) * time.Millisecond
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		for _, point := range sc.Points {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			lat := point.Latitude
			lon := point.Longitude

			if rng.Float64() < req.InvalidRate {
				lat = 999.0
			}

			sendReq := resources.TelemetrySendRequest{
				VehicleID:       vehicleID,
				Latitude:        lat,
				Longitude:       lon,
				DeviceTimestamp: time.Now().UTC(),
			}

			if err := s.sender.Send(ctx, sendReq); err != nil {
				s.logger.Errorw(ctx, "failed to send telemetry", "vehicle_id", vehicleID, "error", err)
			}

			if rng.Float64() < req.DuplicateRate {
				if err := s.sender.Send(ctx, sendReq); err != nil {
					s.logger.Errorw(ctx, "failed to send duplicate telemetry", "vehicle_id", vehicleID, "error", err)
				}
			}

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(interval):
			}
		}
	}
}

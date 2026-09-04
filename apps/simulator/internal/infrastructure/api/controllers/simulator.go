package controllers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"fleet/simulator/internal/core/domain"
	"fleet/simulator/internal/core/ports/services"
	"fleet/simulator/internal/infrastructure/api/dtos"
	apierrors "fleet/simulator/internal/infrastructure/api/errors"
	"fleet/shared/pkg/logger"
)

// Simulator is the HTTP controller for simulator endpoints.
type Simulator struct {
	logger  logger.Logger
	service services.SimulatorService
}

// GetScenarios handles GET /simulator/scenarios requests.
func (c *Simulator) GetScenarios(w http.ResponseWriter, r *http.Request) {
	scenarios := c.service.GetScenarios()

	items := make([]dtos.ScenarioResponse, len(scenarios))
	for i, s := range scenarios {
		items[i] = dtos.ScenarioResponse{
			Name:        s.Name,
			Description: s.Description,
			PointCount:  len(s.Points),
		}
	}

	resp := dtos.ScenariosResponse{Scenarios: items}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode scenarios response", "error", encErr)
	}
}

// GetStatus handles GET /simulator/status requests.
func (c *Simulator) GetStatus(w http.ResponseWriter, r *http.Request) {
	st := c.service.Status()

	resp := dtos.SimulationStatusResponse{
		Running:       st.Running,
		Scenario:      st.Scenario,
		VehicleCount:  st.VehicleCount,
		IntervalMs:    st.IntervalMs,
		DuplicateRate: st.DuplicateRate,
		InvalidRate:   st.InvalidRate,
	}
	if st.StartedAt != nil {
		resp.StartedAt = st.StartedAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode status response", "error", encErr)
	}
}

// Start handles POST /simulator/start requests.
func (c *Simulator) Start(w http.ResponseWriter, r *http.Request) {
	var req dtos.StartSimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	startReq := domain.StartRequest{
		VehicleIDs:    req.VehicleIDs,
		ScenarioName:  req.ScenarioName,
		IntervalMs:    req.IntervalMs,
		DuplicateRate: req.DuplicateRate,
		InvalidRate:   req.InvalidRate,
	}

	if err := c.service.Start(r.Context(), startReq); err != nil {
		c.logger.Errorw(r.Context(), "failed to start simulation", "error", err)
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// StartScenario handles POST /simulator/scenarios/:scenario/start requests.
func (c *Simulator) StartScenario(w http.ResponseWriter, r *http.Request) {
	scenarioName := chi.URLParam(r, "scenario")

	var req dtos.StartSimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}
	req.ScenarioName = scenarioName

	startReq := domain.StartRequest{
		VehicleIDs:    req.VehicleIDs,
		ScenarioName:  req.ScenarioName,
		IntervalMs:    req.IntervalMs,
		DuplicateRate: req.DuplicateRate,
		InvalidRate:   req.InvalidRate,
	}

	if err := c.service.Start(r.Context(), startReq); err != nil {
		c.logger.Errorw(r.Context(), "failed to start scenario simulation", "scenario", scenarioName, "error", err)
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Stop handles POST /simulator/stop requests.
func (c *Simulator) Stop(w http.ResponseWriter, r *http.Request) {
	c.service.Stop()
	w.WriteHeader(http.StatusNoContent)
}

// NewSimulator creates and returns a new Simulator controller.
func NewSimulator(log logger.Logger, svc services.SimulatorService) *Simulator {
	return &Simulator{logger: log, service: svc}
}

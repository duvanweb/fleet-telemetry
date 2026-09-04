package controllers

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"fleet/shared/pkg/logger"
	"fleet/simulator/internal/core/ports/services"
	"fleet/simulator/internal/infrastructure/api/dtos"
	apierrors "fleet/simulator/internal/infrastructure/api/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Health is the HTTP controller for health-related endpoints.
type Health struct {
	logger  logger.Logger
	service services.HealthService
}

// GetHealth handles GET /health and returns liveness status.
func (c *Health) GetHealth(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetHealth(r.Context())
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to get health status", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(dtos.HealthResponse{Status: result.Status}); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode health response", "error", encErr)
	}
}

// GetReady handles GET /ready and returns readiness status.
func (c *Health) GetReady(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetReady(r.Context())
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to get ready status", "error", err)
		apierrors.WriteError(w, http.StatusServiceUnavailable, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(dtos.HealthResponse{Status: result.Status}); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode ready response", "error", encErr)
	}
}

// NewHealth creates and returns a new Health controller.
func NewHealth(log logger.Logger, svc services.HealthService) *Health {
	return &Health{logger: log, service: svc}
}

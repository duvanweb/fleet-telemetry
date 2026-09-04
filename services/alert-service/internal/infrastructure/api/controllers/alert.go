package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"fleet/alert-service/internal/core/domain"
	"fleet/alert-service/internal/core/ports/services"
	"fleet/alert-service/internal/infrastructure/api/dtos"
	apierrors "fleet/alert-service/internal/infrastructure/api/errors"
	"fleet/shared/pkg/logger"
)

// Alert is the HTTP controller for alert-related endpoints.
type Alert struct {
	logger  logger.Logger
	service services.AlertService
}

// @Router /alerts [get]
// @Tags alerts
// @Summary Get all alerts.
// @Success 200 {object} dtos.AlertListResponse "Alerts retrieved successfully."
// @Failure 500 "Unexpected error."
// GetAll handles GET /alerts and returns all alerts.
func (c *Alert) GetAll(w http.ResponseWriter, r *http.Request) {
	alerts, err := c.service.GetAll(r.Context())
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to get all alerts", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := dtos.AlertListResponse{Alerts: make([]dtos.AlertResponse, 0, len(alerts))}
	for _, a := range alerts {
		response.Alerts = append(response.Alerts, toAlertResponse(a))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(response); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode alerts response", "error", encErr)
	}
}

// @Router /vehicles/{vehicleId}/alerts [get]
// @Tags alerts
// @Summary Get open alert for a vehicle.
// @Param vehicleId path string true "Vehicle ID."
// @Success 200 {object} dtos.AlertResponse "Open alert found."
// @Failure 404 "No open alert for vehicle."
// @Failure 500 "Unexpected error."
// GetOpenByVehicle handles GET /vehicles/:vehicleId/alerts and returns the open alert, if any.
func (c *Alert) GetOpenByVehicle(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "vehicleId")

	alert, found, err := c.service.GetOpenByVehicle(r.Context(), vehicleID)
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to get open alert", "vehicle_id", vehicleID, "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !found {
		apierrors.WriteError(w, http.StatusNotFound, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(toAlertResponse(alert)); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode open alert response", "error", encErr)
	}
}

// NewAlert creates and returns a new Alert controller.
func NewAlert(log logger.Logger, svc services.AlertService) *Alert {
	return &Alert{logger: log, service: svc}
}

func toAlertResponse(a domain.Alert) dtos.AlertResponse {
	return dtos.AlertResponse{
		ID:         a.ID,
		VehicleID:  a.VehicleID,
		Type:       a.Type,
		Status:     a.Status,
		StartedAt:  a.StartedAt,
		ResolvedAt: a.ResolvedAt,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"fleet/telemetry-service/internal/core/domain"
	"fleet/telemetry-service/internal/core/ports/services"
	"fleet/telemetry-service/internal/infrastructure/api/dtos"
	apierrors "fleet/telemetry-service/internal/infrastructure/api/errors"
	"fleet/shared/pkg/logger"
)

// Telemetry is the HTTP controller for telemetry-related endpoints.
type Telemetry struct {
	logger  logger.Logger
	service services.TelemetryService
}

// GetByVehicleID handles GET /vehicles/:id/telemetry and returns telemetry points for a vehicle.
func (c *Telemetry) GetByVehicleID(w http.ResponseWriter, r *http.Request) {
	vehicleID := chi.URLParam(r, "id")

	points, err := c.service.GetByVehicleID(r.Context(), vehicleID)
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to get telemetry", "vehicle_id", vehicleID, "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := dtos.VehicleTelemetryResponse{
		VehicleID: vehicleID,
		Points:    make([]dtos.TelemetryPointResponse, 0, len(points)),
	}
	for _, p := range points {
		resp.Points = append(resp.Points, toTelemetryPointResponse(p))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode telemetry response", "error", encErr)
	}
}

// IngestTelemetry handles POST /telemetry and persists a new GPS point.
func (c *Telemetry) IngestTelemetry(w http.ResponseWriter, r *http.Request) {
	var req dtos.IngestTelemetryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	deviceTimestamp, err := time.Parse(time.RFC3339Nano, req.DeviceTimestamp)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, errors.New("device_timestamp must be RFC3339Nano"))
		return
	}

	point := domain.TelemetryPoint{
		VehicleID:       req.VehicleID,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		DeviceTimestamp: deviceTimestamp,
	}

	saved, err := c.service.IngestTelemetry(r.Context(), point)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, domain.ErrInvalidLatitude),
			errors.Is(err, domain.ErrInvalidLongitude),
			errors.Is(err, domain.ErrInvalidTimestamp),
			errors.Is(err, domain.ErrFutureTelemetry):
			status = http.StatusBadRequest
		case errors.Is(err, domain.ErrVehicleNotFound):
			status = http.StatusNotFound
		case errors.Is(err, domain.ErrDuplicateTelemetry):
			status = http.StatusConflict
		}
		c.logger.Errorw(r.Context(), "failed to ingest telemetry", "vehicle_id", req.VehicleID, "error", err)
		apierrors.WriteError(w, status, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if encErr := json.NewEncoder(w).Encode(toTelemetryPointResponse(saved)); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode ingest response", "error", encErr)
	}
}

// NewTelemetry creates and returns a new Telemetry controller.
func NewTelemetry(log logger.Logger, svc services.TelemetryService) *Telemetry {
	return &Telemetry{logger: log, service: svc}
}

func toTelemetryPointResponse(p domain.TelemetryPoint) dtos.TelemetryPointResponse {
	return dtos.TelemetryPointResponse{
		ID:               p.ID,
		VehicleID:        p.VehicleID,
		Latitude:         p.Latitude,
		Longitude:        p.Longitude,
		DeviceTimestamp:  p.DeviceTimestamp.UTC().Format(time.RFC3339Nano),
		ReceivedAt:       p.ReceivedAt.UTC().Format(time.RFC3339Nano),
		DeduplicationKey: p.DeduplicationKey,
	}
}

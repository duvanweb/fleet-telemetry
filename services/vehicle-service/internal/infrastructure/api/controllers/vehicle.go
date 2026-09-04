package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"fleet/vehicle-service/internal/core/domain"
	"fleet/vehicle-service/internal/core/ports/services"
	"fleet/vehicle-service/internal/infrastructure/api/dtos"
	apierrors "fleet/vehicle-service/internal/infrastructure/api/errors"
	"fleet/shared/pkg/logger"
)

// Vehicle is the HTTP controller for vehicle-related endpoints.
type Vehicle struct {
	logger  logger.Logger
	service services.VehicleService
}

// Create handles POST /vehicles and persists a new vehicle.
func (c *Vehicle) Create(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	v := domain.Vehicle{
		ExternalID: req.ExternalID,
		Plate:      req.Plate,
		Name:       req.Name,
	}

	created, err := c.service.Create(r.Context(), v)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicatePlate) {
			apierrors.WriteError(w, http.StatusConflict, err)
			return
		}
		c.logger.Errorw(r.Context(), "failed to create vehicle", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if encErr := json.NewEncoder(w).Encode(toVehicleResponse(created)); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode create vehicle response", "error", encErr)
	}
}

// Delete handles DELETE /vehicles/:id and soft-deletes the vehicle.
func (c *Vehicle) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := c.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrVehicleNotFound) {
			apierrors.WriteError(w, http.StatusNotFound, err)
			return
		}
		c.logger.Errorw(r.Context(), "failed to delete vehicle", "vehicle_id", id, "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAll handles GET /vehicles and returns all non-deleted vehicles.
func (c *Vehicle) GetAll(w http.ResponseWriter, r *http.Request) {
	vehicles, err := c.service.GetAll(r.Context())
	if err != nil {
		c.logger.Errorw(r.Context(), "failed to get all vehicles", "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := make([]dtos.VehicleResponse, len(vehicles))
	for i, v := range vehicles {
		resp[i] = toVehicleResponse(v)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode vehicles response", "error", encErr)
	}
}

// GetByID handles GET /vehicles/:id and returns a single vehicle.
func (c *Vehicle) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	v, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrVehicleNotFound) || errors.Is(err, domain.ErrVehicleDeleted) {
			apierrors.WriteError(w, http.StatusNotFound, err)
			return
		}
		c.logger.Errorw(r.Context(), "failed to get vehicle", "vehicle_id", id, "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(toVehicleResponse(v)); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode vehicle response", "error", encErr)
	}
}

// Update handles PATCH /vehicles/:id and updates plate and name.
func (c *Vehicle) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dtos.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, err)
		return
	}

	v := domain.Vehicle{
		ID:    id,
		Plate: req.Plate,
		Name:  req.Name,
	}

	updated, err := c.service.Update(r.Context(), v)
	if err != nil {
		if errors.Is(err, domain.ErrVehicleNotFound) {
			apierrors.WriteError(w, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, domain.ErrDuplicatePlate) {
			apierrors.WriteError(w, http.StatusConflict, err)
			return
		}
		c.logger.Errorw(r.Context(), "failed to update vehicle", "vehicle_id", id, "error", err)
		apierrors.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encErr := json.NewEncoder(w).Encode(toVehicleResponse(updated)); encErr != nil {
		c.logger.Errorw(r.Context(), "failed to encode update vehicle response", "error", encErr)
	}
}

// NewVehicle creates and returns a new Vehicle controller.
func NewVehicle(log logger.Logger, svc services.VehicleService) *Vehicle {
	return &Vehicle{logger: log, service: svc}
}

func toVehicleResponse(v domain.Vehicle) dtos.VehicleResponse {
	resp := dtos.VehicleResponse{
		ID:         v.ID,
		ExternalID: v.ExternalID,
		Plate:      v.Plate,
		Name:       v.Name,
		CreatedAt:  v.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  v.UpdatedAt.Format(time.RFC3339),
	}
	if v.DeletedAt != nil {
		s := v.DeletedAt.Format(time.RFC3339)
		resp.DeletedAt = &s
	}
	return resp
}

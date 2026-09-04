package vehicle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"fleet/shared/pkg/logger"
	"fleet/vehicle-service/internal/core/domain"
	"fleet/vehicle-service/internal/core/ports/repositories"
	"fleet/vehicle-service/internal/infrastructure/postgres/models"
	vehiclesql "fleet/vehicle-service/internal/infrastructure/postgres/repositories/vehicle/sql"
)

// Dependencies holds the repository's injected dependencies.
type Dependencies struct {
	DB repositories.Databaser
}

type repository struct {
	logger logger.Logger
	db     repositories.Databaser
}

// NewRepository creates and returns a new vehicle Repository.
func NewRepository(log logger.Logger, deps Dependencies) repositories.VehicleRepository {
	return &repository{logger: log, db: deps.DB}
}

// Create inserts a new vehicle row and returns the persisted entity.
func (r *repository) Create(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
	var row models.VehicleRow

	err := r.db.QueryRowContext(ctx, vehiclesql.InsertVehicle,
		v.ID, v.ExternalID, v.Plate, v.Name, v.CreatedAt, v.UpdatedAt,
	).Scan(&row.ID, &row.ExternalID, &row.Plate, &row.Name, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.Vehicle{}, domain.ErrDuplicatePlate
		}
		return domain.Vehicle{}, fmt.Errorf("inserting vehicle: %w", err)
	}

	return row.ToDomain(), nil
}

// ExistsByExternalID reports whether a non-deleted vehicle with the given external ID exists.
func (r *repository) ExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, vehiclesql.ExistsByExternalID, externalID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking external id existence: %w", err)
	}
	return exists, nil
}

// ExistsByPlate reports whether a non-deleted vehicle with the given plate exists.
func (r *repository) ExistsByPlate(ctx context.Context, plate string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, vehiclesql.ExistsByPlate, plate).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking plate existence: %w", err)
	}
	return exists, nil
}

// GetAll returns all non-deleted vehicles ordered by creation date descending.
func (r *repository) GetAll(ctx context.Context) ([]domain.Vehicle, error) {
	rows, err := r.db.QueryContext(ctx, vehiclesql.SelectAll)
	if err != nil {
		return nil, fmt.Errorf("querying vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []domain.Vehicle
	for rows.Next() {
		var row models.VehicleRow
		if err := rows.Scan(&row.ID, &row.ExternalID, &row.Plate, &row.Name, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt); err != nil {
			return nil, fmt.Errorf("scanning vehicle row: %w", err)
		}
		vehicles = append(vehicles, row.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating vehicle rows: %w", err)
	}

	return vehicles, nil
}

// GetByID returns the vehicle with the given ID regardless of deletion status.
func (r *repository) GetByID(ctx context.Context, id string) (domain.Vehicle, error) {
	var row models.VehicleRow

	err := r.db.QueryRowContext(ctx, vehiclesql.SelectByID, id).
		Scan(&row.ID, &row.ExternalID, &row.Plate, &row.Name, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Vehicle{}, domain.ErrVehicleNotFound
	}
	if err != nil {
		return domain.Vehicle{}, fmt.Errorf("querying vehicle by id: %w", err)
	}

	return row.ToDomain(), nil
}

// SoftDelete marks the vehicle with the given ID as deleted.
func (r *repository) SoftDelete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, vehiclesql.SoftDelete, id)
	if err != nil {
		return fmt.Errorf("soft deleting vehicle: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("reading rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrVehicleNotFound
	}

	return nil
}

// Update persists changes to plate and name for the given vehicle.
func (r *repository) Update(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
	var row models.VehicleRow

	err := r.db.QueryRowContext(ctx, vehiclesql.UpdateVehicle, v.Plate, v.Name, v.ID).
		Scan(&row.ID, &row.ExternalID, &row.Plate, &row.Name, &row.CreatedAt, &row.UpdatedAt, &row.DeletedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Vehicle{}, domain.ErrVehicleNotFound
	}
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.Vehicle{}, domain.ErrDuplicatePlate
		}
		return domain.Vehicle{}, fmt.Errorf("updating vehicle: %w", err)
	}

	return row.ToDomain(), nil
}

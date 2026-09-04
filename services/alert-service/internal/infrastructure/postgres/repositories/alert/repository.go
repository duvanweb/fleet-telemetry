package alert

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/fx"

	"fleet/alert-service/internal/core/domain"
	"fleet/alert-service/internal/core/ports/repositories"
	"fleet/alert-service/internal/infrastructure/postgres/models"
	alertsql "fleet/alert-service/internal/infrastructure/postgres/repositories/alert/sql"
	"fleet/shared/pkg/logger"
)

// Dependencies holds injected dependencies for the alert Repository.
type Dependencies struct {
	fx.In

	DB repositories.Databaser
}

// Repository implements repositories.AlertRepository using PostgreSQL.
type Repository struct {
	logger logger.Logger
	db     repositories.Databaser
}

// NewRepository creates and returns a new alert Repository.
func NewRepository(log logger.Logger, deps Dependencies) *Repository {
	return &Repository{logger: log, db: deps.DB}
}

// Create persists a new alert record and returns the saved entity.
func (r *Repository) Create(ctx context.Context, alert domain.Alert) (domain.Alert, error) {
	row := r.db.QueryRowContext(ctx, alertsql.InsertAlert,
		alert.ID,
		alert.VehicleID,
		alert.Type,
		alert.Status,
		alert.StartedAt,
		alert.ResolvedAt,
		nil,
		alert.CreatedAt,
		alert.UpdatedAt,
	)

	var m models.AlertRow
	if err := row.Scan(
		&m.ID, &m.VehicleID, &m.Type, &m.Status,
		&m.StartedAt, &m.ResolvedAt, &m.Metadata,
		&m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		return domain.Alert{}, fmt.Errorf("scanning created alert: %w", err)
	}

	return m.ToDomain(), nil
}

// GetAll retrieves all alerts ordered by creation time descending.
func (r *Repository) GetAll(ctx context.Context) ([]domain.Alert, error) {
	rows, err := r.db.QueryContext(ctx, alertsql.SelectAllAlerts)
	if err != nil {
		return nil, fmt.Errorf("querying all alerts: %w", err)
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var m models.AlertRow
		if err := rows.Scan(
			&m.ID, &m.VehicleID, &m.Type, &m.Status,
			&m.StartedAt, &m.ResolvedAt, &m.Metadata,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning alert row: %w", err)
		}
		alerts = append(alerts, m.ToDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating alert rows: %w", err)
	}

	return alerts, nil
}

// GetOpenByVehicle returns the open alert for the given vehicle ID, if any.
func (r *Repository) GetOpenByVehicle(ctx context.Context, vehicleID string) (domain.Alert, bool, error) {
	row := r.db.QueryRowContext(ctx, alertsql.SelectOpenByVehicle, vehicleID)

	var m models.AlertRow
	if err := row.Scan(
		&m.ID, &m.VehicleID, &m.Type, &m.Status,
		&m.StartedAt, &m.ResolvedAt, &m.Metadata,
		&m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Alert{}, false, nil
		}
		return domain.Alert{}, false, fmt.Errorf("querying open alert by vehicle: %w", err)
	}

	return m.ToDomain(), true, nil
}

// MarkResolved sets the status of the given alert to RESOLVED.
func (r *Repository) MarkResolved(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, alertsql.UpdateMarkResolved, id)
	if err != nil {
		return fmt.Errorf("marking alert resolved: %w", err)
	}

	return nil
}

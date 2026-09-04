package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"fleet/shared/pkg/logger"
	"fleet/telemetry-service/internal/core/domain"
	"fleet/telemetry-service/internal/core/ports/repositories"
	sqlqueries "fleet/telemetry-service/internal/infrastructure/postgres/repositories/telemetry/sql"
)

// Dependencies holds the repository's injected dependencies.
type Dependencies struct {
	DB repositories.Databaser
}

// Repository implements TelemetryRepository using PostgreSQL.
type Repository struct {
	logger logger.Logger
	db     repositories.Databaser
}

// NewRepository creates and returns a new telemetry Repository.
func NewRepository(log logger.Logger, deps Dependencies) *Repository {
	return &Repository{logger: log, db: deps.DB}
}

// Create persists a new telemetry point and returns it.
func (r *Repository) Create(ctx context.Context, point domain.TelemetryPoint) (domain.TelemetryPoint, error) {
	var saved domain.TelemetryPoint
	err := r.db.QueryRowContext(ctx, sqlqueries.InsertTelemetry,
		point.ID,
		point.VehicleID,
		point.Latitude,
		point.Longitude,
		point.DeviceTimestamp,
		point.ReceivedAt,
		point.DeduplicationKey,
	).Scan(
		&saved.ID,
		&saved.VehicleID,
		&saved.Latitude,
		&saved.Longitude,
		&saved.DeviceTimestamp,
		&saved.ReceivedAt,
		&saved.DeduplicationKey,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.TelemetryPoint{}, domain.ErrDuplicateTelemetry
		}
		return domain.TelemetryPoint{}, fmt.Errorf("inserting telemetry point: %w", err)
	}

	return saved, nil
}

// CreateWithOutbox persists a telemetry point and an outbox event atomically in a single transaction.
func (r *Repository) CreateWithOutbox(ctx context.Context, point domain.TelemetryPoint, event domain.OutboxEvent) (domain.TelemetryPoint, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.TelemetryPoint{}, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var saved domain.TelemetryPoint
	err = tx.QueryRowContext(ctx, sqlqueries.InsertTelemetry,
		point.ID,
		point.VehicleID,
		point.Latitude,
		point.Longitude,
		point.DeviceTimestamp,
		point.ReceivedAt,
		point.DeduplicationKey,
	).Scan(
		&saved.ID,
		&saved.VehicleID,
		&saved.Latitude,
		&saved.Longitude,
		&saved.DeviceTimestamp,
		&saved.ReceivedAt,
		&saved.DeduplicationKey,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.TelemetryPoint{}, domain.ErrDuplicateTelemetry
		}
		return domain.TelemetryPoint{}, fmt.Errorf("inserting telemetry point: %w", err)
	}

	_, err = tx.ExecContext(ctx, sqlqueries.InsertOutboxEventInTx,
		event.ID,
		event.EventType,
		event.Payload,
		event.Status,
		event.CreatedAt,
	)
	if err != nil {
		return domain.TelemetryPoint{}, fmt.Errorf("inserting outbox event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return domain.TelemetryPoint{}, fmt.Errorf("committing transaction: %w", err)
	}

	return saved, nil
}

// GetByVehicleID returns all telemetry points for the given vehicle ordered by device_timestamp ASC.
func (r *Repository) GetByVehicleID(ctx context.Context, vehicleID string) ([]domain.TelemetryPoint, error) {
	rows, err := r.db.QueryContext(ctx, sqlqueries.SelectByVehicleID, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("querying telemetry by vehicle id: %w", err)
	}
	defer rows.Close()

	var points []domain.TelemetryPoint
	for rows.Next() {
		var point domain.TelemetryPoint
		if err := rows.Scan(
			&point.ID,
			&point.VehicleID,
			&point.Latitude,
			&point.Longitude,
			&point.DeviceTimestamp,
			&point.ReceivedAt,
			&point.DeduplicationKey,
		); err != nil {
			return nil, fmt.Errorf("scanning telemetry row: %w", err)
		}
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating telemetry rows: %w", err)
	}

	return points, nil
}

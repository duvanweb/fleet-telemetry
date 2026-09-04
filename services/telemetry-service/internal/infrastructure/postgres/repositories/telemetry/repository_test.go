package telemetry_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"fleet/telemetry-service/internal/core/domain"
	"fleet/telemetry-service/internal/infrastructure/postgres/repositories/telemetry"
	sqlqueries "fleet/telemetry-service/internal/infrastructure/postgres/repositories/telemetry/sql"
	testdata "fleet/telemetry-service/test/data"
	"fleet/shared/pkg/logger"
)

var errTest = errors.New("unexpected error")

func newRepo(t *testing.T, db *sql.DB) *telemetry.Repository {
	t.Helper()
	log, _ := logger.NewLogger()
	return telemetry.NewRepository(log, telemetry.Dependencies{DB: db})
}

var telemetryCols = []string{
	"id", "vehicle_id", "latitude", "longitude",
	"device_timestamp", "received_at", "deduplication_key",
}

func TestRepository_GetByVehicleID(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestTelemetryPoint()

	tests := []struct {
		name          string
		setup         func(m sqlmock.Sqlmock)
		expected      []domain.TelemetryPoint
		expectedError error
	}{
		{
			name: "works correctly",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(telemetryCols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Latitude, fixture.Longitude,
					fixture.DeviceTimestamp, fixture.ReceivedAt, fixture.DeduplicationKey,
				)
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.SelectByVehicleID)).
					WithArgs(fixture.VehicleID).
					WillReturnRows(rows)
			},
			expected: []domain.TelemetryPoint{fixture},
		},
		{
			name: "works correctly when no results",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.SelectByVehicleID)).
					WithArgs(fixture.VehicleID).
					WillReturnRows(sqlmock.NewRows(telemetryCols))
			},
			expected: nil,
		},
		{
			name: "fails when query fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.SelectByVehicleID)).
					WithArgs(fixture.VehicleID).
					WillReturnError(errTest)
			},
			expectedError: errors.New("querying telemetry by vehicle id: unexpected error"),
		},
		{
			name: "fails when scan fails",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(telemetryCols).AddRow(
					fixture.ID, fixture.VehicleID, "not-a-float", fixture.Longitude,
					fixture.DeviceTimestamp, fixture.ReceivedAt, fixture.DeduplicationKey,
				)
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.SelectByVehicleID)).
					WithArgs(fixture.VehicleID).
					WillReturnRows(rows)
			},
			expectedError: errors.New("scanning telemetry row"),
		},
		{
			name: "fails when rows.Err fails",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(telemetryCols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Latitude, fixture.Longitude,
					fixture.DeviceTimestamp, fixture.ReceivedAt, fixture.DeduplicationKey,
				).RowError(0, errTest)
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.SelectByVehicleID)).
					WithArgs(fixture.VehicleID).
					WillReturnRows(rows)
			},
			expectedError: errors.New("iterating telemetry rows: unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setup(mock)

			repo := newRepo(t, db)
			result, err := repo.GetByVehicleID(context.Background(), fixture.VehicleID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error()[:20])
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateWithOutbox(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestTelemetryPoint()

	tests := []struct {
		name          string
		setup         func(m sqlmock.Sqlmock)
		expected      domain.TelemetryPoint
		expectedError bool
	}{
		{
			name: "works correctly",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				rows := sqlmock.NewRows(telemetryCols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Latitude, fixture.Longitude,
					fixture.DeviceTimestamp, fixture.ReceivedAt, fixture.DeduplicationKey,
				)
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.InsertTelemetry)).WillReturnRows(rows)
				m.ExpectExec(regexp.QuoteMeta(sqlqueries.InsertOutboxEventInTx)).WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
			expected: fixture,
		},
		{
			name: "fails when begin tx fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(errTest)
			},
			expectedError: true,
		},
		{
			name: "fails when insert telemetry fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.InsertTelemetry)).WillReturnError(errTest)
				m.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name: "fails when insert outbox fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				rows := sqlmock.NewRows(telemetryCols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Latitude, fixture.Longitude,
					fixture.DeviceTimestamp, fixture.ReceivedAt, fixture.DeduplicationKey,
				)
				m.ExpectQuery(regexp.QuoteMeta(sqlqueries.InsertTelemetry)).WillReturnRows(rows)
				m.ExpectExec(regexp.QuoteMeta(sqlqueries.InsertOutboxEventInTx)).WillReturnError(errTest)
				m.ExpectRollback()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setup(mock)

			repo := newRepo(t, db)
			result, err := repo.CreateWithOutbox(context.Background(), fixture, domain.OutboxEvent{})

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

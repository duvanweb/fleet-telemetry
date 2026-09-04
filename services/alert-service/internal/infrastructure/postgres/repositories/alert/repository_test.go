package alert_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"fleet/alert-service/internal/core/domain"
	"fleet/alert-service/internal/core/ports/repositories"
	"fleet/alert-service/internal/infrastructure/postgres/repositories/alert"
	alertsql "fleet/alert-service/internal/infrastructure/postgres/repositories/alert/sql"
	testdata "fleet/alert-service/test/data"
	"fleet/shared/pkg/logger"
)

var errTest = errors.New("unexpected error")

func newRepo(t *testing.T, db *sql.DB) *alert.Repository {
	t.Helper()
	log, _ := logger.NewLogger()
	return alert.NewRepository(log, alert.Dependencies{DB: db})
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		expected      domain.Alert
		expectedError error
	}{
		{
			name: "works correctly",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "type", "status",
					"started_at", "resolved_at", "metadata",
					"created_at", "updated_at",
				}).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Type, fixture.Status,
					fixture.StartedAt, nil, nil,
					fixture.CreatedAt, fixture.UpdatedAt,
				)
				m.ExpectQuery(regexp.QuoteMeta(alertsql.InsertAlert)).
					WithArgs(
						fixture.ID, fixture.VehicleID, fixture.Type, fixture.Status,
						fixture.StartedAt, nil, nil,
						fixture.CreatedAt, fixture.UpdatedAt,
					).
					WillReturnRows(rows)
			},
			expected: fixture,
		},
		{
			name: "fails when query fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(alertsql.InsertAlert)).
					WillReturnError(errTest)
			},
			expectedError: errors.New("scanning created alert: sql: expected 0 destination arguments in Scan, not 9"),
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
			result, err := repo.Create(context.Background(), fixture)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetAll(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	cols := []string{
		"id", "vehicle_id", "type", "status",
		"started_at", "resolved_at", "metadata",
		"created_at", "updated_at",
	}

	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		expected      []domain.Alert
		expectedError error
	}{
		{
			name: "works correctly",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Type, fixture.Status,
					fixture.StartedAt, nil, nil,
					fixture.CreatedAt, fixture.UpdatedAt,
				)
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectAllAlerts)).WillReturnRows(rows)
			},
			expected: []domain.Alert{fixture},
		},
		{
			name: "works correctly when no results",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectAllAlerts)).
					WillReturnRows(sqlmock.NewRows(cols))
			},
			expected: nil,
		},
		{
			name: "fails when query fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectAllAlerts)).WillReturnError(errTest)
			},
			expectedError: errors.New("querying all alerts: unexpected error"),
		},
		{
			name: "fails when scan fails",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Type, fixture.Status,
					"not-a-time", nil, nil,
					fixture.CreatedAt, fixture.UpdatedAt,
				)
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectAllAlerts)).WillReturnRows(rows)
			},
			expectedError: errors.New("scanning alert row: sql: Scan error on column index 4"),
		},
		{
			name: "fails when rows.Err fails",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Type, fixture.Status,
					fixture.StartedAt, nil, nil,
					fixture.CreatedAt, fixture.UpdatedAt,
				).RowError(0, errTest)
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectAllAlerts)).WillReturnRows(rows)
			},
			expectedError: errors.New("iterating alert rows: unexpected error"),
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
			result, err := repo.GetAll(context.Background())

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error()[:30])
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetOpenByVehicle(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	cols := []string{
		"id", "vehicle_id", "type", "status",
		"started_at", "resolved_at", "metadata",
		"created_at", "updated_at",
	}

	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		expected      domain.Alert
		expectedFound bool
		expectedError error
	}{
		{
			name: "works correctly when alert found",
			setup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(cols).AddRow(
					fixture.ID, fixture.VehicleID, fixture.Type, fixture.Status,
					fixture.StartedAt, nil, nil,
					fixture.CreatedAt, fixture.UpdatedAt,
				)
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectOpenByVehicle)).
					WithArgs(fixture.VehicleID).
					WillReturnRows(rows)
			},
			expected:      fixture,
			expectedFound: true,
		},
		{
			name: "works correctly when no alert found",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectOpenByVehicle)).
					WithArgs(fixture.VehicleID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
		},
		{
			name: "fails when query fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(alertsql.SelectOpenByVehicle)).
					WithArgs(fixture.VehicleID).
					WillReturnError(errTest)
			},
			expectedError: errors.New("querying open alert by vehicle: unexpected error"),
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
			result, found, err := repo.GetOpenByVehicle(context.Background(), fixture.VehicleID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, "querying open alert by vehicle")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFound, found)
				if tt.expectedFound {
					assert.Equal(t, tt.expected, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_MarkResolved(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "works correctly",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(alertsql.UpdateMarkResolved)).
					WithArgs(fixture.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "fails when exec fails",
			setup: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(alertsql.UpdateMarkResolved)).
					WithArgs(fixture.ID).
					WillReturnError(errTest)
			},
			expectedError: errors.New("marking alert resolved: unexpected error"),
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
			err = repo.MarkResolved(context.Background(), fixture.ID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, "marking alert resolved")
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Ensure Repository implements repositories.AlertRepository.
var _ repositories.AlertRepository = (*alert.Repository)(nil)

// unused but kept for type safety check
var _ = time.Now

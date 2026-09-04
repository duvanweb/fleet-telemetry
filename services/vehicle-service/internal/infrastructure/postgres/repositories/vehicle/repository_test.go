package vehicle_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"fleet/vehicle-service/internal/core/domain"
	"fleet/vehicle-service/internal/core/ports/repositories"
	vehiclerepo "fleet/vehicle-service/internal/infrastructure/postgres/repositories/vehicle"
	vehiclesql "fleet/vehicle-service/internal/infrastructure/postgres/repositories/vehicle/sql"
	"fleet/vehicle-service/test/data"
	"fleet/shared/pkg/logger"
)

var errTest = errors.New("unexpected db error")

func newTestRepo(t *testing.T) (repositories.VehicleRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	log, _ := logger.NewLogger()
	repo := vehiclerepo.NewRepository(log, vehiclerepo.Dependencies{DB: db})
	return repo, mock
}

func vehicleColumns() []string {
	return []string{"id", "external_id", "plate", "name", "created_at", "updated_at", "deleted_at"}
}

func vehicleRow(v domain.Vehicle) *sqlmock.Rows {
	return sqlmock.NewRows(vehicleColumns()).
		AddRow(v.ID, v.ExternalID, v.Plate, v.Name, v.CreatedAt, v.UpdatedAt, v.DeletedAt)
}

func TestVehicleRepository_Create(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		expectedVehicle domain.Vehicle
		expectedError error
	}{
		{
			name: "works correctly",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.InsertVehicle)).
					WithArgs(vehicle.ID, vehicle.ExternalID, vehicle.Plate, vehicle.Name, vehicle.CreatedAt, vehicle.UpdatedAt).
					WillReturnRows(vehicleRow(vehicle))
			},
			expectedVehicle: vehicle,
		},
		{
			name: "fails when plate is duplicate",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.InsertVehicle)).
					WithArgs(vehicle.ID, vehicle.ExternalID, vehicle.Plate, vehicle.Name, vehicle.CreatedAt, vehicle.UpdatedAt).
					WillReturnError(&pq.Error{Code: "23505"})
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrDuplicatePlate,
		},
		{
			name: "fails when db query fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.InsertVehicle)).
					WithArgs(vehicle.ID, vehicle.ExternalID, vehicle.Plate, vehicle.Name, vehicle.CreatedAt, vehicle.UpdatedAt).
					WillReturnError(errTest)
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			result, err := repo.Create(t.Context(), vehicle)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedVehicle, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicle, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVehicleRepository_ExistsByExternalID(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		expectedExists bool
		expectedError  error
	}{
		{
			name: "works correctly when vehicle exists",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.ExistsByExternalID)).
					WithArgs(vehicle.ExternalID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedExists: true,
		},
		{
			name: "works correctly when vehicle does not exist",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.ExistsByExternalID)).
					WithArgs(vehicle.ExternalID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedExists: false,
		},
		{
			name: "fails when db query fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.ExistsByExternalID)).
					WithArgs(vehicle.ExternalID).
					WillReturnError(errTest)
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			exists, err := repo.ExistsByExternalID(t.Context(), vehicle.ExternalID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVehicleRepository_ExistsByPlate(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name           string
		setup          func(mock sqlmock.Sqlmock)
		expectedExists bool
		expectedError  error
	}{
		{
			name: "works correctly when plate exists",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.ExistsByPlate)).
					WithArgs(vehicle.Plate).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expectedExists: true,
		},
		{
			name: "works correctly when plate does not exist",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.ExistsByPlate)).
					WithArgs(vehicle.Plate).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expectedExists: false,
		},
		{
			name: "fails when db query fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.ExistsByPlate)).
					WithArgs(vehicle.Plate).
					WillReturnError(errTest)
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			exists, err := repo.ExistsByPlate(t.Context(), vehicle.Plate)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVehicleRepository_GetAll(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name             string
		setup            func(mock sqlmock.Sqlmock)
		expectedVehicles []domain.Vehicle
		expectedError    error
	}{
		{
			name: "works correctly",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectAll)).
					WillReturnRows(vehicleRow(vehicle))
			},
			expectedVehicles: []domain.Vehicle{vehicle},
		},
		{
			name: "works correctly when no vehicles exist",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectAll)).
					WillReturnRows(sqlmock.NewRows(vehicleColumns()))
			},
			expectedVehicles: nil,
		},
		{
			name: "fails when db query fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectAll)).
					WillReturnError(errTest)
			},
			expectedError: errTest,
		},
		{
			name: "fails when scan fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectAll)).
					WillReturnRows(sqlmock.NewRows(vehicleColumns()).
						AddRow("id", "ext", "plate", "name", "bad-time", time.Now(), nil))
			},
			expectedError: errors.New("scanning vehicle row"),
		},
		{
			name: "fails when rows.Err fails",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows(vehicleColumns()).AddRow(
					vehicle.ID, vehicle.ExternalID, vehicle.Plate, vehicle.Name,
					vehicle.CreatedAt, vehicle.UpdatedAt, vehicle.DeletedAt,
				).RowError(0, errTest)
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectAll)).WillReturnRows(rows)
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			result, err := repo.GetAll(t.Context())

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicles, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVehicleRepository_GetByID(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name            string
		setup           func(mock sqlmock.Sqlmock)
		expectedVehicle domain.Vehicle
		expectedError   error
	}{
		{
			name: "works correctly",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectByID)).
					WithArgs(vehicle.ID).
					WillReturnRows(vehicleRow(vehicle))
			},
			expectedVehicle: vehicle,
		},
		{
			name: "fails when vehicle not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectByID)).
					WithArgs(vehicle.ID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrVehicleNotFound,
		},
		{
			name: "fails when db query fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.SelectByID)).
					WithArgs(vehicle.ID).
					WillReturnError(errTest)
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			result, err := repo.GetByID(t.Context(), vehicle.ID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedVehicle, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicle, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVehicleRepository_SoftDelete(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name          string
		setup         func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "works correctly",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(vehiclesql.SoftDelete)).
					WithArgs(vehicle.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "fails when vehicle not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(vehiclesql.SoftDelete)).
					WithArgs(vehicle.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: domain.ErrVehicleNotFound,
		},
		{
			name: "fails when db exec fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(vehiclesql.SoftDelete)).
					WithArgs(vehicle.ID).
					WillReturnError(errTest)
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			err := repo.SoftDelete(t.Context(), vehicle.ID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVehicleRepository_Update(t *testing.T) {
	t.Parallel()

	vehicle := data.GetTestVehicle()

	tests := []struct {
		name            string
		setup           func(mock sqlmock.Sqlmock)
		expectedVehicle domain.Vehicle
		expectedError   error
	}{
		{
			name: "works correctly",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.UpdateVehicle)).
					WithArgs(vehicle.Plate, vehicle.Name, vehicle.ID).
					WillReturnRows(vehicleRow(vehicle))
			},
			expectedVehicle: vehicle,
		},
		{
			name: "fails when vehicle not found",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.UpdateVehicle)).
					WithArgs(vehicle.Plate, vehicle.Name, vehicle.ID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrVehicleNotFound,
		},
		{
			name: "fails when plate is duplicate",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.UpdateVehicle)).
					WithArgs(vehicle.Plate, vehicle.Name, vehicle.ID).
					WillReturnError(&pq.Error{Code: "23505"})
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrDuplicatePlate,
		},
		{
			name: "fails when db query fails",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(vehiclesql.UpdateVehicle)).
					WithArgs(vehicle.Plate, vehicle.Name, vehicle.ID).
					WillReturnError(errTest)
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepo(t)
			tt.setup(mock)

			result, err := repo.Update(t.Context(), vehicle)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedVehicle, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicle, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

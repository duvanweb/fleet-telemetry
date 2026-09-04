package telemetry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fleet/telemetry-service/internal/core/domain"
	"fleet/telemetry-service/internal/core/telemetry"
	repomocks "fleet/telemetry-service/internal/core/ports/repositories/mocks"
	resmocks "fleet/telemetry-service/internal/core/ports/resources/mocks"
	testdata "fleet/telemetry-service/test/data"
	"fleet/shared/pkg/logger"
)

var (
	errTest = errors.New("unexpected error")
	anyctx  = mock.Anything
	anypoint = mock.Anything
)

type Dependencies struct {
	Telemetry      *repomocks.TelemetryRepository
	VehicleChecker *resmocks.VehicleChecker
}

func newService(t *testing.T, deps Dependencies) *telemetry.Service {
	t.Helper()
	log, _ := logger.NewLogger()
	return telemetry.NewService(log,
		telemetry.Repositories{Telemetry: deps.Telemetry},
		telemetry.Resources{VehicleChecker: deps.VehicleChecker},
	)
}

func TestService_GetByVehicleID(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestTelemetryPoint()
	repoSuccessMock := &repomocks.TelemetryRepository{}
	repoSuccessMock.On("GetByVehicleID", anyctx, fixture.VehicleID).Return([]domain.TelemetryPoint{fixture}, nil)

	tests := []struct {
		name          string
		vehicleID     string
		dependencies  Dependencies
		expected      []domain.TelemetryPoint
		expectedError error
	}{
		{
			name:         "works correctly",
			vehicleID:    fixture.VehicleID,
			dependencies: Dependencies{Telemetry: repoSuccessMock},
			expected:     []domain.TelemetryPoint{fixture},
		},
		{
			name:      "fails when repository GetByVehicleID fails",
			vehicleID: fixture.VehicleID,
			dependencies: Dependencies{
				Telemetry: func() *repomocks.TelemetryRepository {
					mockRepo := &repomocks.TelemetryRepository{}
					mockRepo.On("GetByVehicleID", anyctx, fixture.VehicleID).Return(nil, errTest)
					return mockRepo
				}(),
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newService(t, tt.dependencies)
			result, err := svc.GetByVehicleID(context.Background(), tt.vehicleID)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expected, result)

			if tt.dependencies.Telemetry != nil {
				tt.dependencies.Telemetry.AssertExpectations(t)
			}
		})
	}
}

func TestService_IngestTelemetry(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestTelemetryPoint()
	validPoint := domain.TelemetryPoint{
		VehicleID:       fixture.VehicleID,
		Latitude:        fixture.Latitude,
		Longitude:       fixture.Longitude,
		DeviceTimestamp: time.Now().Add(-10 * time.Second),
	}

	checkerSuccessMock := &resmocks.VehicleChecker{}
	checkerSuccessMock.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)

	repoSuccessMock := &repomocks.TelemetryRepository{}
	repoSuccessMock.On("Create", anyctx, anypoint).Return(fixture, nil)

	tests := []struct {
		name          string
		point         domain.TelemetryPoint
		dependencies  Dependencies
		expected      domain.TelemetryPoint
		expectedError error
	}{
		{
			name:  "works correctly",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: checkerSuccessMock,
				Telemetry:      repoSuccessMock,
			},
			expected: fixture,
		},
		{
			name: "handles correctly when latitude is invalid",
			point: domain.TelemetryPoint{
				VehicleID:       fixture.VehicleID,
				Latitude:        91.0,
				Longitude:       fixture.Longitude,
				DeviceTimestamp: time.Now().Add(-1 * time.Second),
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrInvalidLatitude,
		},
		{
			name: "handles correctly when longitude is invalid",
			point: domain.TelemetryPoint{
				VehicleID:       fixture.VehicleID,
				Latitude:        fixture.Latitude,
				Longitude:       181.0,
				DeviceTimestamp: time.Now().Add(-1 * time.Second),
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrInvalidLongitude,
		},
		{
			name: "handles correctly when timestamp is zero",
			point: domain.TelemetryPoint{
				VehicleID: fixture.VehicleID,
				Latitude:  fixture.Latitude,
				Longitude: fixture.Longitude,
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrInvalidTimestamp,
		},
		{
			name: "handles correctly when timestamp is too far in the future",
			point: domain.TelemetryPoint{
				VehicleID:       fixture.VehicleID,
				Latitude:        fixture.Latitude,
				Longitude:       fixture.Longitude,
				DeviceTimestamp: time.Now().Add(10 * time.Minute),
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrFutureTelemetry,
		},
		{
			name:  "handles correctly when vehicle is not found",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					mockChecker := &resmocks.VehicleChecker{}
					mockChecker.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(false, nil)
					return mockChecker
				}(),
			},
			expectedError: domain.ErrVehicleNotFound,
		},
		{
			name:  "fails when VehicleChecker ExistsAndActive fails",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					mockChecker := &resmocks.VehicleChecker{}
					mockChecker.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(false, errTest)
					return mockChecker
				}(),
			},
			expectedError: errTest,
		},
		{
			name:  "handles correctly when telemetry is duplicate",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					mockChecker := &resmocks.VehicleChecker{}
					mockChecker.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return mockChecker
				}(),
				Telemetry: func() *repomocks.TelemetryRepository {
					mockRepo := &repomocks.TelemetryRepository{}
					mockRepo.On("Create", anyctx, anypoint).Return(domain.TelemetryPoint{}, domain.ErrDuplicateTelemetry)
					return mockRepo
				}(),
			},
			expectedError: domain.ErrDuplicateTelemetry,
		},
		{
			name:  "fails when repository Create fails",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					mockChecker := &resmocks.VehicleChecker{}
					mockChecker.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return mockChecker
				}(),
				Telemetry: func() *repomocks.TelemetryRepository {
					mockRepo := &repomocks.TelemetryRepository{}
					mockRepo.On("Create", anyctx, anypoint).Return(domain.TelemetryPoint{}, errTest)
					return mockRepo
				}(),
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newService(t, tt.dependencies)
			result, err := svc.IngestTelemetry(context.Background(), tt.point)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expected, result)
			}

			if tt.dependencies.VehicleChecker != nil {
				tt.dependencies.VehicleChecker.AssertExpectations(t)
			}
			if tt.dependencies.Telemetry != nil {
				tt.dependencies.Telemetry.AssertExpectations(t)
			}
		})
	}
}

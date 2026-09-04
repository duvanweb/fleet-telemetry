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
	errTest  = errors.New("unexpected error")
	anyctx   = mock.Anything
	anypoint = mock.Anything
	anyevent = mock.Anything
	anyttl   = mock.Anything
	anystr   = mock.Anything
)

// Dependencies lists mocks in execution order of IngestTelemetry.
type Dependencies struct {
	VehicleChecker *resmocks.VehicleChecker
	Cache          *resmocks.TelemetryCache
	Telemetry      *repomocks.TelemetryRepository
}

func newService(t *testing.T, deps Dependencies) *telemetry.Service {
	t.Helper()
	log, _ := logger.NewLogger()
	return telemetry.NewService(log,
		telemetry.Repositories{Telemetry: deps.Telemetry},
		telemetry.Resources{VehicleChecker: deps.VehicleChecker, Cache: deps.Cache},
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

	cacheSuccessMock := &resmocks.TelemetryCache{}
	cacheSuccessMock.On("CheckDedup", anyctx, anystr).Return(false, nil)
	cacheSuccessMock.On("SetDedup", anyctx, anystr, anyttl).Return(nil)
	cacheSuccessMock.On("GetLastPosition", anyctx, fixture.VehicleID).Return(domain.TelemetryPoint{}, false, nil)
	cacheSuccessMock.On("SetLastPosition", anyctx, fixture.VehicleID, anypoint, anyttl).Return(nil)

	repoSuccessMock := &repomocks.TelemetryRepository{}
	repoSuccessMock.On("CreateWithOutbox", anyctx, anypoint, anyevent).Return(fixture, nil)

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
				Cache:          cacheSuccessMock,
				Telemetry:      repoSuccessMock,
			},
			expected: fixture,
		},
		{
			name: "works correctly when point is out of order",
			point: domain.TelemetryPoint{
				VehicleID:       fixture.VehicleID,
				Latitude:        fixture.Latitude,
				Longitude:       fixture.Longitude,
				DeviceTimestamp: time.Now().Add(-2 * time.Minute),
			},
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return m
				}(),
				Cache: func() *resmocks.TelemetryCache {
					m := &resmocks.TelemetryCache{}
					m.On("CheckDedup", anyctx, anystr).Return(false, nil)
					m.On("SetDedup", anyctx, anystr, anyttl).Return(nil)
					m.On("GetLastPosition", anyctx, fixture.VehicleID).Return(
						domain.TelemetryPoint{DeviceTimestamp: time.Now().Add(-1 * time.Minute)},
						true, nil,
					)
					return m
				}(),
				Telemetry: func() *repomocks.TelemetryRepository {
					m := &repomocks.TelemetryRepository{}
					m.On("CreateWithOutbox", anyctx, anypoint, anyevent).Return(fixture, nil)
					return m
				}(),
			},
			expected: fixture,
		},
		{
			name: "handles correctly when latitude is invalid",
			point: domain.TelemetryPoint{
				VehicleID: fixture.VehicleID, Latitude: 91.0, Longitude: fixture.Longitude,
				DeviceTimestamp: time.Now().Add(-1 * time.Second),
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrInvalidLatitude,
		},
		{
			name: "handles correctly when longitude is invalid",
			point: domain.TelemetryPoint{
				VehicleID: fixture.VehicleID, Latitude: fixture.Latitude, Longitude: 181.0,
				DeviceTimestamp: time.Now().Add(-1 * time.Second),
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrInvalidLongitude,
		},
		{
			name: "handles correctly when timestamp is zero",
			point: domain.TelemetryPoint{
				VehicleID: fixture.VehicleID, Latitude: fixture.Latitude, Longitude: fixture.Longitude,
			},
			dependencies:  Dependencies{},
			expectedError: domain.ErrInvalidTimestamp,
		},
		{
			name: "handles correctly when timestamp is too far in the future",
			point: domain.TelemetryPoint{
				VehicleID: fixture.VehicleID, Latitude: fixture.Latitude, Longitude: fixture.Longitude,
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
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(false, nil)
					return m
				}(),
			},
			expectedError: domain.ErrVehicleNotFound,
		},
		{
			name:  "handles correctly when Redis dedup key exists",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return m
				}(),
				Cache: func() *resmocks.TelemetryCache {
					m := &resmocks.TelemetryCache{}
					m.On("CheckDedup", anyctx, anystr).Return(true, nil)
					return m
				}(),
			},
			expectedError: domain.ErrDuplicateTelemetry,
		},
		{
			name:  "handles correctly when DB returns duplicate telemetry",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return m
				}(),
				Cache: func() *resmocks.TelemetryCache {
					m := &resmocks.TelemetryCache{}
					m.On("CheckDedup", anyctx, anystr).Return(false, nil)
					m.On("SetDedup", anyctx, anystr, anyttl).Return(nil)
					return m
				}(),
				Telemetry: func() *repomocks.TelemetryRepository {
					m := &repomocks.TelemetryRepository{}
					m.On("CreateWithOutbox", anyctx, anypoint, anyevent).Return(domain.TelemetryPoint{}, domain.ErrDuplicateTelemetry)
					return m
				}(),
			},
			expectedError: domain.ErrDuplicateTelemetry,
		},
		{
			name:  "fails when VehicleChecker ExistsAndActive fails",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(false, errTest)
					return m
				}(),
			},
			expectedError: errTest,
		},
		{
			name:  "fails when Cache CheckDedup fails",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return m
				}(),
				Cache: func() *resmocks.TelemetryCache {
					m := &resmocks.TelemetryCache{}
					m.On("CheckDedup", anyctx, anystr).Return(false, errTest)
					return m
				}(),
			},
			expectedError: errTest,
		},
		{
			name:  "fails when Cache SetDedup fails",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return m
				}(),
				Cache: func() *resmocks.TelemetryCache {
					m := &resmocks.TelemetryCache{}
					m.On("CheckDedup", anyctx, anystr).Return(false, nil)
					m.On("SetDedup", anyctx, anystr, anyttl).Return(errTest)
					return m
				}(),
			},
			expectedError: errTest,
		},
		{
			name:  "fails when repository CreateWithOutbox fails",
			point: validPoint,
			dependencies: Dependencies{
				VehicleChecker: func() *resmocks.VehicleChecker {
					m := &resmocks.VehicleChecker{}
					m.On("ExistsAndActive", anyctx, fixture.VehicleID).Return(true, nil)
					return m
				}(),
				Cache: func() *resmocks.TelemetryCache {
					m := &resmocks.TelemetryCache{}
					m.On("CheckDedup", anyctx, anystr).Return(false, nil)
					m.On("SetDedup", anyctx, anystr, anyttl).Return(nil)
					return m
				}(),
				Telemetry: func() *repomocks.TelemetryRepository {
					m := &repomocks.TelemetryRepository{}
					m.On("CreateWithOutbox", anyctx, anypoint, anyevent).Return(domain.TelemetryPoint{}, errTest)
					return m
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
			if tt.dependencies.Cache != nil {
				tt.dependencies.Cache.AssertExpectations(t)
			}
			if tt.dependencies.Telemetry != nil {
				tt.dependencies.Telemetry.AssertExpectations(t)
			}
		})
	}
}

package alerts_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fleet/alert-service/internal/core/alerts"
	"fleet/alert-service/internal/core/domain"
	repomocks "fleet/alert-service/internal/core/ports/repositories/mocks"
	testdata "fleet/alert-service/test/data"
	"fleet/shared/pkg/logger"
)

var (
	errTest  = errors.New("unexpected error")
	anyctx   = mock.Anything
	anyalert = mock.Anything
)

func newService(t *testing.T, repo *repomocks.AlertRepository) *alerts.Service {
	t.Helper()
	log, _ := logger.NewLogger()
	return alerts.NewService(log, alerts.Dependencies{Repository: repo})
}

func TestService_CreateAlert(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()
	startedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	repositorySuccessMock := func() *repomocks.AlertRepository {
		m := &repomocks.AlertRepository{}
		m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, nil)
		m.On("Create", anyctx, anyalert).Return(fixture, nil)
		return m
	}

	tests := []struct {
		name          string
		dependencies  func() *repomocks.AlertRepository
		vehicleID     string
		alertType     string
		startedAt     time.Time
		expected      domain.Alert
		expectedError error
	}{
		{
			name:         "works correctly",
			dependencies: repositorySuccessMock,
			vehicleID:    fixture.VehicleID,
			alertType:    domain.AlertTypeVehicleStopped,
			startedAt:    startedAt,
			expected:     fixture,
		},
		{
			name: "handles correctly when open alert already exists",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(fixture, true, nil)
				return m
			},
			vehicleID: fixture.VehicleID,
			alertType: domain.AlertTypeVehicleStopped,
			startedAt: startedAt,
			expected:  fixture,
		},
		{
			name: "fails when AlertRepository GetOpenByVehicle fails",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, errTest)
				return m
			},
			vehicleID:     fixture.VehicleID,
			alertType:     domain.AlertTypeVehicleStopped,
			startedAt:     startedAt,
			expected:      domain.Alert{},
			expectedError: errTest,
		},
		{
			name: "fails when AlertRepository Create fails",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, nil)
				m.On("Create", anyctx, anyalert).Return(domain.Alert{}, errTest)
				return m
			},
			vehicleID:     fixture.VehicleID,
			alertType:     domain.AlertTypeVehicleStopped,
			startedAt:     startedAt,
			expected:      domain.Alert{},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tt.dependencies()
			svc := newService(t, repoMock)

			result, err := svc.CreateAlert(context.Background(), tt.vehicleID, tt.alertType, tt.startedAt)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedError, err)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_GetAll(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	repositorySuccessMock := func() *repomocks.AlertRepository {
		m := &repomocks.AlertRepository{}
		m.On("GetAll", anyctx).Return([]domain.Alert{fixture}, nil)
		return m
	}

	tests := []struct {
		name          string
		dependencies  func() *repomocks.AlertRepository
		expected      []domain.Alert
		expectedError error
	}{
		{
			name:         "works correctly",
			dependencies: repositorySuccessMock,
			expected:     []domain.Alert{fixture},
		},
		{
			name: "works correctly when no alerts exist",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetAll", anyctx).Return(nil, nil)
				return m
			},
			expected: nil,
		},
		{
			name: "fails when AlertRepository GetAll fails",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetAll", anyctx).Return(nil, errTest)
				return m
			},
			expected:      nil,
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tt.dependencies()
			svc := newService(t, repoMock)

			result, err := svc.GetAll(context.Background())

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedError, err)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_GetOpenByVehicle(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	repositorySuccessMock := func() *repomocks.AlertRepository {
		m := &repomocks.AlertRepository{}
		m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(fixture, true, nil)
		return m
	}

	tests := []struct {
		name          string
		dependencies  func() *repomocks.AlertRepository
		vehicleID     string
		expected      domain.Alert
		expectedFound bool
		expectedError error
	}{
		{
			name:          "works correctly when alert exists",
			dependencies:  repositorySuccessMock,
			vehicleID:     fixture.VehicleID,
			expected:      fixture,
			expectedFound: true,
		},
		{
			name: "works correctly when no open alert",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, nil)
				return m
			},
			vehicleID:     fixture.VehicleID,
			expected:      domain.Alert{},
			expectedFound: false,
		},
		{
			name: "fails when AlertRepository GetOpenByVehicle fails",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, errTest)
				return m
			},
			vehicleID:     fixture.VehicleID,
			expected:      domain.Alert{},
			expectedFound: false,
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tt.dependencies()
			svc := newService(t, repoMock)

			result, found, err := svc.GetOpenByVehicle(context.Background(), tt.vehicleID)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedFound, found)
			assert.Equal(t, tt.expectedError, err)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestService_ResolveAlert(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestAlert()

	repositorySuccessMock := func() *repomocks.AlertRepository {
		m := &repomocks.AlertRepository{}
		m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(fixture, true, nil)
		m.On("MarkResolved", anyctx, fixture.ID).Return(nil)
		return m
	}

	tests := []struct {
		name          string
		dependencies  func() *repomocks.AlertRepository
		vehicleID     string
		expectedError error
	}{
		{
			name:         "works correctly",
			dependencies: repositorySuccessMock,
			vehicleID:    fixture.VehicleID,
		},
		{
			name: "handles correctly when no open alert exists",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, nil)
				return m
			},
			vehicleID: fixture.VehicleID,
		},
		{
			name: "fails when AlertRepository GetOpenByVehicle fails",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(domain.Alert{}, false, errTest)
				return m
			},
			vehicleID:     fixture.VehicleID,
			expectedError: errTest,
		},
		{
			name: "fails when AlertRepository MarkResolved fails",
			dependencies: func() *repomocks.AlertRepository {
				m := &repomocks.AlertRepository{}
				m.On("GetOpenByVehicle", anyctx, fixture.VehicleID).Return(fixture, true, nil)
				m.On("MarkResolved", anyctx, fixture.ID).Return(errTest)
				return m
			},
			vehicleID:     fixture.VehicleID,
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tt.dependencies()
			svc := newService(t, repoMock)

			err := svc.ResolveAlert(context.Background(), tt.vehicleID)

			assert.Equal(t, tt.expectedError, err)
			repoMock.AssertExpectations(t)
		})
	}
}

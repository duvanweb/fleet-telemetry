package stopped_detection_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fleet/alert-service/internal/core/domain"
	svcmocks "fleet/alert-service/internal/core/ports/services/mocks"
	"fleet/alert-service/internal/core/stopped_detection"
	testdata "fleet/alert-service/test/data"
	"fleet/shared/pkg/logger"
)

var (
	errTest = errors.New("unexpected error")
	anyctx  = mock.Anything
)

const testThreshold = time.Minute

var (
	t0      = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	posA    = domain.TelemetryEvent{VehicleID: "vehicle-1", Latitude: 1.0, Longitude: 1.0}
	posB    = domain.TelemetryEvent{VehicleID: "vehicle-1", Latitude: 4.711, Longitude: -74.0721}
	posDiff = domain.TelemetryEvent{VehicleID: "vehicle-1", Latitude: 5.0, Longitude: 5.0}
)

func newService(t *testing.T, alerts *svcmocks.AlertService, clockTimes []time.Time) *stopped_detection.Service {
	t.Helper()
	log, _ := logger.NewLogger()
	idx := 0
	clockFn := func() time.Time {
		if idx >= len(clockTimes) {
			return clockTimes[len(clockTimes)-1]
		}
		ct := clockTimes[idx]
		idx++
		return ct
	}
	return stopped_detection.NewService(log, alerts, testThreshold, clockFn)
}

func TestService_Evaluate(t *testing.T) {
	t.Parallel()

	fixture := testdata.GetTestTelemetryEvent()

	tests := []struct {
		name          string
		clockTimes    []time.Time
		events        []domain.TelemetryEvent
		mockSetup     func(m *svcmocks.AlertService)
		expectedError error
	}{
		{
			name:       "works correctly when same position under threshold",
			clockTimes: []time.Time{t0, t0, t0, t0.Add(30 * time.Second)},
			events:     []domain.TelemetryEvent{posA, posB, posB, posB},
			mockSetup:  func(m *svcmocks.AlertService) {},
		},
		{
			name:       "works correctly when same position at threshold creates alert",
			clockTimes: []time.Time{t0, t0, t0, t0.Add(2 * time.Minute)},
			events:     []domain.TelemetryEvent{posA, posB, posB, posB},
			mockSetup: func(m *svcmocks.AlertService) {
				m.On("CreateAlert", anyctx, fixture.VehicleID, domain.AlertTypeVehicleStopped, t0).
					Return(testdata.GetTestAlert(), nil)
			},
		},
		{
			name:       "works correctly when same position with open alert no duplicate",
			clockTimes: []time.Time{t0, t0, t0, t0.Add(2 * time.Minute), t0.Add(3 * time.Minute)},
			events:     []domain.TelemetryEvent{posA, posB, posB, posB, posB},
			mockSetup: func(m *svcmocks.AlertService) {
				m.On("CreateAlert", anyctx, fixture.VehicleID, domain.AlertTypeVehicleStopped, t0).
					Return(testdata.GetTestAlert(), nil).Once()
			},
		},
		{
			name:       "works correctly when different position resolves alert",
			clockTimes: []time.Time{t0, t0, t0, t0.Add(2 * time.Minute), t0.Add(3 * time.Minute)},
			events:     []domain.TelemetryEvent{posA, posB, posB, posB, posDiff},
			mockSetup: func(m *svcmocks.AlertService) {
				m.On("CreateAlert", anyctx, fixture.VehicleID, domain.AlertTypeVehicleStopped, t0).
					Return(testdata.GetTestAlert(), nil).Once()
				m.On("ResolveAlert", anyctx, fixture.VehicleID).Return(nil)
			},
		},
		{
			name:       "fails when AlertService CreateAlert fails",
			clockTimes: []time.Time{t0, t0, t0, t0.Add(2 * time.Minute)},
			events:     []domain.TelemetryEvent{posA, posB, posB, posB},
			mockSetup: func(m *svcmocks.AlertService) {
				m.On("CreateAlert", anyctx, fixture.VehicleID, domain.AlertTypeVehicleStopped, t0).
					Return(domain.Alert{}, errTest)
			},
			expectedError: errTest,
		},
		{
			name:       "fails when AlertService ResolveAlert fails",
			clockTimes: []time.Time{t0, t0, t0, t0.Add(2 * time.Minute), t0.Add(3 * time.Minute)},
			events:     []domain.TelemetryEvent{posA, posB, posB, posB, posDiff},
			mockSetup: func(m *svcmocks.AlertService) {
				m.On("CreateAlert", anyctx, fixture.VehicleID, domain.AlertTypeVehicleStopped, t0).
					Return(testdata.GetTestAlert(), nil).Once()
				m.On("ResolveAlert", anyctx, fixture.VehicleID).Return(errTest)
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			alertMock := &svcmocks.AlertService{}
			tt.mockSetup(alertMock)

			svc := newService(t, alertMock, tt.clockTimes)

			var err error
			for _, event := range tt.events {
				err = svc.Evaluate(context.Background(), event)
				if err != nil {
					break
				}
			}

			assert.Equal(t, tt.expectedError, err)
			alertMock.AssertExpectations(t)
		})
	}
}

package vehicle_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fleet/vehicle-service/internal/core/domain"
	"fleet/vehicle-service/internal/core/ports/repositories/mocks"
	"fleet/vehicle-service/internal/core/vehicle"
	"fleet/vehicle-service/test/data"
	"fleet/shared/pkg/logger"
)

var (
	errTest    = errors.New("unexpected error")
	anyctx     = mock.Anything
	anyvehicle = mock.Anything
)

type Dependencies struct {
	Vehicle *mocks.VehicleRepository
}

func newService(t *testing.T, deps Dependencies) *vehicle.Service {
	t.Helper()
	log, _ := logger.NewLogger()
	return vehicle.NewService(log, vehicle.Repositories{Vehicle: deps.Vehicle})
}

func TestVehicleService_Create(t *testing.T) {
	t.Parallel()

	v := data.GetTestVehicle()

	repoSuccessMock := &mocks.VehicleRepository{}
	repoSuccessMock.On("ExistsByPlate", anyctx, v.Plate).Return(false, nil)
	repoSuccessMock.On("ExistsByExternalID", anyctx, v.ExternalID).Return(false, nil)
	repoSuccessMock.On("Create", anyctx, anyvehicle).Return(v, nil)

	tests := []struct {
		name          string
		dependencies  Dependencies
		expectedError error
	}{
		{
			name:         "works correctly",
			dependencies: Dependencies{Vehicle: repoSuccessMock},
		},
		{
			name: "handles correctly when plate already exists",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("ExistsByPlate", anyctx, v.Plate).Return(true, nil)
					return mockRepo
				}(),
			},
			expectedError: domain.ErrDuplicatePlate,
		},
		{
			name: "fails when ExistsByPlate fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("ExistsByPlate", anyctx, v.Plate).Return(false, errTest)
					return mockRepo
				}(),
			},
			expectedError: errTest,
		},
		{
			name: "fails when ExistsByExternalID fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("ExistsByPlate", anyctx, v.Plate).Return(false, nil)
					mockRepo.On("ExistsByExternalID", anyctx, v.ExternalID).Return(false, errTest)
					return mockRepo
				}(),
			},
			expectedError: errTest,
		},
		{
			name: "fails when Create fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("ExistsByPlate", anyctx, v.Plate).Return(false, nil)
					mockRepo.On("ExistsByExternalID", anyctx, v.ExternalID).Return(false, nil)
					mockRepo.On("Create", anyctx, anyvehicle).Return(domain.Vehicle{}, errTest)
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

			input := domain.Vehicle{ExternalID: v.ExternalID, Plate: v.Plate, Name: v.Name}
			_, err := svc.Create(t.Context(), input)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			tt.dependencies.Vehicle.AssertExpectations(t)
		})
	}
}

func TestVehicleService_Delete(t *testing.T) {
	t.Parallel()

	v := data.GetTestVehicle()

	repoSuccessMock := &mocks.VehicleRepository{}
	repoSuccessMock.On("SoftDelete", anyctx, v.ID).Return(nil)

	tests := []struct {
		name          string
		dependencies  Dependencies
		expectedError error
	}{
		{
			name:         "works correctly",
			dependencies: Dependencies{Vehicle: repoSuccessMock},
		},
		{
			name: "fails when vehicle not found",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("SoftDelete", anyctx, v.ID).Return(domain.ErrVehicleNotFound)
					return mockRepo
				}(),
			},
			expectedError: domain.ErrVehicleNotFound,
		},
		{
			name: "fails when SoftDelete fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("SoftDelete", anyctx, v.ID).Return(errTest)
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

			err := svc.Delete(t.Context(), v.ID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			tt.dependencies.Vehicle.AssertExpectations(t)
		})
	}
}

func TestVehicleService_GetAll(t *testing.T) {
	t.Parallel()

	v := data.GetTestVehicle()

	repoSuccessMock := &mocks.VehicleRepository{}
	repoSuccessMock.On("GetAll", anyctx).Return([]domain.Vehicle{v}, nil)

	tests := []struct {
		name             string
		dependencies     Dependencies
		expectedVehicles []domain.Vehicle
		expectedError    error
	}{
		{
			name:             "works correctly",
			dependencies:     Dependencies{Vehicle: repoSuccessMock},
			expectedVehicles: []domain.Vehicle{v},
		},
		{
			name: "fails when GetAll fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("GetAll", anyctx).Return(nil, errTest)
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

			result, err := svc.GetAll(t.Context())

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicles, result)
			}
			tt.dependencies.Vehicle.AssertExpectations(t)
		})
	}
}

func TestVehicleService_GetByID(t *testing.T) {
	t.Parallel()

	v := data.GetTestVehicle()

	repoSuccessMock := &mocks.VehicleRepository{}
	repoSuccessMock.On("GetByID", anyctx, v.ID).Return(v, nil)

	tests := []struct {
		name            string
		dependencies    Dependencies
		expectedVehicle domain.Vehicle
		expectedError   error
	}{
		{
			name:            "works correctly",
			dependencies:    Dependencies{Vehicle: repoSuccessMock},
			expectedVehicle: v,
		},
		{
			name: "handles correctly when vehicle is deleted",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					deleted := data.GetTestDeletedVehicle()
					mockRepo.On("GetByID", anyctx, deleted.ID).Return(deleted, nil)
					return mockRepo
				}(),
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrVehicleDeleted,
		},
		{
			name: "fails when vehicle not found",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("GetByID", anyctx, v.ID).Return(domain.Vehicle{}, domain.ErrVehicleNotFound)
					return mockRepo
				}(),
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrVehicleNotFound,
		},
		{
			name: "fails when GetByID fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("GetByID", anyctx, v.ID).Return(domain.Vehicle{}, errTest)
					return mockRepo
				}(),
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := v.ID
			if tt.expectedError != nil && errors.Is(tt.expectedError, domain.ErrVehicleDeleted) {
				id = data.GetTestDeletedVehicle().ID
			}

			svc := newService(t, tt.dependencies)
			result, err := svc.GetByID(t.Context(), id)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedVehicle, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicle, result)
			}
			tt.dependencies.Vehicle.AssertExpectations(t)
		})
	}
}

func TestVehicleService_Update(t *testing.T) {
	t.Parallel()

	v := data.GetTestVehicle()

	repoSuccessMock := &mocks.VehicleRepository{}
	repoSuccessMock.On("Update", anyctx, anyvehicle).Return(v, nil)

	tests := []struct {
		name            string
		dependencies    Dependencies
		expectedVehicle domain.Vehicle
		expectedError   error
	}{
		{
			name:            "works correctly",
			dependencies:    Dependencies{Vehicle: repoSuccessMock},
			expectedVehicle: v,
		},
		{
			name: "fails when vehicle not found",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("Update", anyctx, anyvehicle).Return(domain.Vehicle{}, domain.ErrVehicleNotFound)
					return mockRepo
				}(),
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   domain.ErrVehicleNotFound,
		},
		{
			name: "fails when Update fails",
			dependencies: Dependencies{
				Vehicle: func() *mocks.VehicleRepository {
					mockRepo := &mocks.VehicleRepository{}
					mockRepo.On("Update", anyctx, anyvehicle).Return(domain.Vehicle{}, errTest)
					return mockRepo
				}(),
			},
			expectedVehicle: domain.Vehicle{},
			expectedError:   errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := newService(t, tt.dependencies)

			result, err := svc.Update(t.Context(), v)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedVehicle, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVehicle, result)
			}
			tt.dependencies.Vehicle.AssertExpectations(t)
		})
	}
}

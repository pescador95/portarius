package resident_test

import (
	"testing"

	"portarius/internal/resident/domain"
	mock_repository "portarius/internal/resident/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*gomock.Controller, *mock_repository.MockIResidentRepository, *domain.ResidentService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockIResidentRepository(ctrl)
	service := domain.NewResidentService(mockRepo)
	return ctrl, mockRepo, service
}

func validResident() *domain.Resident {
	return &domain.Resident{
		Name:         "João Silva",
		Document:     "123.456.789-00",
		Phone:        "(11) 98765-4321",
		ResidentType: domain.Tenant,
		Email:        "joao@example.com",
		Apartment:    "42",
		Block:        "B",
	}
}

func TestResidentService_CreateResident(t *testing.T) {
	t.Run("should create resident successfully", func(t *testing.T) {
		ctrl, mockRepo, service := setupTest(t)
		defer ctrl.Finish()

		res := validResident()

		mockRepo.EXPECT().
			Create(gomock.Any()).
			DoAndReturn(func(r *domain.Resident) error {
				assert.Equal(t, "João Silva", r.Name)
				return nil
			})

		err := service.CreateResident(res)
		assert.NoError(t, err)
	})
}

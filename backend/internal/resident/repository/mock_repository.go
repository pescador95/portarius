package repository

import (
	"portarius/internal/resident/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockResidentRepository struct {
	GetAllFunc  func() ([]domain.Resident, error)
	GetByIDFunc func(id uint) (*domain.Resident, error)
	CreateFunc  func(resident *domain.Resident) error
	UpdateFunc  func(resident *domain.Resident) error
	DeleteFunc  func(id uint) error
}

func (m *MockResidentRepository) GetAll() ([]domain.Resident, error) {
	return m.GetAllFunc()
}

func (m *MockResidentRepository) GetByID(id uint) (*domain.Resident, error) {
	return m.GetByIDFunc(id)
}

func (m *MockResidentRepository) Create(resident *domain.Resident) error {
	return m.CreateFunc(resident)
}

func (m *MockResidentRepository) Update(resident *domain.Resident) error {
	return m.UpdateFunc(resident)
}

func (m *MockResidentRepository) Delete(id uint) error {
	return m.DeleteFunc(id)
}

func TestMockResidentRepository(t *testing.T) {

	mockRepo := &MockResidentRepository{
		GetAllFunc: func() ([]domain.Resident, error) {
			return []domain.Resident{{Name: "John Doe"}}, nil
		},
		GetByIDFunc: func(id uint) (*domain.Resident, error) {
			return &domain.Resident{Name: "John Doe"}, nil
		},
		CreateFunc: func(resident *domain.Resident) error {
			return nil
		},
		UpdateFunc: func(resident *domain.Resident) error {
			return nil
		},
		DeleteFunc: func(id uint) error {
			return nil
		},
	}

	residents, err := mockRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, residents)
	assert.Len(t, residents, 1)

	resident, err := mockRepo.GetByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, resident)
	assert.Equal(t, uint(1), resident.ID)

	err = mockRepo.Create(&domain.Resident{Name: "Jane Doe"})
	assert.NoError(t, err)

	err = mockRepo.Update(&domain.Resident{Name: "John Doe Updated"})
	assert.NoError(t, err)

	err = mockRepo.Delete(1)
	assert.NoError(t, err)
}

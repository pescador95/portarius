package resident

import "portarius/internal/resident/domain"

type IResidentService interface {
	GetAllResidents() ([]domain.Resident, error)
	GetResidentByID(id uint) (*domain.Resident, error)
	CreateResident(resident *domain.Resident) error
	UpdateResident(resident *domain.Resident) error
	DeleteResident(id uint) error
}

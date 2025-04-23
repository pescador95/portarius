package repository

import (
	"portarius/internal/resident/domain"

	"gorm.io/gorm"
)

type residentRepository struct {
	db *gorm.DB
}

func NewResidentRepository(db *gorm.DB) domain.IResidentRepository {
	return &residentRepository{db: db}
}

func (r *residentRepository) GetAll() ([]domain.Resident, error) {
	var residents []domain.Resident
	err := r.db.Find(&residents).Error
	return residents, err
}

func (r *residentRepository) GetByID(id uint) (*domain.Resident, error) {
	var resident domain.Resident
	err := r.db.First(&resident, id).Error
	return &resident, err
}

func (r *residentRepository) Create(resident *domain.Resident) error {
	return r.db.Create(resident).Error
}

func (r *residentRepository) Update(resident *domain.Resident) error {
	return r.db.Save(resident).Error
}

func (r *residentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Resident{}, id).Error
}

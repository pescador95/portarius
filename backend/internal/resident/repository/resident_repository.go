package repository

import (
	packageDomain "portarius/internal/package/domain"
	reservationDomain "portarius/internal/reservation/domain"
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

func (r *residentRepository) GetPhoneByReservationID(reservationID uint) (string, error) {
	var phone string
	err := r.db.Model(&reservationDomain.Reservation{}).
		Joins("JOIN residents ON residents.id = reservations.resident_id").
		Where("reservations.id = ?", reservationID).
		Select("residents.phone").
		Scan(&phone).Error
	if err != nil {
		return "", err
	}
	return phone, nil
}

func (r *residentRepository) GetPhoneByPackageID(packageID uint) (string, error) {
	var phone string
	err := r.db.Model(&packageDomain.Package{}).
		Joins("JOIN residents ON residents.id = packages.resident_id").
		Where("packages.id = ?", packageID).
		Select("residents.phone").
		Scan(&phone).Error
	if err != nil {
		return "", err
	}
	return phone, nil
}

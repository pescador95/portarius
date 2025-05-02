package repository

import (
	"fmt"
	"portarius/internal/infra"
	"portarius/internal/reservation/domain"
	"time"

	"gorm.io/gorm"
)

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) domain.IReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) GetAll(page, pageSize int) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	err := r.db.Preload("Resident").Scopes(infra.Paginate(page, pageSize)).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetByID(id uint) (*domain.Reservation, error) {
	var reservation domain.Reservation
	err := r.db.Preload("Resident").First(&reservation, id).Error
	return &reservation, err
}

func (r *reservationRepository) Create(reservation *domain.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) Update(reservation *domain.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Reservation{}, id).Error
}

func (r *reservationRepository) GetByResident(residentID uint) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	err := r.db.Where("resident_id = ?", residentID).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetBySpace(space string) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	err := r.db.Where("space = ?", space).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetByStatus(status string) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	err := r.db.Where("status = ?", status).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByDateRange(startDate, endDate time.Time) ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	err := r.db.Where("start_time BETWEEN ? AND ?", startDate, endDate).Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindUpcomingReservations() ([]domain.Reservation, error) {
	var reservations []domain.Reservation
	err := r.db.
		Where("start_time > ? AND status != ?", time.Now(), domain.StatusCancelled).
		Order("start_time ASC").
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	var reservation domain.Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		return err
	}

	reservation.Status = domain.ReservationStatus(status)
	return r.db.Save(&reservation).Error
}

func (r *reservationRepository) ImportSalonReservations(reservations []domain.Reservation) error {
	return r.db.Create(&reservations).Error
}

func (r *reservationRepository) CheckReservationConflict(space string, startTime, endTime time.Time, excludeID uint) error {
	var count int64
	query := r.db.Model(&domain.Reservation{}).
		Where("space = ? AND status NOT IN ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?))",
			space,
			[]domain.ReservationStatus{domain.StatusCancelled, domain.StatusKeysReturned},
			startTime,
			endTime,
			startTime,
			endTime,
			startTime,
			endTime)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("já existe uma reserva para este salão no horário selecionado")
	}

	return nil
}

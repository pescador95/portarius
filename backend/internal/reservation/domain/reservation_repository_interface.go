package domain

import "time"

type IReservationRepository interface {
	GetAll(page, pageSize int) ([]Reservation, error)
	GetByID(id uint) (*Reservation, error)
	Create(reservation *Reservation) error
	Update(reservation *Reservation) error
	Delete(id uint) error
	GetByResident(residentID uint) ([]Reservation, error)
	GetBySpace(space string) ([]Reservation, error)
	GetByStatus(status string) ([]Reservation, error)
	FindByDateRange(startDate, endDate time.Time) ([]Reservation, error)
	FindUpcomingReservations() ([]Reservation, error)
	UpdateStatus(id uint, status string) error
	ImportSalonReservations(reservations []Reservation) error
	CheckReservationConflict(space string, startTime, endTime time.Time, excludeID uint) error
}

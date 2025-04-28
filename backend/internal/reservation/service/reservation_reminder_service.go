package reservation

import (
	reminderDomain "portarius/internal/reminder/domain"
	reservationDomain "portarius/internal/reservation/domain"
)

type ReservationReminderService struct {
	repo reminderDomain.IReminderRepository
}

func NewReservationReminderService(repo reminderDomain.IReminderRepository) *ReservationReminderService {
	return &ReservationReminderService{repo: repo}
}

func (s *ReservationReminderService) CreateReminderForReservation(reservation *reservationDomain.Reservation) error {

	reminder := reminderDomain.Reminder{
		Recipient:     reservation.Resident.Phone,
		ReservationID: &reservation.ID,
		PackageID:     nil,
		Channel:       reminderDomain.ReminderChannelWhatsApp,
		Status:        reminderDomain.ReminderStatusPending,
	}
	return s.repo.Create(&reminder)
}

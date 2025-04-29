package domain

import (
	"portarius/internal/eventbus"
)

func (r *Reminder) PublishPackageReminder() {
	if r.PackageID != nil && r.Status == ReminderStatusPending || r.Status == ReminderStatusFailed {

		eventbus.Publish("SendPackageReminder", &eventbus.ReminderEvent{
			ReminderID: &r.ID,
			PackageID:  r.PackageID,
			Phone:      r.Recipient,
		})
	}
}

func (r *Reminder) PublishReservationReminder() {
	if r.ReservationID != nil && r.Status == ReminderStatusPending || r.Status == ReminderStatusFailed {

		eventbus.Publish("SendReservationReminder", &eventbus.ReminderEvent{
			ReminderID:    &r.ID,
			ReservationID: r.ReservationID,
			Phone:         r.Recipient,
		})
	}
}

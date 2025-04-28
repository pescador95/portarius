package listeners

import (
	"portarius/internal/eventbus"
	"portarius/internal/reminder/domain"
	residentDomain "portarius/internal/resident/domain"
)

var reminderRepo domain.IReminderRepository
var residentRepo residentDomain.IResidentRepository

func RegisterReminderListeners(reminderRepository domain.IReminderRepository, residentRepository residentDomain.IResidentRepository) {
	reminderRepo = reminderRepository
	residentRepo = residentRepository

	eventbus.Subscribe("PackageCreated", onPackageCreated)
	eventbus.Subscribe("ReservationCreated", onReservationCreated)
}

func onPackageCreated(e eventbus.Event) {
	event := e.(eventbus.PackageCreatedEvent)

	phone, err := residentRepo.GetPhoneByPackageID(event.PackageID)
	if err != nil {

		return
	}

	reminder := domain.Reminder{
		PackageID: &event.PackageID,
		Recipient: string(phone),
		Channel:   domain.ReminderChannel(event.Channel),
		Status:    domain.ReminderStatusPending,
	}

	_ = reminderRepo.Create(&reminder)
}

func onReservationCreated(e eventbus.Event) {
	event := e.(eventbus.ReservationCreatedEvent)

	phone, err := residentRepo.GetPhoneByReservationID(event.ReservationID)
	if err != nil {

		return
	}

	reminder := domain.Reminder{
		ReservationID: &event.ReservationID,
		Recipient:     string(phone),
		Channel:       domain.ReminderChannel(event.Channel),
		Status:        domain.ReminderStatusPending,
	}

	_ = reminderRepo.Create(&reminder)
}

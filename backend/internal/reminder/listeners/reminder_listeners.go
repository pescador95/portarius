package listeners

import (
	"portarius/internal/eventbus"
	holydayHandler "portarius/internal/holyday/handler"
	packageDomain "portarius/internal/package/domain"
	reminderDomain "portarius/internal/reminder/domain"
	reservationDomain "portarius/internal/reservation/domain"
	residentDomain "portarius/internal/resident/domain"
	whatsAppDomain "portarius/internal/whatsapp/domain"
	"time"
)

var reminderRepo reminderDomain.IReminderRepository
var residentRepo residentDomain.IResidentRepository
var packageRepo packageDomain.IPackageRepository
var reservationRepo reservationDomain.IReservationRepository
var whatsappHandler whatsAppDomain.IWhatsAppHandler

func RegisterReminderListeners(reminderRepository reminderDomain.IReminderRepository, residentRepository residentDomain.IResidentRepository, packageRepository packageDomain.IPackageRepository, reservationRepository reservationDomain.IReservationRepository, handler whatsAppDomain.IWhatsAppHandler) {
	reminderRepo = reminderRepository
	residentRepo = residentRepository
	packageRepo = packageRepository
	reservationRepo = reservationRepository
	whatsappHandler = handler

	eventbus.Subscribe("PackageCreated", onPackageCreated)
	eventbus.Subscribe("ReservationCreated", onReservationCreated)
	eventbus.Subscribe("SendPackageReminder", onSendPackageReminder)
	eventbus.Subscribe("SendReservationReminder", onSendReservationReminder)
	eventbus.Subscribe("UpdateStatusReminder", onUpdateStatusReminder)
}

func onPackageCreated(e eventbus.Event) {
	event := e.(*eventbus.PackageCreatedEvent)

	phone, err := residentRepo.GetPhoneByPackageID(*event.PackageID)
	if err != nil {

		return
	}

	reminder := reminderDomain.Reminder{
		PackageID:   event.PackageID,
		Recipient:   string(phone),
		Channel:     reminderDomain.ReminderChannel(event.Channel),
		Status:      reminderDomain.ReminderStatusPending,
		ScheduledAt: time.Now(),
	}

	_ = reminderRepo.Create(&reminder)
}

func onReservationCreated(e eventbus.Event) {
	event := e.(*eventbus.ReservationCreatedEvent)

	phone, err := residentRepo.GetPhoneByReservationID(*event.ReservationID)
	if err != nil {

		return
	}

	scheduledAt := reservationDomain.GetReminderScheduleDate(event.StartTime, holydayHandler.IsHolyday)

	reminder := reminderDomain.Reminder{
		ReservationID: event.ReservationID,
		Recipient:     string(phone),
		Channel:       reminderDomain.ReminderChannel(event.Channel),
		Status:        reminderDomain.ReminderStatusPending,
		ScheduledAt:   scheduledAt,
	}

	_ = reminderRepo.Create(&reminder)
	reminder.PublishPackageReminder()
}

func onSendPackageReminder(e eventbus.Event) {
	event := e.(*eventbus.ReminderEvent)

	pck, err := packageRepo.GetByID(*event.PackageID)
	if err != nil {
		return
	}

	whatsappHandler.SendPackageNotification(*event.ReminderID, event.Phone, pck.Resident.Name)
}

func onSendReservationReminder(e eventbus.Event) {
	event := e.(*eventbus.ReminderEvent)

	reservation, err := reservationRepo.GetByID(*event.ReservationID)

	if err != nil {
		return
	}

	whatsappHandler.SendReservationKeyReminder(*event.ReminderID, event.Phone, reservation.Resident.Name, reservation.GetLastCharFromSalon())
}

func onUpdateStatusReminder(e eventbus.Event) {
	event := e.(*eventbus.ReminderEvent)

	reminder, err := reminderRepo.GetByID(*event.ReminderID)
	if err == nil {
		reminder.Status = reminderDomain.ReminderStatus(event.ReminderStatus)
		reminder.SentAt = time.Now()

		_ = reminderRepo.Update(reminder)

	}

}

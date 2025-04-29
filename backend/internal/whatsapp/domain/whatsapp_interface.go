package domain

type IWhatsAppHandler interface {
	SendPackageNotification(reminderID uint, phone, name string) error
	SendReservationKeyReminder(reminderID uint, phone, name, hall string) error
}

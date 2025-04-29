package eventbus

import "time"

type PackageCreatedEvent struct {
	PackageID *uint
	Channel   string
	Recipient string
}

type ReservationCreatedEvent struct {
	ReservationID *uint
	StartTime     time.Time
	Channel       string
	Recipient     string
}

type ReminderEvent struct {
	ReminderID     *uint
	ReservationID  *uint
	PackageID      *uint
	ReminderStatus string
	Phone          string
	Name           string
	Hall           string
}

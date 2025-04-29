package domain

import (
	"time"

	"gorm.io/gorm"
)

type ReminderChannel string

const (
	ReminderChannelWhatsApp  ReminderChannel = "WHATSAPP"
	ReminderChannelEmail     ReminderChannel = "EMAIL"
	ReminderChannelSMS       ReminderChannel = "SMS"
	ReminderChannelTelegram  ReminderChannel = "TELEGRAM"
	ReminderChannelInstagram ReminderChannel = "INSTAGRAM"
	ReminderChannelFacebook  ReminderChannel = "FACEBOOK"
	ReminderChannelDiscord   ReminderChannel = "DISCORD"
)

type ReminderStatus string

const (
	ReminderStatusPending   ReminderStatus = "PENDING"
	ReminderStatusSent      ReminderStatus = "SENT"
	ReminderStatusFailed    ReminderStatus = "FAILED"
	ReminderStatusCancelled ReminderStatus = "CANCELLED"
)

type Reminder struct {
	gorm.Model
	Recipient     string          `json:"recipient" gorm:"type:varchar(50);not null"`
	ScheduledAt   time.Time       `json:"scheduled_at"`
	SentAt        time.Time       `json:"sent_at"`
	ReservationID *uint           `json:"reservation_id"`
	PackageID     *uint           `json:"package_id"`
	Channel       ReminderChannel `json:"channel" gorm:"type:varchar(10);not null;default:'WHATSAPP'"`
	Status        ReminderStatus  `json:"status" gorm:"type:varchar(10);not null;default:'PENDING'"`
}

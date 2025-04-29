package domain

import (
	"portarius/internal/eventbus"
	reminderDomain "portarius/internal/reminder/domain"
)

func (s *WhatsAppMessage) PublishReminderSentEvent() {

	eventbus.Publish("UpdateStatusReminder", &eventbus.ReminderEvent{
		ReminderID:     &s.ReminderID,
		ReminderStatus: string(reminderDomain.ReminderStatusSent),
	})
}

func (s *WhatsAppMessage) PublishReminderFailedEvent() {

	eventbus.Publish("UpdateStatusReminder", &eventbus.ReminderEvent{
		ReminderID:     &s.ReminderID,
		ReminderStatus: string(reminderDomain.ReminderStatusFailed),
	})
}

package handler

import (
	"portarius/internal/whatsapp/domain"
)

type WhatsAppHandler struct {
	WhatsAppService *domain.WhatsAppService
}

func NewWhatsAppHandler(service *domain.WhatsAppService) domain.IWhatsAppHandler {
	return &WhatsAppHandler{
		WhatsAppService: service,
	}
}

func (h *WhatsAppHandler) SendReservationKeyReminder(reminderId uint, phone, name, hall string) error {
	message := domain.WhatsAppMessage{
		ReminderID:       reminderId,
		MessagingProduct: "whatsapp",
		To:               phone,
		Type:             "template",
		Template: domain.Template{
			Name: "reservation_key_reminder",
			Language: domain.Language{
				Code: "pt_BR",
			},
			Components: []domain.Component{
				{
					Type: "body",
					Parameters: []domain.Param{
						{
							Type: "text",
							Text: name,
						},
						{
							Type: "text",
							Text: hall,
						},
					},
				},
			},
		},
	}

	return h.WhatsAppService.SendMessage(message)
}

func (h *WhatsAppHandler) SendPackageNotification(reminderId uint, phone, name string) error {
	message := domain.WhatsAppMessage{
		ReminderID:       reminderId,
		MessagingProduct: "whatsapp",
		To:               phone,
		Type:             "template",
		Template: domain.Template{
			Name: "package_notification",
			Language: domain.Language{
				Code: "pt_BR",
			},
			Components: []domain.Component{
				{
					Type: "body",
					Parameters: []domain.Param{
						{
							Type: "text",
							Text: name,
						},
					},
				},
			},
		},
	}

	return h.WhatsAppService.SendMessage(message)
}

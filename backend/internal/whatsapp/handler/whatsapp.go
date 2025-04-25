package handler

import (
	"fmt"
	"portarius/internal/whatsapp/domain"
)

type WhatsAppHandler struct {
	WhatsAppService *domain.WhatsAppService
}

func NewWhatsAppHandler(service *domain.WhatsAppService) *WhatsAppHandler {
	return &WhatsAppHandler{
		WhatsAppService: service,
	}
}

func (h *WhatsAppHandler) SendReservationKeyReminder(phone, name, hall string) error {
	message := domain.WhatsAppMessage{
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

func (h *WhatsAppHandler) SendPackageNotification(phone, name string) error {
	fmt.Println("Sending package notification to:", phone, "for resident:", name)
	message := domain.WhatsAppMessage{
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

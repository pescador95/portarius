package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WhatsAppService struct {
	apiKey     string
	apiBaseURL string
}

type WhatsAppMessage struct {
	To       string `json:"to"`
	Type     string `json:"type"`
	Template struct {
		Name       string      `json:"name"`
		Language   string      `json:"language"`
		Components []Component `json:"components"`
	} `json:"template"`
}

type Component struct {
	Type       string  `json:"type"`
	Parameters []Param `json:"parameters"`
}

type Param struct {
	Type  string `json:"type"`
	Text  string `json:"text,omitempty"`
	Value string `json:"value,omitempty"`
}

func NewWhatsAppService() *WhatsAppService {
	return &WhatsAppService{
		apiKey:     os.Getenv("WHATSAPP_API_KEY"),
		apiBaseURL: "https://graph.facebook.com/v17.0/YOUR_PHONE_NUMBER_ID",
	}
}

func (s *WhatsAppService) SendReservationKeyReminder(phone, name string) error {
	message := WhatsAppMessage{
		To:   phone,
		Type: "template",
		Template: struct {
			Name       string      `json:"name"`
			Language   string      `json:"language"`
			Components []Component `json:"components"`
		}{
			Name:     "reservation_key_reminder",
			Language: "pt_BR",
			Components: []Component{
				{
					Type: "body",
					Parameters: []Param{
						{
							Type: "text",
							Text: name,
						},
					},
				},
			},
		},
	}

	return s.sendMessage(message)
}

func (s *WhatsAppService) SendPackageNotification(phone, name string) error {
	message := WhatsAppMessage{
		To:   phone,
		Type: "template",
		Template: struct {
			Name       string      `json:"name"`
			Language   string      `json:"language"`
			Components []Component `json:"components"`
		}{
			Name:     "package_notification",
			Language: "pt_BR",
			Components: []Component{
				{
					Type: "body",
					Parameters: []Param{
						{
							Type: "text",
							Text: name,
						},
					},
				},
			},
		},
	}

	return s.sendMessage(message)
}

func (s *WhatsAppService) sendMessage(message WhatsAppMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	req, err := http.NewRequest("POST", s.apiBaseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from WhatsApp API: %d", resp.StatusCode)
	}

	return nil
}

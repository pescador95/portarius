package domain

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
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Template         struct {
		Name       string      `json:"name"`
		Language   Language    `json:"language"`
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

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

type Language struct {
	Code string `json:"code"`
}

func NewWhatsAppService() *WhatsAppService {
	phoneId := os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	apiKeyVal := os.Getenv("WHATSAPP_API_KEY")
	return &WhatsAppService{
		apiKey:     apiKeyVal,
		apiBaseURL: "https://graph.facebook.com/v22.0/" + phoneId,
	}
}

func (s *WhatsAppService) SendMessage(message WhatsAppMessage) error {
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

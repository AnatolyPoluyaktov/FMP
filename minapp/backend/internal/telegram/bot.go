package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Bot struct {
	Token  string
	Client *http.Client
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type TelegramWebAppData struct {
	QueryID  string       `json:"query_id"`
	User     TelegramUser `json:"user"`
	AuthDate int64        `json:"auth_date"`
	Hash     string       `json:"hash"`
}

type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type SendMessageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}

func InitializeBot(token string) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is required")
	}

	return &Bot{
		Token: token,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.Token)

	reqBody := SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := b.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

func (b *Bot) ValidateWebAppData(data string) (*TelegramWebAppData, error) {
	// In a real implementation, you would validate the hash here
	// For now, we'll just parse the data
	var webAppData TelegramWebAppData
	if err := json.Unmarshal([]byte(data), &webAppData); err != nil {
		return nil, err
	}

	return &webAppData, nil
}

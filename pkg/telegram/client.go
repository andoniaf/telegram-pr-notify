package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const apiBase = "https://api.telegram.org"

type Client struct {
	botToken string
	chatID   string
	topicID  string
	apiURL   string
}

type sendMessageRequest struct {
	ChatID                string       `json:"chat_id"`
	Text                  string       `json:"text"`
	ParseMode             string       `json:"parse_mode"`
	DisableWebPagePreview bool         `json:"disable_web_page_preview"`
	MessageThreadID       *int         `json:"message_thread_id,omitempty"`
	ReplyMarkup           *replyMarkup `json:"reply_markup,omitempty"`
}

type replyMarkup struct {
	InlineKeyboard [][]inlineButton `json:"inline_keyboard"`
}

type inlineButton struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type apiResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

// NewClient creates a new Telegram client.
func NewClient(botToken, chatID, topicID string) *Client {
	return &Client{
		botToken: botToken,
		chatID:   chatID,
		topicID:  topicID,
		apiURL:   apiBase,
	}
}

// SendMessage sends an HTML message with an optional inline keyboard button.
func (c *Client) SendMessage(text, buttonText, buttonURL string) error {
	req := sendMessageRequest{
		ChatID:                c.chatID,
		Text:                  text,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}

	if c.topicID != "" {
		tid, err := strconv.Atoi(c.topicID)
		if err != nil {
			return fmt.Errorf("invalid topic_id %q: %w", c.topicID, err)
		}
		req.MessageThreadID = &tid
	}

	if buttonText != "" && buttonURL != "" {
		req.ReplyMarkup = &replyMarkup{
			InlineKeyboard: [][]inlineButton{
				{{Text: buttonText, URL: buttonURL}},
			},
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", c.apiURL, c.botToken)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if !apiResp.OK {
		return fmt.Errorf("telegram API error: %s", apiResp.Description)
	}

	return nil
}

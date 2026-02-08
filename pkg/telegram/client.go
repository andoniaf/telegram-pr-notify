package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const apiBase = "https://api.telegram.org"

const telegramMaxMessageLength = 4096
const truncationMarker = "\n\n[message truncated]"

// Button represents an inline keyboard button.
type Button struct {
	Text string
	URL  string
}

type Client struct {
	botToken   string
	chatID     string
	topicID    string
	apiURL     string
	httpClient *http.Client
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
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WithHTTPClient sets a custom http.Client, useful for testing.
func (c *Client) WithHTTPClient(hc *http.Client) *Client {
	c.httpClient = hc
	return c
}

// SendMessage sends an HTML message with optional inline keyboard buttons.
func (c *Client) SendMessage(text string, buttons []Button) error {
	if len([]rune(text)) > telegramMaxMessageLength {
		runes := []rune(text)
		maxLen := telegramMaxMessageLength - len([]rune(truncationMarker))
		text = string(runes[:maxLen]) + truncationMarker
	}

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

	if len(buttons) > 0 {
		row := make([]inlineButton, len(buttons))
		for i, b := range buttons {
			row[i] = inlineButton{Text: b.Text, URL: b.URL}
		}
		req.ReplyMarkup = &replyMarkup{
			InlineKeyboard: [][]inlineButton{row},
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", c.apiURL, c.botToken)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sending request to Telegram API: %w", sanitizeErr(err, c.botToken))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
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

// sanitizeErr removes the bot token from error messages to prevent log leakage.
func sanitizeErr(err error, token string) error {
	if err == nil || token == "" {
		return err
	}
	cleaned := strings.ReplaceAll(err.Error(), token, "[REDACTED]")
	return fmt.Errorf("%s", cleaned)
}

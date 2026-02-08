package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(serverURL string) *Client {
	return &Client{
		botToken:   "test-token",
		chatID:     "-100123",
		apiURL:     serverURL,
		httpClient: &http.Client{},
	}
}

func TestSendMessageSuccess(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	err := client.SendMessage("Hello", []Button{{Text: "View PR", URL: "https://github.com/pr/1"}})
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	if received.ChatID != "-100123" {
		t.Errorf("ChatID = %q, want %q", received.ChatID, "-100123")
	}
	if received.Text != "Hello" {
		t.Errorf("Text = %q, want %q", received.Text, "Hello")
	}
	if received.ParseMode != "HTML" {
		t.Errorf("ParseMode = %q, want %q", received.ParseMode, "HTML")
	}
	if !received.DisableWebPagePreview {
		t.Error("DisableWebPagePreview = false, want true")
	}
	if received.MessageThreadID != nil {
		t.Errorf("MessageThreadID = %v, want nil", received.MessageThreadID)
	}
	if received.ReplyMarkup == nil {
		t.Fatal("ReplyMarkup is nil, want inline keyboard")
	}
	if len(received.ReplyMarkup.InlineKeyboard) != 1 || len(received.ReplyMarkup.InlineKeyboard[0]) != 1 {
		t.Fatal("unexpected inline keyboard structure")
	}
	btn := received.ReplyMarkup.InlineKeyboard[0][0]
	if btn.Text != "View PR" {
		t.Errorf("button text = %q, want %q", btn.Text, "View PR")
	}
	if btn.URL != "https://github.com/pr/1" {
		t.Errorf("button URL = %q, want %q", btn.URL, "https://github.com/pr/1")
	}
}

func TestSendMessageMultipleButtons(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	buttons := []Button{
		{Text: "View PR", URL: "https://github.com/pr/1"},
		{Text: "Issue #15", URL: "https://github.com/repo/issues/15"},
	}
	err := client.SendMessage("Hello", buttons)
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	if received.ReplyMarkup == nil {
		t.Fatal("ReplyMarkup is nil, want inline keyboard")
	}
	if len(received.ReplyMarkup.InlineKeyboard) != 1 {
		t.Fatalf("expected 1 row, got %d", len(received.ReplyMarkup.InlineKeyboard))
	}
	row := received.ReplyMarkup.InlineKeyboard[0]
	if len(row) != 2 {
		t.Fatalf("expected 2 buttons in row, got %d", len(row))
	}
	if row[0].Text != "View PR" {
		t.Errorf("button[0].Text = %q, want %q", row[0].Text, "View PR")
	}
	if row[1].Text != "Issue #15" {
		t.Errorf("button[1].Text = %q, want %q", row[1].Text, "Issue #15")
	}
	if row[1].URL != "https://github.com/repo/issues/15" {
		t.Errorf("button[1].URL = %q, want %q", row[1].URL, "https://github.com/repo/issues/15")
	}
}

func TestSendMessageWithTopicID(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.topicID = "456"

	err := client.SendMessage("Hello", nil)
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	if received.MessageThreadID == nil {
		t.Fatal("MessageThreadID is nil, want 456")
	}
	if *received.MessageThreadID != 456 {
		t.Errorf("MessageThreadID = %d, want %d", *received.MessageThreadID, 456)
	}
	if received.ReplyMarkup != nil {
		t.Error("ReplyMarkup should be nil when no button provided")
	}
}

func TestSendMessageAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"ok": false, "description": "Bad Request: chat not found"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.chatID = "invalid"

	err := client.SendMessage("Hello", nil)
	if err == nil {
		t.Fatal("SendMessage() expected error for API error")
	}

	expected := "telegram API error: Bad Request: chat not found"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestSendMessageInvalidTopicID(t *testing.T) {
	client := &Client{
		botToken:   "test-token",
		chatID:     "-100123",
		topicID:    "not-a-number",
		apiURL:     "http://localhost",
		httpClient: &http.Client{},
	}

	err := client.SendMessage("Hello", nil)
	if err == nil {
		t.Fatal("SendMessage() expected error for invalid topic_id")
	}
}

func TestSendMessageRequestPath(t *testing.T) {
	var requestPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.botToken = "my-secret-token"

	err := client.SendMessage("Hello", nil)
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	expected := "/botmy-secret-token/sendMessage"
	if requestPath != expected {
		t.Errorf("request path = %q, want %q", requestPath, expected)
	}
}

func TestSendMessageTruncatesLongText(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	longText := strings.Repeat("x", 5000)
	err := client.SendMessage(longText, nil)
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	runes := []rune(received.Text)
	if len(runes) > telegramMaxMessageLength {
		t.Errorf("message length = %d runes, want <= %d", len(runes), telegramMaxMessageLength)
	}
	if !strings.HasSuffix(received.Text, truncationMarker) {
		t.Errorf("truncated message should end with %q", truncationMarker)
	}
}

func TestSendMessageDoesNotTruncateShortText(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	shortText := "Hello, World!"
	err := client.SendMessage(shortText, nil)
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	if received.Text != shortText {
		t.Errorf("Text = %q, want %q", received.Text, shortText)
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{}
	client := NewClient("token", "chat", "").WithHTTPClient(customClient)
	if client.httpClient != customClient {
		t.Error("WithHTTPClient did not set the custom http client")
	}
}

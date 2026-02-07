package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessageSuccess(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := &Client{
		botToken: "test-token",
		chatID:   "-100123",
		apiURL:   server.URL,
	}

	err := client.SendMessage("Hello", "View PR", "https://github.com/pr/1")
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

func TestSendMessageWithTopicID(t *testing.T) {
	var received sendMessageRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := &Client{
		botToken: "test-token",
		chatID:   "-100123",
		topicID:  "456",
		apiURL:   server.URL,
	}

	err := client.SendMessage("Hello", "", "")
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

	client := &Client{
		botToken: "test-token",
		chatID:   "invalid",
		apiURL:   server.URL,
	}

	err := client.SendMessage("Hello", "", "")
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
		botToken: "test-token",
		chatID:   "-100123",
		topicID:  "not-a-number",
		apiURL:   "http://localhost",
	}

	err := client.SendMessage("Hello", "", "")
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

	client := &Client{
		botToken: "my-secret-token",
		chatID:   "-100123",
		apiURL:   server.URL,
	}

	err := client.SendMessage("Hello", "", "")
	if err != nil {
		t.Fatalf("SendMessage() error: %v", err)
	}

	expected := "/botmy-secret-token/sendMessage"
	if requestPath != expected {
		t.Errorf("request path = %q, want %q", requestPath, expected)
	}
}

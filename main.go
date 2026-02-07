package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/andoniaf/telegram-pr-notify/pkg/events"
	"github.com/andoniaf/telegram-pr-notify/pkg/templates"
	"github.com/andoniaf/telegram-pr-notify/pkg/telegram"
)

var chatIDPattern = regexp.MustCompile(`^-?\d+$`)

func main() {
	botToken := os.Getenv("INPUT_BOT_TOKEN")
	chatID := os.Getenv("INPUT_CHAT_ID")
	topicID := os.Getenv("INPUT_TOPIC_ID")
	customTemplate := os.Getenv("INPUT_CUSTOM_TEMPLATE")
	eventPayload := os.Getenv("INPUT_EVENT_PAYLOAD")

	if botToken != "" {
		fmt.Fprintf(os.Stdout, "::add-mask::%s\n", botToken)
	}

	if botToken == "" {
		fatal("bot_token is required")
	}
	if chatID == "" {
		fatal("chat_id is required")
	}
	if !chatIDPattern.MatchString(chatID) {
		fatal("chat_id must be a numeric value (e.g., -100123456789)")
	}
	if eventPayload == "" {
		fatal("event_payload is required")
	}

	data, err := events.Parse([]byte(eventPayload))
	if err != nil {
		fatal("parsing event: %v", err)
	}

	message, err := templates.Render(data, customTemplate)
	if err != nil {
		fatal("rendering template: %v", err)
	}

	client := telegram.NewClient(botToken, chatID, topicID)
	if err := client.SendMessage(message, data.ButtonText(), data.RelevantURL()); err != nil {
		fatal("sending message: %v", err)
	}

	fmt.Println("Notification sent successfully")
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "::error::"+format+"\n", args...)
	os.Exit(1)
}

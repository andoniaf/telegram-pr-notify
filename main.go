package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/andoniaf/telegram-pr-notify/pkg/events"
	"github.com/andoniaf/telegram-pr-notify/pkg/telegram"
	"github.com/andoniaf/telegram-pr-notify/pkg/templates"
)

var chatIDPattern = regexp.MustCompile(`^-?\d+$`)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "::error::%v\n", err)
		os.Exit(1)
	}
	fmt.Println("Notification sent successfully")
}

func run() error {
	botToken := os.Getenv("INPUT_BOT_TOKEN")
	chatID := os.Getenv("INPUT_CHAT_ID")
	topicID := os.Getenv("INPUT_TOPIC_ID")
	customTemplate := os.Getenv("INPUT_CUSTOM_TEMPLATE")
	eventPayload := os.Getenv("INPUT_EVENT_PAYLOAD")

	if botToken != "" {
		fmt.Fprintf(os.Stdout, "::add-mask::%s\n", botToken)
	}

	if botToken == "" {
		return fmt.Errorf("bot_token is required")
	}
	if chatID == "" {
		return fmt.Errorf("chat_id is required")
	}
	if !chatIDPattern.MatchString(chatID) {
		return fmt.Errorf("chat_id must be a numeric value (e.g., -100123456789)")
	}
	if eventPayload == "" {
		return fmt.Errorf("event_payload is required")
	}

	data, err := events.Parse([]byte(eventPayload))
	if err != nil {
		return fmt.Errorf("parsing event: %w", err)
	}

	message, err := templates.Render(data, customTemplate)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	buttons := []telegram.Button{{Text: data.ButtonText(), URL: data.RelevantURL()}}
	for _, issue := range data.LinkedIssues() {
		buttons = append(buttons, telegram.Button{Text: issue.Text, URL: issue.URL})
	}

	client := telegram.NewClient(botToken, chatID, topicID)
	if err := client.SendMessage(message, buttons); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

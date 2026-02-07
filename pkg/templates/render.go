package templates

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/andoniaf/telegram-pr-notify/pkg/events"
)

var funcMap = template.FuncMap{
	"truncate": func(s string, max int) string {
		if len(s) <= max {
			return s
		}
		return s[:max] + "..."
	},
}

// Render executes a template against the given data.
// If customTpl is non-empty, it is used as the template string.
// Otherwise, a default template is selected based on event type and action.
func Render(data *events.TemplateData, customTpl string) (string, error) {
	tplStr := customTpl
	if tplStr == "" {
		tplStr = selectDefault(data)
		if tplStr == "" {
			return "", fmt.Errorf("no template for event %s action %s", data.EventName, data.Action)
		}
	}

	tpl, err := template.New("msg").Funcs(funcMap).Parse(tplStr)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func selectDefault(data *events.TemplateData) string {
	if data.IsMerged() {
		return defaultTemplates["pull_request:merged"]
	}

	key := data.EventName + ":" + data.Action
	if data.EventName == "pull_request_review" && data.Action == "submitted" {
		key = data.EventName + ":" + data.Review.State
	}

	return defaultTemplates[key]
}

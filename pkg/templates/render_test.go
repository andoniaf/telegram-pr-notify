package templates

import (
	"strings"
	"testing"

	"github.com/andoniaf/telegram-pr-notify/pkg/events"
)

func samplePRData() *events.TemplateData {
	return &events.TemplateData{
		EventName: "pull_request",
		Action:    "opened",
		Actor: events.User{
			Login:   "octocat",
			HTMLURL: "https://github.com/octocat",
		},
		Repo: events.Repository{
			FullName: "octocat/Hello-World",
			HTMLURL:  "https://github.com/octocat/Hello-World",
		},
		PR: events.PullRequest{
			Number:  42,
			Title:   "Add new feature",
			HTMLURL: "https://github.com/octocat/Hello-World/pull/42",
			Head:    events.Branch{Ref: "feature-branch"},
			Base:    events.Branch{Ref: "main"},
		},
	}
}

func TestRenderDefaultPROpened(t *testing.T) {
	data := samplePRData()
	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	expectations := []string{
		"New Pull Request",
		"#42",
		"Add new feature",
		"octocat",
		"octocat/Hello-World",
		"feature-branch",
		"main",
	}

	for _, exp := range expectations {
		if !strings.Contains(result, exp) {
			t.Errorf("result missing %q:\n%s", exp, result)
		}
	}
}

func TestRenderDefaultMerged(t *testing.T) {
	data := samplePRData()
	data.Action = "closed"
	data.PR.Merged = true

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(result, "Merged") {
		t.Errorf("result missing 'Merged':\n%s", result)
	}
}

func TestRenderDefaultClosed(t *testing.T) {
	data := samplePRData()
	data.Action = "closed"
	data.PR.Merged = false

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(result, "Closed") {
		t.Errorf("result missing 'Closed':\n%s", result)
	}
}

func TestRenderReviewApproved(t *testing.T) {
	data := samplePRData()
	data.EventName = "pull_request_review"
	data.Action = "submitted"
	data.Review = events.Review{
		State: "approved",
		Body:  "LGTM!",
	}

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(result, "Approved") {
		t.Errorf("result missing 'Approved':\n%s", result)
	}
	if !strings.Contains(result, "LGTM!") {
		t.Errorf("result missing review body:\n%s", result)
	}
}

func TestRenderReviewChangesRequested(t *testing.T) {
	data := samplePRData()
	data.EventName = "pull_request_review"
	data.Action = "submitted"
	data.Review = events.Review{
		State: "changes_requested",
		Body:  "Please fix the tests",
	}

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(result, "Changes Requested") {
		t.Errorf("result missing 'Changes Requested':\n%s", result)
	}
}

func TestRenderCustomTemplate(t *testing.T) {
	data := samplePRData()
	custom := "PR #{{.PR.Number}} by {{.Actor.Login}}"

	result, err := Render(data, custom)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	expected := "PR #42 by octocat"
	if result != expected {
		t.Errorf("result = %q, want %q", result, expected)
	}
}

func TestRenderHTMLEscaping(t *testing.T) {
	data := samplePRData()
	data.PR.Title = "Fix <script>alert('xss')</script>"

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(result, "<script>") {
		t.Errorf("result contains unescaped HTML:\n%s", result)
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	data := samplePRData()
	_, err := Render(data, "{{.Invalid")
	if err == nil {
		t.Error("Render() expected error for invalid template")
	}
}

func TestRenderNoTemplateForEvent(t *testing.T) {
	data := &events.TemplateData{
		EventName: "pull_request",
		Action:    "unknown_action",
	}
	_, err := Render(data, "")
	if err == nil {
		t.Error("Render() expected error for unknown action")
	}
}

func TestRenderTruncate(t *testing.T) {
	data := samplePRData()
	data.EventName = "pull_request_review"
	data.Action = "submitted"
	data.Review = events.Review{
		State: "approved",
		Body:  strings.Repeat("a", 600),
	}

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if len(result) > 1000 && !strings.Contains(result, "...") {
		t.Error("expected truncation with '...' for long review body")
	}
}

func TestRenderReviewCommentCreated(t *testing.T) {
	data := samplePRData()
	data.EventName = "pull_request_review_comment"
	data.Action = "created"
	data.Comment = events.Comment{
		Body: "This needs a fix",
		Path: "pkg/main.go",
	}

	result, err := Render(data, "")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(result, "pkg/main.go") {
		t.Errorf("result missing file path:\n%s", result)
	}
	if !strings.Contains(result, "This needs a fix") {
		t.Errorf("result missing comment body:\n%s", result)
	}
}

package events

import (
	"os"
	"testing"
)

func TestParsePullRequestOpened(t *testing.T) {
	payload, err := os.ReadFile("../../testdata/pull_request_opened.json")
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	data, err := Parse(payload)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if data.EventName != "pull_request" {
		t.Errorf("EventName = %q, want %q", data.EventName, "pull_request")
	}
	if data.Action != "opened" {
		t.Errorf("Action = %q, want %q", data.Action, "opened")
	}
	if data.PR.Number != 42 {
		t.Errorf("PR.Number = %d, want %d", data.PR.Number, 42)
	}
	if data.PR.Title != "Add new feature" {
		t.Errorf("PR.Title = %q, want %q", data.PR.Title, "Add new feature")
	}
	if data.Actor.Login != "octocat" {
		t.Errorf("Actor.Login = %q, want %q", data.Actor.Login, "octocat")
	}
	if data.Repo.FullName != "octocat/Hello-World" {
		t.Errorf("Repo.FullName = %q, want %q", data.Repo.FullName, "octocat/Hello-World")
	}
	if data.PR.Head.Ref != "feature-branch" {
		t.Errorf("PR.Head.Ref = %q, want %q", data.PR.Head.Ref, "feature-branch")
	}
	if data.PR.Base.Ref != "main" {
		t.Errorf("PR.Base.Ref = %q, want %q", data.PR.Base.Ref, "main")
	}
}

func TestParsePullRequestClosedMerged(t *testing.T) {
	payload, err := os.ReadFile("../../testdata/pull_request_closed_merged.json")
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	data, err := Parse(payload)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if data.Action != "closed" {
		t.Errorf("Action = %q, want %q", data.Action, "closed")
	}
	if !data.IsMerged() {
		t.Error("IsMerged() = false, want true")
	}
}

func TestParseReviewApproved(t *testing.T) {
	payload, err := os.ReadFile("../../testdata/pull_request_review_approved.json")
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	data, err := Parse(payload)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if data.EventName != "pull_request_review" {
		t.Errorf("EventName = %q, want %q", data.EventName, "pull_request_review")
	}
	if data.Action != "submitted" {
		t.Errorf("Action = %q, want %q", data.Action, "submitted")
	}
	if data.Review.State != "approved" {
		t.Errorf("Review.State = %q, want %q", data.Review.State, "approved")
	}
	if data.Review.Body != "Looks good to me!" {
		t.Errorf("Review.Body = %q, want %q", data.Review.Body, "Looks good to me!")
	}
	if data.Actor.Login != "reviewer" {
		t.Errorf("Actor.Login = %q, want %q", data.Actor.Login, "reviewer")
	}
}

func TestParseReviewComment(t *testing.T) {
	payload, err := os.ReadFile("../../testdata/pull_request_review_comment.json")
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}

	data, err := Parse(payload)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if data.EventName != "pull_request_review_comment" {
		t.Errorf("EventName = %q, want %q", data.EventName, "pull_request_review_comment")
	}
	if data.Comment.Path != "src/main.go" {
		t.Errorf("Comment.Path = %q, want %q", data.Comment.Path, "src/main.go")
	}
	if data.Comment.Body != "Consider using a constant here" {
		t.Errorf("Comment.Body = %q, want %q", data.Comment.Body, "Consider using a constant here")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	_, err := Parse([]byte("not json"))
	if err == nil {
		t.Error("Parse() expected error for invalid JSON")
	}
}

func TestParseUnsupportedEvent(t *testing.T) {
	payload := []byte(`{"event_name": "push", "event": {}}`)
	_, err := Parse(payload)
	if err == nil {
		t.Error("Parse() expected error for unsupported event")
	}
}

func TestParseMissingEvent(t *testing.T) {
	payload := []byte(`{"event_name": "pull_request"}`)
	_, err := Parse(payload)
	if err == nil {
		t.Error("Parse() expected error for missing event payload")
	}
}

func TestIsMergedFalseForOpenPR(t *testing.T) {
	data := &TemplateData{
		EventName: "pull_request",
		Action:    "opened",
		PR:        PullRequest{Merged: false},
	}
	if data.IsMerged() {
		t.Error("IsMerged() = true, want false for opened PR")
	}
}

func TestRelevantURL(t *testing.T) {
	tests := []struct {
		name string
		data TemplateData
		want string
	}{
		{
			name: "pull_request returns PR URL",
			data: TemplateData{
				EventName: "pull_request",
				PR:        PullRequest{HTMLURL: "https://github.com/pr/1"},
			},
			want: "https://github.com/pr/1",
		},
		{
			name: "review returns review URL",
			data: TemplateData{
				EventName: "pull_request_review",
				PR:        PullRequest{HTMLURL: "https://github.com/pr/1"},
				Review:    Review{HTMLURL: "https://github.com/pr/1#review-1"},
			},
			want: "https://github.com/pr/1#review-1",
		},
		{
			name: "review comment returns comment URL",
			data: TemplateData{
				EventName: "pull_request_review_comment",
				PR:        PullRequest{HTMLURL: "https://github.com/pr/1"},
				Comment:   Comment{HTMLURL: "https://github.com/pr/1#comment-1"},
			},
			want: "https://github.com/pr/1#comment-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.RelevantURL(); got != tt.want {
				t.Errorf("RelevantURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestButtonText(t *testing.T) {
	tests := []struct {
		eventName string
		want      string
	}{
		{"pull_request", "View Pull Request"},
		{"pull_request_review", "View Review"},
		{"pull_request_review_comment", "View Comment"},
	}

	for _, tt := range tests {
		t.Run(tt.eventName, func(t *testing.T) {
			data := &TemplateData{EventName: tt.eventName}
			if got := data.ButtonText(); got != tt.want {
				t.Errorf("ButtonText() = %q, want %q", got, tt.want)
			}
		})
	}
}

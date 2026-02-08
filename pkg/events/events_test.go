package events

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		fixture    string
		wantEvent  string
		wantAction string
		check      func(t *testing.T, data *TemplateData)
	}{
		{
			name:       "pull_request opened",
			fixture:    "../../testdata/pull_request_opened.json",
			wantEvent:  "pull_request",
			wantAction: "opened",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if data.PR.Number != 42 {
					t.Errorf("PR.Number = %d, want 42", data.PR.Number)
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
			},
		},
		{
			name:       "pull_request closed (not merged)",
			fixture:    "../../testdata/pull_request_closed.json",
			wantEvent:  "pull_request",
			wantAction: "closed",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if data.IsMerged() {
					t.Error("IsMerged() = true, want false")
				}
			},
		},
		{
			name:       "pull_request closed merged",
			fixture:    "../../testdata/pull_request_closed_merged.json",
			wantEvent:  "pull_request",
			wantAction: "closed",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if !data.IsMerged() {
					t.Error("IsMerged() = false, want true")
				}
			},
		},
		{
			name:       "pull_request reopened",
			fixture:    "../../testdata/pull_request_reopened.json",
			wantEvent:  "pull_request",
			wantAction: "reopened",
		},
		{
			name:       "pull_request synchronize",
			fixture:    "../../testdata/pull_request_synchronize.json",
			wantEvent:  "pull_request",
			wantAction: "synchronize",
		},
		{
			name:       "pull_request ready_for_review",
			fixture:    "../../testdata/pull_request_ready_for_review.json",
			wantEvent:  "pull_request",
			wantAction: "ready_for_review",
		},
		{
			name:       "pull_request converted_to_draft",
			fixture:    "../../testdata/pull_request_converted_to_draft.json",
			wantEvent:  "pull_request",
			wantAction: "converted_to_draft",
		},
		{
			name:       "review approved",
			fixture:    "../../testdata/pull_request_review_approved.json",
			wantEvent:  "pull_request_review",
			wantAction: "approved",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if data.Review.State != "approved" {
					t.Errorf("Review.State = %q, want %q", data.Review.State, "approved")
				}
				if data.Review.Body != "Looks good to me!" {
					t.Errorf("Review.Body = %q, want %q", data.Review.Body, "Looks good to me!")
				}
				if data.Actor.Login != "reviewer" {
					t.Errorf("Actor.Login = %q, want %q", data.Actor.Login, "reviewer")
				}
			},
		},
		{
			name:       "review changes_requested",
			fixture:    "../../testdata/pull_request_review_changes_requested.json",
			wantEvent:  "pull_request_review",
			wantAction: "changes_requested",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if data.Review.State != "changes_requested" {
					t.Errorf("Review.State = %q, want %q", data.Review.State, "changes_requested")
				}
				if data.Review.Body != "Please fix the tests" {
					t.Errorf("Review.Body = %q, want %q", data.Review.Body, "Please fix the tests")
				}
			},
		},
		{
			name:       "review commented",
			fixture:    "../../testdata/pull_request_review_commented.json",
			wantEvent:  "pull_request_review",
			wantAction: "commented",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if data.Review.State != "commented" {
					t.Errorf("Review.State = %q, want %q", data.Review.State, "commented")
				}
				if data.Review.Body != "Looks interesting, a few thoughts..." {
					t.Errorf("Review.Body = %q, want %q", data.Review.Body, "Looks interesting, a few thoughts...")
				}
			},
		},
		{
			name:       "review comment",
			fixture:    "../../testdata/pull_request_review_comment.json",
			wantEvent:  "pull_request_review_comment",
			wantAction: "created",
			check: func(t *testing.T, data *TemplateData) {
				t.Helper()
				if data.Comment.Path != "src/main.go" {
					t.Errorf("Comment.Path = %q, want %q", data.Comment.Path, "src/main.go")
				}
				if data.Comment.Body != "Consider using a constant here" {
					t.Errorf("Comment.Body = %q, want %q", data.Comment.Body, "Consider using a constant here")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := os.ReadFile(tt.fixture)
			if err != nil {
				t.Fatalf("reading testdata: %v", err)
			}
			data, err := Parse(payload)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			if data.EventName != tt.wantEvent {
				t.Errorf("EventName = %q, want %q", data.EventName, tt.wantEvent)
			}
			if data.Action != tt.wantAction {
				t.Errorf("Action = %q, want %q", data.Action, tt.wantAction)
			}
			if tt.check != nil {
				tt.check(t, data)
			}
		})
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
			name: "review with empty HTMLURL falls back to PR URL",
			data: TemplateData{
				EventName: "pull_request_review",
				PR:        PullRequest{HTMLURL: "https://github.com/pr/1"},
				Review:    Review{HTMLURL: ""},
			},
			want: "https://github.com/pr/1",
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
		{
			name: "review comment with empty HTMLURL falls back to PR URL",
			data: TemplateData{
				EventName: "pull_request_review_comment",
				PR:        PullRequest{HTMLURL: "https://github.com/pr/1"},
				Comment:   Comment{HTMLURL: ""},
			},
			want: "https://github.com/pr/1",
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

func TestLinkedIssues(t *testing.T) {
	base := TemplateData{
		Repo: Repository{HTMLURL: "https://github.com/octocat/Hello-World"},
	}

	tests := []struct {
		name     string
		body     string
		wantLen  int
		wantText []string
		wantURL  []string
	}{
		{
			name:    "no issues",
			body:    "Just a regular PR body",
			wantLen: 0,
		},
		{
			name:     "single closes",
			body:     "Closes #15",
			wantLen:  1,
			wantText: []string{"Issue #15"},
			wantURL:  []string{"https://github.com/octocat/Hello-World/issues/15"},
		},
		{
			name:     "multiple issues",
			body:     "Fixes #10 and resolves #20",
			wantLen:  2,
			wantText: []string{"Issue #10", "Issue #20"},
			wantURL: []string{
				"https://github.com/octocat/Hello-World/issues/10",
				"https://github.com/octocat/Hello-World/issues/20",
			},
		},
		{
			name:     "duplicates are removed",
			body:     "Fixes #42\nAlso fixes #42",
			wantLen:  1,
			wantText: []string{"Issue #42"},
			wantURL:  []string{"https://github.com/octocat/Hello-World/issues/42"},
		},
		{
			name:     "case insensitive",
			body:     "CLOSES #5\nFIXES #6\nRESOLVES #7",
			wantLen:  3,
			wantText: []string{"Issue #5", "Issue #6", "Issue #7"},
		},
		{
			name:     "all keyword variants",
			body:     "close #1\ncloses #2\nclosed #3\nfix #4\nfixes #5\nfixed #6\nresolve #7\nresolves #8\nresolved #9",
			wantLen:  9,
			wantText: []string{"Issue #1", "Issue #2", "Issue #3", "Issue #4", "Issue #5", "Issue #6", "Issue #7", "Issue #8", "Issue #9"},
		},
		{
			name:     "colon variant",
			body:     "Closes: #11\nFixes: #12",
			wantLen:  2,
			wantText: []string{"Issue #11", "Issue #12"},
			wantURL: []string{
				"https://github.com/octocat/Hello-World/issues/11",
				"https://github.com/octocat/Hello-World/issues/12",
			},
		},
		{
			name:     "refs keyword",
			body:     "Refs #113\nRef #50",
			wantLen:  2,
			wantText: []string{"Issue #113", "Issue #50"},
			wantURL: []string{
				"https://github.com/octocat/Hello-World/issues/113",
				"https://github.com/octocat/Hello-World/issues/50",
			},
		},
		{
			name:    "empty body",
			body:    "",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := base
			data.PR = PullRequest{Body: tt.body}
			issues := data.LinkedIssues()
			if len(issues) != tt.wantLen {
				t.Fatalf("LinkedIssues() returned %d issues, want %d", len(issues), tt.wantLen)
			}
			for i, issue := range issues {
				if i < len(tt.wantText) && issue.Text != tt.wantText[i] {
					t.Errorf("issue[%d].Text = %q, want %q", i, issue.Text, tt.wantText[i])
				}
				if tt.wantURL != nil && i < len(tt.wantURL) && issue.URL != tt.wantURL[i] {
					t.Errorf("issue[%d].URL = %q, want %q", i, issue.URL, tt.wantURL[i])
				}
			}
		})
	}
}

func TestLinkedIssuesFromFixture(t *testing.T) {
	payload, err := os.ReadFile("../../testdata/pull_request_opened_with_issue.json")
	if err != nil {
		t.Fatalf("reading testdata: %v", err)
	}
	data, err := Parse(payload)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	issues := data.LinkedIssues()
	if len(issues) != 1 {
		t.Fatalf("LinkedIssues() = %d issues, want 1", len(issues))
	}
	if issues[0].Text != "Issue #15" {
		t.Errorf("issue text = %q, want %q", issues[0].Text, "Issue #15")
	}
	if issues[0].URL != "https://github.com/octocat/Hello-World/issues/15" {
		t.Errorf("issue URL = %q, want %q", issues[0].URL, "https://github.com/octocat/Hello-World/issues/15")
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

package events

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// GitHubContext represents the top-level structure of toJSON(github).
type GitHubContext struct {
	EventName string          `json:"event_name"`
	Event     json.RawMessage `json:"event"`
}

type User struct {
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
}

type Repository struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

type PullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	Body    string `json:"body"`
	Draft   bool   `json:"draft"`
	Merged  bool   `json:"merged"`
	User    User   `json:"user"`
	Base    Branch `json:"base"`
	Head    Branch `json:"head"`
}

type Branch struct {
	Ref string `json:"ref"`
}

type Review struct {
	State   string `json:"state"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    User   `json:"user"`
}

type Comment struct {
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	Path    string `json:"path"`
	User    User   `json:"user"`
}

type pullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

type reviewEvent struct {
	Action      string      `json:"action"`
	Review      Review      `json:"review"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

type reviewCommentEvent struct {
	Action      string      `json:"action"`
	Comment     Comment     `json:"comment"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}

// TemplateData is the flattened view passed to templates.
type TemplateData struct {
	EventName string
	Action    string
	Actor     User
	Repo      Repository
	PR        PullRequest
	Review    Review
	Comment   Comment
}

var linkedIssuePattern = regexp.MustCompile(`(?i)\b(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?|refs?):?\s+#(\d+)\b`)

// IssueLink represents a linked issue parsed from a PR body.
type IssueLink struct {
	Text string
	URL  string
}

// LinkedIssues parses GitHub closing keywords from the PR body and returns
// deduplicated issue links for the same repository.
func (d *TemplateData) LinkedIssues() []IssueLink {
	if d.PR.Body == "" || d.Repo.HTMLURL == "" {
		return nil
	}

	matches := linkedIssuePattern.FindAllStringSubmatch(d.PR.Body, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[int]bool)
	var links []IssueLink
	for _, m := range matches {
		num, err := strconv.Atoi(m[1])
		if err != nil || seen[num] {
			continue
		}
		seen[num] = true
		links = append(links, IssueLink{
			Text: fmt.Sprintf("Issue #%d", num),
			URL:  fmt.Sprintf("%s/issues/%d", d.Repo.HTMLURL, num),
		})
	}
	return links
}

// IsMerged returns true if the PR was closed via merge.
func (d *TemplateData) IsMerged() bool {
	return d.EventName == "pull_request" && d.Action == "closed" && d.PR.Merged
}

// RelevantURL returns the most relevant URL for the event.
func (d *TemplateData) RelevantURL() string {
	switch d.EventName {
	case "pull_request_review":
		if d.Review.HTMLURL != "" {
			return d.Review.HTMLURL
		}
	case "pull_request_review_comment":
		if d.Comment.HTMLURL != "" {
			return d.Comment.HTMLURL
		}
	}
	return d.PR.HTMLURL
}

// ButtonText returns a label for the inline keyboard button.
func (d *TemplateData) ButtonText() string {
	switch d.EventName {
	case "pull_request_review":
		return "View Review"
	case "pull_request_review_comment":
		return "View Comment"
	default:
		return "View Pull Request"
	}
}

// Parse parses a GitHub context JSON payload into TemplateData.
func Parse(payload []byte) (*TemplateData, error) {
	var ctx GitHubContext
	if err := json.Unmarshal(payload, &ctx); err != nil {
		return nil, fmt.Errorf("parsing github context: %w", err)
	}

	if ctx.Event == nil {
		return nil, fmt.Errorf("missing event payload")
	}

	switch ctx.EventName {
	case "pull_request":
		return parsePullRequest(ctx.Event)
	case "pull_request_review":
		return parseReview(ctx.Event)
	case "pull_request_review_comment":
		return parseReviewComment(ctx.Event)
	default:
		return nil, fmt.Errorf("unsupported event: %s", ctx.EventName)
	}
}

func parsePullRequest(raw json.RawMessage) (*TemplateData, error) {
	var e pullRequestEvent
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("parsing pull_request event: %w", err)
	}
	return &TemplateData{
		EventName: "pull_request",
		Action:    e.Action,
		Actor:     e.Sender,
		Repo:      e.Repository,
		PR:        e.PullRequest,
	}, nil
}

func parseReview(raw json.RawMessage) (*TemplateData, error) {
	var e reviewEvent
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("parsing pull_request_review event: %w", err)
	}

	action := e.Action
	if action == "submitted" {
		action = e.Review.State
	}

	return &TemplateData{
		EventName: "pull_request_review",
		Action:    action,
		Actor:     e.Sender,
		Repo:      e.Repository,
		PR:        e.PullRequest,
		Review:    e.Review,
	}, nil
}

func parseReviewComment(raw json.RawMessage) (*TemplateData, error) {
	var e reviewCommentEvent
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("parsing pull_request_review_comment event: %w", err)
	}
	return &TemplateData{
		EventName: "pull_request_review_comment",
		Action:    e.Action,
		Actor:     e.Sender,
		Repo:      e.Repository,
		PR:        e.PullRequest,
		Comment:   e.Comment,
	}, nil
}

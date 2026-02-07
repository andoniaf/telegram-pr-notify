# Telegram PR Notify

A GitHub Action that sends Telegram notifications for Pull Request events. Built with Go, zero external dependencies.

## Features

- Notifications for PR opens, closes, merges, reopens, updates, and draft changes
- Review notifications (approved, changes requested, commented)
- Review comment notifications
- Customizable message templates using Go `html/template` syntax
- Telegram forum/topic support
- Inline keyboard button linking to the PR/review/comment
- Minimal Docker image (distroless)

## Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `bot_token` | Yes | - | Telegram Bot API token |
| `chat_id` | Yes | - | Telegram chat ID |
| `topic_id` | No | `""` | Telegram forum topic/thread ID |
| `custom_template` | No | `""` | Go template string to override default message |
| `event_payload` | No | `${{ toJSON(github) }}` | GitHub context JSON payload. Automatically populated by GitHub Actions. For the payload schema, see [GitHub Actions context documentation](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs#github-context). Override only for testing. |

## Supported Events

| Event | Actions |
|-------|---------|
| `pull_request` | `opened`, `closed` (merged detection), `reopened`, `synchronize`, `ready_for_review`, `converted_to_draft` |
| `pull_request_review` | `submitted` (approved, changes_requested, commented) |
| `pull_request_review_comment` | `created` |

> **Note:** The `pull_request_review` event only has the `submitted` action type. To filter by review state (e.g., only approvals), add a condition to your workflow step:
>
> ```yaml
> - uses: andoniaf/telegram-pr-notify@v1
>   if: github.event.review.state == 'approved'
>   with:
>     bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
>     chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
> ```

## Usage Examples

### All PR Events

```yaml
name: PR Notifications
on:
  pull_request:
    types: [opened, closed, reopened, synchronize, ready_for_review, converted_to_draft]
  pull_request_review:
    types: [submitted]
  pull_request_review_comment:
    types: [created]

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - uses: andoniaf/telegram-pr-notify@v1
        with:
          bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
```

### Only New PRs (with Topic Support)

```yaml
name: New PR Notifications
on:
  pull_request:
    types: [opened]

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - uses: andoniaf/telegram-pr-notify@v1
        with:
          bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
          topic_id: "12345"
```

### PRs + Reviews with Custom Template

```yaml
name: PR + Review Notifications
on:
  pull_request:
    types: [opened, closed]
  pull_request_review:
    types: [submitted]

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - uses: andoniaf/telegram-pr-notify@v1
        with:
          bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
          custom_template: |
            [{{.Repo.FullName}}] {{.Actor.Login}}: {{.Action}} PR #{{.PR.Number}} {{.PR.Title}}
```

## Common Workflow Patterns

### Only notify for PRs targeting `main`

```yaml
on:
  pull_request:
    types: [opened, closed]
    branches: [main]
```

### Skip draft PRs

```yaml
- uses: andoniaf/telegram-pr-notify@v1
  if: github.event.pull_request.draft == false
  with:
    bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
    chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
```

### Skip bot PRs (e.g., Dependabot)

```yaml
- uses: andoniaf/telegram-pr-notify@v1
  if: github.actor != 'dependabot[bot]'
  with:
    bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
    chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
```

### Different templates per event type

```yaml
jobs:
  notify-pr:
    if: github.event_name == 'pull_request'
    runs-on: ubuntu-latest
    steps:
      - uses: andoniaf/telegram-pr-notify@v1
        with:
          bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
          custom_template: "PR {{.Action}}: #{{.PR.Number}} {{.PR.Title}}"
  notify-review:
    if: github.event_name == 'pull_request_review'
    runs-on: ubuntu-latest
    steps:
      - uses: andoniaf/telegram-pr-notify@v1
        with:
          bot_token: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          chat_id: ${{ secrets.TELEGRAM_CHAT_ID }}
          custom_template: "Review {{.Review.State}} on #{{.PR.Number}}"
```

## Custom Templates

Templates use Go [`html/template`](https://pkg.go.dev/html/template) syntax. Since `html/template` is used, user-generated content (PR titles, usernames, etc.) is automatically HTML-escaped for safe use in Telegram HTML messages.

> **Security:** The `custom_template` input should only contain trusted input defined by the workflow author. Do not pass user-controlled data (e.g., PR body, branch names) into the template string itself. Note that `PR.Body` may contain sensitive content submitted by external contributors; use `{{truncate .PR.Body N}}` to limit its length and avoid leaking large amounts of text into notifications.

### Available Fields

| Variable | Type | Description |
|----------|------|-------------|
| `{{.EventName}}` | string | GitHub event name (`pull_request`, `pull_request_review`, `pull_request_review_comment`) |
| `{{.Action}}` | string | Event action (`opened`, `closed`, `submitted`, etc.) |
| `{{.Actor.Login}}` | string | User who triggered the event |
| `{{.Actor.HTMLURL}}` | string | URL to the actor's profile |
| `{{.Repo.FullName}}` | string | Repository full name (e.g., `owner/repo`) |
| `{{.Repo.HTMLURL}}` | string | URL to the repository |
| `{{.PR.Number}}` | int | Pull request number |
| `{{.PR.Title}}` | string | Pull request title |
| `{{.PR.HTMLURL}}` | string | URL to the pull request |
| `{{.PR.Body}}` | string | Pull request body/description |
| `{{.PR.Draft}}` | bool | Whether the PR is a draft |
| `{{.PR.Merged}}` | bool | Whether the PR was merged |
| `{{.PR.Head.Ref}}` | string | Source branch name |
| `{{.PR.Base.Ref}}` | string | Target branch name |
| `{{.Review.State}}` | string | Review state (`approved`, `changes_requested`, `commented`) |
| `{{.Review.Body}}` | string | Review body text |
| `{{.Comment.Body}}` | string | Review comment body text |
| `{{.Comment.Path}}` | string | File path of the review comment |

### Available Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `{{.IsMerged}}` | bool | `true` when a `pull_request` / `closed` event has `PR.Merged == true` |

### Template Functions

| Function | Description |
|----------|-------------|
| `{{truncate .Field 100}}` | Truncate a string to a maximum length, appending `...` if truncated |

### Conditional Examples

Merged vs. closed:

```yaml
custom_template: |
  {{if .IsMerged}}Merged{{else if eq .Action "closed"}}Closed{{else}}{{.Action}}{{end}} PR #{{.PR.Number}} {{.PR.Title}}
```

Include review body only when present:

```yaml
custom_template: |
  Review on #{{.PR.Number}}: {{.Review.State}}
  {{- if .Review.Body}}
  Comment: {{truncate .Review.Body 200}}
  {{- end}}
```

### Default Template Example

The default template for a new PR looks like this:

```
<emoji> <b>New Pull Request</b>
<a href="{{.PR.HTMLURL}}">#{{.PR.Number}}</a> {{.PR.Title}}
{{.PR.Head.Ref}} -> {{.PR.Base.Ref}}

by <a href="{{.Actor.HTMLURL}}">{{.Actor.Login}}</a> in <b>{{.Repo.FullName}}</b>
```

See [`pkg/templates/defaults.go`](pkg/templates/defaults.go) for all default templates.

## Setup

### 1. Create a Telegram Bot

1. Message [@BotFather](https://t.me/BotFather) on Telegram
2. Send `/newbot` and follow the instructions
3. Copy the bot token

### 2. Get Your Chat ID

1. Add the bot to your group/channel
2. Send a message in the group
3. Visit `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`
4. Find the `chat.id` field in the response

> **Note:** Group chat IDs are negative numbers. Supergroup IDs start with `-100` (e.g., `-1001234567890`). Private chat IDs are positive numbers. Make sure to include the full number including the minus sign.

### 3. Get Topic ID (Optional)

For Telegram groups with topics/forums enabled:

1. Open the topic in Telegram Web
2. The topic ID is in the URL: `https://web.telegram.org/.../<TOPIC_ID>`

### 4. Add Secrets

Add `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` as [repository secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets).

## Troubleshooting

| Problem | Cause | Solution |
|---------|-------|----------|
| `telegram API error: Bad Request: chat not found` | Bot not added to the group, or incorrect `chat_id` | Ensure the bot is a member of the group. Verify `chat_id` using the `/getUpdates` API (see Setup). |
| `telegram API error: Bad Request: message thread not found` | `topic_id` does not exist or topics are not enabled | Verify the topic exists and that the group has topics/forums enabled. |
| `telegram API error: Forbidden: bot was blocked by the user` | Bot lacks permissions or was removed | Re-add the bot to the group and ensure it has permission to send messages. |
| `parsing template: ...` error | Invalid Go template syntax in `custom_template` | Check your template syntax against the [Go template docs](https://pkg.go.dev/html/template). Common issues: unmatched `{{`, missing closing `{{end}}`, referencing non-existent fields. |
| `unsupported event: <name>` | Workflow triggers an event this action does not handle | Only `pull_request`, `pull_request_review`, and `pull_request_review_comment` events are supported. |
| `no template for event <name> action <action>` | Valid event but unrecognized action | Check the Supported Events table. Ensure your workflow `types` filter matches supported actions. |

## Versioning

This action follows [semantic versioning](https://semver.org/). Use a major version tag for stability:

```yaml
uses: andoniaf/telegram-pr-notify@v1
```

The `v1` tag always points to the latest `v1.x.x` release. You can also pin to a specific version (e.g., `@v1.0.0`) for maximum reproducibility.

## License

[MIT](LICENSE)

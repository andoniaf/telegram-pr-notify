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
| `event_payload` | No | `${{ toJSON(github) }}` | GitHub context JSON payload |

## Supported Events

| Event | Actions |
|-------|---------|
| `pull_request` | `opened`, `closed` (merged detection), `reopened`, `synchronize`, `ready_for_review`, `converted_to_draft` |
| `pull_request_review` | `submitted` (approved, changes_requested, commented) |
| `pull_request_review_comment` | `created` |

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

## Custom Template Variables

Templates use Go `html/template` syntax. Available fields:

| Variable | Description |
|----------|-------------|
| `{{.EventName}}` | GitHub event name (`pull_request`, `pull_request_review`, `pull_request_review_comment`) |
| `{{.Action}}` | Event action (`opened`, `closed`, `submitted`, etc.) |
| `{{.Actor.Login}}` | User who triggered the event |
| `{{.Actor.HTMLURL}}` | URL to the actor's profile |
| `{{.Repo.FullName}}` | Repository full name (e.g., `owner/repo`) |
| `{{.Repo.HTMLURL}}` | URL to the repository |
| `{{.PR.Number}}` | Pull request number |
| `{{.PR.Title}}` | Pull request title |
| `{{.PR.HTMLURL}}` | URL to the pull request |
| `{{.PR.Body}}` | Pull request body/description |
| `{{.PR.Draft}}` | Whether the PR is a draft |
| `{{.PR.Merged}}` | Whether the PR was merged |
| `{{.PR.Head.Ref}}` | Source branch name |
| `{{.PR.Base.Ref}}` | Target branch name |
| `{{.Review.State}}` | Review state (`approved`, `changes_requested`, `commented`) |
| `{{.Review.Body}}` | Review body text |
| `{{.Comment.Body}}` | Review comment body text |
| `{{.Comment.Path}}` | File path of the review comment |

Template functions:
- `{{truncate .Field 100}}` — Truncate a string to a maximum length

Since `html/template` is used, user-generated content (PR titles, usernames, etc.) is automatically HTML-escaped.

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

### 3. Get Topic ID (Optional)

For Telegram groups with topics/forums enabled:

1. Open the topic in Telegram Web
2. The topic ID is in the URL: `https://web.telegram.org/.../<TOPIC_ID>`

### 4. Add Secrets

Add `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` as [repository secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets).

## License

[MIT](LICENSE)

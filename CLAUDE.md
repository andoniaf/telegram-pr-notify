# Telegram PR Notify

GitHub Action (Go + Docker) that sends Telegram notifications for Pull Request events.

## Architecture

- **Language**: Go 1.23, zero external dependencies (stdlib only)
- **Runtime**: Docker container action (multi-stage build, distroless)
- **API**: Direct HTTP POST to Telegram Bot API (`/bot{token}/sendMessage`)
- **Templates**: `html/template` for auto-escaping user content

## Project Structure

```
main.go                      # Entry point: env → parse → render → send
pkg/events/events.go         # GitHub context parsing, TemplateData model
pkg/templates/defaults.go    # 11 default HTML template strings
pkg/templates/render.go      # Template selection + execution
pkg/telegram/client.go       # HTTP client: sendMessage with inline keyboard
testdata/                    # Sample GitHub context JSON payloads
```

## Development Commands

```bash
make test          # Run all tests
make test-v        # Verbose tests
make lint          # go vet + gofmt check
make build         # Build binary
make docker        # Build Docker image
go test ./... -race  # Tests with race detector
```

## Local Testing

```bash
cp .envrc.example .envrc    # Fill in bot token + chat ID
source .envrc               # Or use direnv
go run .                    # Sends test notification
```

## Key Design Decisions

- **`html/template` not `text/template`**: Auto-escapes `{{.Field}}` values while leaving template HTML intact. Required for safe Telegram messages.
- **`env` mapping not `args` in action.yml**: Reliable for multi-line content (custom templates, large JSON payloads).
- **Inputs via `INPUT_*` env vars**: `INPUT_BOT_TOKEN`, `INPUT_CHAT_ID`, `INPUT_TOPIC_ID`, `INPUT_CUSTOM_TEMPLATE`, `INPUT_EVENT_PAYLOAD`.
- **Bot token is masked**: `::add-mask::` is called before any processing.

## Code Conventions

- Zero external dependencies. stdlib only.
- Use conventional commits (`feat:`, `fix:`, `ci:`, `docs:`).
- Use `html/template` for any Telegram message rendering.
- Run `make lint` before committing. Pre-commit hooks handle `go fmt`, `go vet`, `go build`, and TruffleHog.
- Tests use table-driven patterns where applicable.
- Never commit `.envrc` or any file containing tokens/secrets.
- Pin GitHub Actions by SHA, Docker images by digest.

## Supported Events

| Event | Actions |
|-------|---------|
| `pull_request` | `opened`, `closed` (detects merged), `reopened`, `synchronize`, `ready_for_review`, `converted_to_draft` |
| `pull_request_review` | `submitted` → resolves to `approved`, `changes_requested`, `commented` |
| `pull_request_review_comment` | `created` |

## Adding a New Event Type

1. Add struct + parse function in `pkg/events/events.go`
2. Add template constant in `pkg/templates/defaults.go`
3. Register in `defaultTemplates` map with `event_name:action` key
4. Add testdata JSON fixture in `testdata/`
5. Add test cases in `pkg/events/events_test.go` and `pkg/templates/render_test.go`

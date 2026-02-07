# Contributing to Telegram PR Notify

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) (for building the container image)
- [direnv](https://direnv.net/) (optional, for local testing)

## Development Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/andoniaf/telegram-pr-notify.git
   cd telegram-pr-notify
   ```

2. Copy the environment example and fill in your test credentials:

   ```bash
   cp .envrc.example .envrc
   # Edit .envrc with your Telegram bot token and chat ID
   direnv allow  # if using direnv
   ```

3. Run the tests:

   ```bash
   make test
   ```

## Project Structure

```
.
├── main.go                  # Entry point, reads env vars and orchestrates
├── pkg/
│   ├── events/              # GitHub event parsing and TemplateData model
│   ├── templates/           # Template rendering and default templates
│   └── telegram/            # Telegram Bot API client
├── testdata/                # JSON fixtures for event parsing tests
├── action.yml               # GitHub Action definition
└── Dockerfile               # Multi-stage build for the action container
```

## Running Tests

```bash
make test        # Run all tests
make test-v      # Run tests with verbose output
make lint        # Run go vet + format check
make build       # Build the binary
```

## Local Testing

With `direnv` and a configured `.envrc`:

```bash
go run .
```

This sends a test notification using the payload in `.envrc` (defaults to `testdata/pull_request_opened.json`).

## Docker Build

```bash
make docker
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Keep external dependencies at zero (stdlib only).
- Use `html/template` (not `text/template`) for Telegram HTML messages.

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.

## Pull Requests

- Include tests for new event types or template changes.
- Update README.md if adding user-facing features.
- Ensure `make test` and `make lint` pass.

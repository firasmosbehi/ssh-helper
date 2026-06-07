# Agent Guidance

## Project

`ssh-helper` is a Go project that provides CLI, TUI and GUI interfaces for SSH and rsync workflows. It also implements the Model Context Protocol (MCP) as both a server and a client.

## Tech stack

- Go 1.23+
- CLI: cobra + viper
- TUI: bubbletea + lipgloss + bubbles
- GUI: fyne.io/fyne/v2
- SSH: golang.org/x/crypto/ssh
- SSH config: kevinburke/ssh_config
- MCP: github.com/mark3labs/mcp-go
- Store: JSON / SQLite
- Secrets: zalando/go-keyring

## Build

```bash
go build -o ssh-helper ./cmd/ssh-helper
```

## Tests

```bash
go test ./...
```

## Conventions

- Keep `internal/` packages focused and small.
- Expose reusable logic under `pkg/api/` if needed.
- Use `slog` for all logging.
- Write table-driven tests for exported functions.
- Tag issues and PRs with `area/cli`, `area/tui`, `area/gui`, or `area/mcp`.

# ssh-helper

[![CI](https://github.com/firasmosbehi/ssh-helper/actions/workflows/ci.yml/badge.svg)](https://github.com/firasmosbehi/ssh-helper/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/firasmosbehi/ssh-helper)](https://github.com/firasmosbehi/ssh-helper/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

CLI + TUI + GUI tool for managing SSH connections, keys, configs and running rsync transfers.

Built in Go. Ships as a single binary.

## Features

- **SSH Host Management** — parse, list, add, edit and connect to hosts from `~/.ssh/config`.
- **SSH Key Management** — list, generate, remove keys and copy public keys to remote hosts.
- **rsync Jobs** — run one-off transfers or save/sync recurring jobs with live progress.
- **Terminal UI** — interactive Bubble Tea app for browsing hosts, jobs and keys.
- **Desktop GUI** — native Fyne interface for host, job, key and MCP management (CGO builds).
- **MCP Server & Client** — expose SSH/rsync tools to MCP consumers and call external MCP servers.
- **Cross-Platform** — Linux, macOS and Windows binaries.

## Install

### Homebrew (macOS / Linux)

```bash
brew tap firasmosbehi/tap
brew install ssh-helper
```

### Scoop (Windows)

```powershell
scoop bucket add ssh-helper https://github.com/firasmosbehi/scoop-bucket
scoop install ssh-helper
```

### Go

```bash
go install github.com/firasmosbehi/ssh-helper/cmd/ssh-helper@latest
```

### Download

Grab a pre-built binary from the [releases page](https://github.com/firasmosbehi/ssh-helper/releases).

## Quick start

```bash
# Initialize configuration
ssh-helper config init

# List hosts from ~/.ssh/config
ssh-helper host list

# Connect to a host
ssh-helper host connect myserver

# Run a one-off rsync transfer
ssh-helper sync run ./local/dir/ user@host:/remote/dir/

# Launch the TUI
ssh-helper tui

# Launch the GUI (CGO-enabled builds only)
ssh-helper gui

# Start the MCP server
ssh-helper mcp serve
```

## Documentation

- [Full documentation](https://firasmosbehi.github.io/ssh-helper/)
- [Installation guide](docs/installation.md)
- [CLI reference](docs/cli.md)
- [TUI cheatsheet](docs/tui.md)
- [MCP schema](docs/mcp.md)
- [Security model](docs/security.md)

## Project layout

```
cmd/ssh-helper/      # Application entrypoint
internal/            # Private application code
  config/            # Configuration management
  core/              # Domain models
  ssh/               # SSH client and config parsing
  rsync/             # rsync runner and progress parser
  store/             # Persistence layer
  tui/               # Terminal user interface
  gui/               # Desktop GUI (Fyne)
  mcp/               # MCP server and client
  platform/          # OS-specific helpers
pkg/api/             # Public library API
docs/                # Documentation
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT

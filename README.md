# ssh-helper

CLI + TUI + GUI tool for managing SSH connections, keys, configs and running rsync transfers.

Built in Go. Ships as a single binary.

## Interfaces

- **CLI**: fast commands for scripts and terminal power users (`ssh-helper connect`, `ssh-helper sync`, ...).
- **TUI**: interactive terminal UI powered by Bubble Tea.
- **GUI**: native desktop UI powered by Fyne.
- **MCP**: acts as an MCP server (expose SSH/rsync tools to AI agents) and MCP client (call external MCP servers).

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
scripts/             # Development scripts
```

## Quick start

### Build

```bash
go build -o ssh-helper ./cmd/ssh-helper
```

### Run

```bash
# Show help
./ssh-helper --help

# Initialize config
./ssh-helper config init

# List hosts from ~/.ssh/config
./ssh-helper host list

# Start TUI
./ssh-helper tui

# Start GUI
./ssh-helper gui

# Run MCP server
./ssh-helper mcp serve
```

## Documentation

See [docs/](docs/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT

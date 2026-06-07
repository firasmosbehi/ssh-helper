# ssh-helper

CLI + TUI + GUI tool for managing SSH connections, keys, configs and running rsync transfers.

Built in Go. Ships as a single binary.

## Interfaces

- **CLI**: fast commands for scripts and terminal power users (`ssh-helper connect`, `ssh-helper sync`, ...).
- **TUI**: interactive terminal UI powered by Bubble Tea.
- **GUI**: native desktop UI powered by Fyne.
- **MCP**: acts as an MCP server (expose SSH/rsync tools to AI agents) and MCP client (call external MCP servers).

## Quick start

```bash
# Install
brew install firasmosbehi/tap/ssh-helper

# Initialize config
ssh-helper config init

# List hosts from ~/.ssh/config
ssh-helper host list

# Start TUI
ssh-helper tui

# Start GUI
ssh-helper gui

# Run MCP server
ssh-helper mcp serve
```

## Documentation

See [docs/](docs/).

## License

MIT

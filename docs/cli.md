# CLI Reference

ssh-helper is a unified CLI for SSH and rsync management.

## Global flags

| Flag | Description |
|------|-------------|
| `-v`, `--verbose` | Increase verbosity (use `-vv` for debug) |
| `--json` | Output logs as JSON |
| `--version` | Print version |

## Commands

### `config`
Manage application configuration.

- `config init` — Initialize default config files.
- `config get <key>` — Get a config value.
- `config set <key> <value>` — Set a config value.
- `config edit` — Open config in `$EDITOR`.

### `host`
Manage SSH hosts.

- `host list` — List all hosts from `~/.ssh/config`.
- `host show <name>` — Show host details.
- `host add <name> <hostname>` — Add a new host.
- `host remove <name>` — Remove a host.
- `host edit` — Open `~/.ssh/config` in `$EDITOR`.
- `host connect <name>` — Connect to a host interactively.

### `key`
Manage SSH keys.

- `key list` — List keys with fingerprints.
- `key generate <name>` — Generate a new key (ed25519 or rsa).
- `key remove <name>` — Remove a key pair.
- `key copy-id <host> <key>` — Copy public key to a remote host.
- `key agent-add <name>` — Add a key to the SSH agent.

### `sync`
Run and manage rsync jobs.

- `sync run <src> <dst>` — Run a one-off rsync transfer.
- `sync job create <name> <src> <dst>` — Create a saved job.
- `sync job list` — List saved jobs.
- `sync job delete <id>` — Delete a saved job.
- `sync job run <id>` — Run a saved job.

### `tui`
Launch the interactive terminal UI.

### `gui`
Launch the native desktop GUI.

### `mcp`
MCP server and client commands.

- `mcp serve` — Start the built-in MCP server over stdio.
- `mcp servers` — List configured external MCP servers.
- `mcp call <server> <tool> --args '{}'` — Call a tool on an external MCP server.

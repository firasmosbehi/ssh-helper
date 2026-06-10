# MCP Schema

ssh-helper supports the Model Context Protocol (MCP) both as a server and a client.

## Built-in MCP Server

Start the server:

```bash
ssh-helper mcp serve
```

### Tools exposed

| Tool | Description |
|------|-------------|
| `list_hosts` | List SSH hosts from `~/.ssh/config` |
| `show_host` | Show details for a specific host |
| `run_command` | Run a command on a remote host via SSH |
| `transfer_files` | Transfer files via rsync |
| `list_keys` | List SSH keys |
| `generate_key` | Generate a new SSH key |
| `list_sync_jobs` | List saved rsync jobs |
| `run_sync_job` | Run a saved rsync job |

## External MCP Clients

Configure external MCP servers in `~/.config/ssh-helper/config.yaml`:

```yaml
mcp_servers:
  my-server:
    command: npx
    args: ["-y", "@modelcontextprotocol/server-filesystem"]
    env:
      SOME_VAR: value
```

### CLI commands

- `ssh-helper mcp servers` — List configured servers.
- `ssh-helper mcp call <server> <tool> --args '{}'` — Invoke a tool.

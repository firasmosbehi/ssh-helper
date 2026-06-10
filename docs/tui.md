# TUI Cheatsheet

Launch the TUI with:

```bash
ssh-helper tui
```

## Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select / Drill down |
| `Esc` | Go back |
| `q` / `Ctrl+C` | Quit |

## Tabs

### Hosts
- Browse SSH hosts from `~/.ssh/config`.
- Type to filter by name or hostname.
- Press `Enter` to see connection details.

### Sync Jobs
- Browse saved rsync jobs.
- Press `Enter` to run a job.

### Keys
- Browse SSH keys with fingerprints.

### MCP
- Browse configured external MCP servers.
- `Enter` on a server to list its tools.
- `Enter` on a tool to input JSON arguments and invoke it.

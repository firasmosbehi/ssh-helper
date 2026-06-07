#!/bin/bash
set -e
REPO="firasmosbehi/ssh-helper"

create_issue() {
  local title="$1"
  local body="$2"
  local milestone="$3"
  local labels="$4"
  echo "Creating issue: $title"
  gh issue create --repo "$REPO" --title "$title" --body "$body" --milestone "$milestone" --label "$labels"
}

# M0
create_issue "Initialize Go module, folder layout, and README" "Bootstrap the repository with the agreed internal/ folder structure, initialize go.mod, and ensure README describes the project." "M0 — Repo Bootstrap" "kind/feature,area/cli"
create_issue "Add AGENTS.md and contributor docs" "Write AGENTS.md with conventions, build/test instructions, and label usage. Add a basic CONTRIBUTING.md." "M0 — Repo Bootstrap" "kind/feature"
create_issue "Set up GitHub Actions: lint, test, build matrix" "Add a workflow that runs golangci-lint, go test ./..., and builds on macOS, Linux, Windows." "M0 — Repo Bootstrap" "kind/feature,area/cli"
create_issue "Add issue/PR templates and labels" "Create bug and feature request templates, plus PR template. Ensure labels kind/feature, kind/bug, area/cli, area/tui, area/gui, area/mcp exist." "M0 — Repo Bootstrap" "kind/feature"

# M1
create_issue "Implement config init and global config loading with Viper" "Add ssh-helper config init and load user settings from ~/.config/ssh-helper/config.yaml." "M1 — CLI Core + Config" "kind/feature,area/cli"
create_issue "Implement config get/set/edit" "CLI subcommands to read, update, and open config in the default editor." "M1 — CLI Core + Config" "kind/feature,area/cli"
create_issue "Add persistent data directory and JSON store abstraction" "Define the data dir, create a Store interface, and implement JSON-backed persistence for hosts/keys/jobs." "M1 — CLI Core + Config" "kind/feature,area/cli"
create_issue "Add structured logging and --verbose/--json flags" "Replace fmt prints with slog. Support text and JSON output, plus verbosity levels." "M1 — CLI Core + Config" "kind/feature,area/cli"

# M2
create_issue "Parse ~/.ssh/config and implement host list/show" "Use kevinburke/ssh_config to read the system SSH config and expose host list/show commands." "M2 — SSH Management" "kind/feature,area/cli"
create_issue "Implement host add/remove/edit with backup" "Safely mutate ~/.ssh/config with automatic backups and validation." "M2 — SSH Management" "kind/feature,area/cli"
create_issue "Implement connect <host> command" "Open an SSH session via golang.org/x/crypto/ssh with password, key, and agent auth support." "M2 — SSH Management" "kind/feature,area/cli"
create_issue "Implement key list/generate/remove/copy-id" "Manage SSH keys: enumerate, generate ed25519/RSA, delete, and copy public keys to remote hosts." "M2 — SSH Management" "kind/feature,area/cli"
create_issue "Add SSH agent integration and keyring-backed passphrase cache" "Use the OS SSH agent and store passphrases securely with go-keyring." "M2 — SSH Management" "kind/feature,area/cli"

# M3
create_issue "Design rsync job model and CLI command builder" "Define a Job struct with src/dest/flags/excludes/dry-run and a builder that produces safe rsync command args." "M3 — rsync Runner & Jobs" "kind/feature,area/cli"
create_issue "Implement sync run with live progress and dry-run support" "Run rsync, stream stdout/stderr, and display progress in CLI. Support --dry-run." "M3 — rsync Runner & Jobs" "kind/feature,area/cli"
create_issue "Implement sync job create/list/delete/run/edit" "Persist sync jobs and expose full CRUD commands." "M3 — rsync Runner & Jobs" "kind/feature,area/cli"
create_issue "Parse rsync --info=progress2 output for real-time progress" "Build a parser to extract file counts, bytes, and speed from rsync progress output." "M3 — rsync Runner & Jobs" "kind/feature,area/cli"
create_issue "Persist job history and logs to SQLite/JSON store" "Store executed job history and output logs, searchable by host or date." "M3 — rsync Runner & Jobs" "kind/feature,area/cli"

# M4
create_issue "Bootstrap Bubble Tea app with global keybindings and theme" "Create the main TUI loop, navigation, and a consistent lipgloss theme." "M4 — TUI" "kind/feature,area/tui"
create_issue "TUI host browser with fuzzy search and connect action" "Interactive host list with filter and Enter to connect." "M4 — TUI" "kind/feature,area/tui"
create_issue "TUI sync job browser with run/preview/logs" "Browse sync jobs, run them, preview command, and view logs." "M4 — TUI" "kind/feature,area/tui"
create_issue "TUI key manager" "List, generate, and delete SSH keys from the TUI." "M4 — TUI" "kind/feature,area/tui"
create_issue "TUI real-time progress view for running syncs" "Show live progress bars and transfer stats while rsync runs." "M4 — TUI" "kind/feature,area/tui"

# M5
create_issue "Integrate MCP Go SDK and serve stdio transport" "Add mcp-go dependency and implement ssh-helper mcp serve using stdio transport." "M5 — MCP Server" "kind/feature,area/mcp"
create_issue "Implement list_hosts, show_host, connect_host MCP tools" "Expose host read/connect actions as MCP tools with clear schemas." "M5 — MCP Server" "kind/feature,area/mcp"
create_issue "Implement run_command and transfer_files MCP tools with approval hooks" "Expose remote command execution and rsync transfers; require user approval before destructive ops." "M5 — MCP Server" "kind/feature,area/mcp"
create_issue "Implement list_keys and generate_key MCP tools" "Expose SSH key management as MCP tools." "M5 — MCP Server" "kind/feature,area/mcp"
create_issue "Implement list_sync_jobs and run_sync_job MCP tools" "Expose sync job read/run as MCP tools." "M5 — MCP Server" "kind/feature,area/mcp"
create_issue "Add approval prompts for destructive MCP operations" "Interactive confirmation when MCP tools run commands, transfer files, or generate/rotate keys." "M5 — MCP Server" "kind/feature,area/mcp"

# M6
create_issue "Config schema for external MCP servers" "Allow users to register external MCP servers under mcp_servers in config." "M6 — MCP Client + External Integrations" "kind/feature,area/mcp"
create_issue "Implement mcp list and mcp call CLI commands" "List registered MCP servers and call a tool by name with JSON arguments from CLI." "M6 — MCP Client + External Integrations" "kind/feature,area/mcp,area/cli"
create_issue "TUI MCP server browser and tool invoker" "Browse external MCP servers, inspect tools, and invoke them from the TUI." "M6 — MCP Client + External Integrations" "kind/feature,area/mcp,area/tui"
create_issue "Example integration: fetch secrets from MCP server for SSH passphrase" "Demonstrate consuming an MCP secrets server to unlock SSH keys." "M6 — MCP Client + External Integrations" "kind/feature,area/mcp"

# M7
create_issue "Bootstrap Fyne app and main navigation" "Create the GUI window, menus, and router between host/sync/key/mcp views." "M7 — GUI (Fyne)" "kind/feature,area/gui"
create_issue "GUI host manager (list, add, edit, connect)" "Fyne view for managing SSH hosts." "M7 — GUI (Fyne)" "kind/feature,area/gui"
create_issue "GUI sync job manager with run and progress views" "Fyne view for sync jobs including progress bars." "M7 — GUI (Fyne)" "kind/feature,area/gui"
create_issue "GUI key manager" "Fyne view for SSH key list/generate/remove." "M7 — GUI (Fyne)" "kind/feature,area/gui"
create_issue "GUI MCP server/tool browser" "Fyne view to browse internal and external MCP servers/tools." "M7 — GUI (Fyne)" "kind/feature,area/gui"

# M8
create_issue "Configure goreleaser with Homebrew tap and scoop" "Set up .goreleaser.yaml to publish binaries, checksums, and package managers." "M8 — CI/CD, Packaging, Docs" "kind/feature"
create_issue "Build and sign macOS/Linux/Windows binaries in CI" "Add release workflow with code signing and notarization where possible." "M8 — CI/CD, Packaging, Docs" "kind/feature"
create_issue "Write user guide in docs/: installation, CLI reference, TUI cheatsheet" "Complete Markdown docs for end users." "M8 — CI/CD, Packaging, Docs" "kind/feature"
create_issue "Add automated CLI docs generation" "Generate CLI markdown docs from cobra commands in CI." "M8 — CI/CD, Packaging, Docs" "kind/feature,area/cli"
create_issue "Add code coverage reporting and branch protection rules" "Upload coverage, require PR reviews and passing CI before merge." "M8 — CI/CD, Packaging, Docs" "kind/feature"

# M9
create_issue "Security audit: input sanitization, secret redaction, safe rsync/ssh invocation" "Review all user inputs, avoid injection, redact secrets in logs, validate rsync args." "M9 — Hardening & Production Ready" "kind/feature"
create_issue "Add integration tests with Docker-based SSH server fixture" "Spin up an OpenSSH container in tests to validate connect, copy-id, and rsync flows." "M9 — Hardening & Production Ready" "kind/feature"
create_issue "Performance: large host lists, large sync history, progress rendering at 60 FPS" "Profile TUI and GUI with big datasets; optimize rendering and queries." "M9 — Hardening & Production Ready" "kind/feature,area/tui,area/gui"
create_issue "Accessibility and localization groundwork" "Add keyboard shortcuts, basic a11y labels, and a strings file for future i18n." "M9 — Hardening & Production Ready" "kind/feature,area/gui,area/tui"
create_issue "v1.0.0 release checklist and announcement" "Final release notes, tag v1.0.0, publish packages, and announce." "M9 — Hardening & Production Ready" "kind/feature"

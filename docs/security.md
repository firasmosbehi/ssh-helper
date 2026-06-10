# Security Model

## Data Storage

- Configuration is stored in `~/.config/ssh-helper/`.
- Job history and keys metadata are stored in `~/.local/share/ssh-helper/`.
- SSH config edits create backups before modification.
- Generated private keys are written with `0o600` permissions.

## Secrets

- SSH key passphrases can be cached in the OS keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service).
- No plaintext password storage.

## Process Safety

- rsync and SSH commands are executed via `exec.Command` with explicit argument lists (no shell injection).
- MCP server tools validate inputs against declared schemas.

## Networking

- SSH connections use the standard OpenSSH client library (`golang.org/x/crypto/ssh`).
- Host key verification follows the user's `~/.ssh/known_hosts`.

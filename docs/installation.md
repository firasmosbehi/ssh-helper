# Installation

## Homebrew (macOS/Linux)

```bash
brew tap firasmosbehi/tap
brew install ssh-helper
```

## Scoop (Windows)

```powershell
scoop bucket add ssh-helper https://github.com/firasmosbehi/scoop-bucket
scoop install ssh-helper
```

## Go install

```bash
go install github.com/firasmosbehi/ssh-helper/cmd/ssh-helper@latest
```

## From source

```bash
git clone https://github.com/firasmosbehi/ssh-helper.git
cd ssh-helper
go build ./cmd/ssh-helper
```

## Requirements

- Go 1.25+ (for building from source)
- OpenSSH client (`ssh`, `ssh-keygen`, `ssh-agent`)
- rsync (for sync jobs)

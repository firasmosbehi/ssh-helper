# Contributing to ssh-helper

Thank you for your interest in contributing!

## Getting started

1. Fork the repository.
2. Clone your fork and create a feature branch.
3. Make your changes, add tests, and ensure the build passes.
4. Open a pull request using the provided template.

## Development

```bash
go build -o ssh-helper ./cmd/ssh-helper
go test ./...
golangci-lint run
```

## Conventions

- Keep packages small and focused.
- Use `slog` for logging.
- Write table-driven tests.
- Tag issues and PRs with the appropriate `area/*` label.

## Code of conduct

Be respectful and constructive in all interactions.

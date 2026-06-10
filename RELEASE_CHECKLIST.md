# Release Checklist

## Pre-release

- [ ] All tests pass: `go test ./...`
- [ ] Integration tests pass: `go test -tags=integration ./integration/...`
- [ ] Lint passes: `gofmt -l .` and `go vet ./...`
- [ ] Security audit items resolved (see `docs/security.md`)
- [ ] CHANGELOG.md updated with version notes
- [ ] Version bumped in code (tag will drive `ldflags`)
- [ ] Documentation regenerated: `go run ./cmd/ssh-helper gendocs docs/cli`

## Tag & Build

- [ ] Create Git tag: `git tag -s v1.0.0 -m "Release v1.0.0"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] Verify CI passes on tag
- [ ] Verify goreleaser completes and artifacts are uploaded
- [ ] Verify Homebrew formula updated
- [ ] Verify Scoop manifest updated

## Post-release

- [ ] GitHub Release notes published
- [ ] Binary signatures verified (if signing enabled)
- [ ] Announcement posted (blog, social, etc.)
- [ ] Close release milestone

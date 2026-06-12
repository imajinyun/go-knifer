## Summary

- 

## Change type

- [ ] Bug fix
- [ ] New public API
- [ ] Internal refactor
- [ ] Documentation
- [ ] CI / governance

## Checklist

- [ ] I kept `v*` packages as thin facades over `internal/*`.
- [ ] I avoided new dependencies in public facades, or updated the architecture allowlist with justification.
- [ ] I added or updated tests for changed behavior.
- [ ] I ran `go test -race -shuffle=on -coverprofile=coverage.out ./...`.
- [ ] I ran `bash bin/check_coverage.sh coverage.out`.
- [ ] I ran `bash bin/check_arch.sh`.
- [ ] I ran `golangci-lint run ./...`.
- [ ] I updated `CHANGELOG.md` for user-visible changes.

## Security review

- [ ] This change does not touch security-sensitive code.
- [ ] This change touches security-sensitive code and includes regression tests.
- [ ] Any security linter suppression is narrow and documented.

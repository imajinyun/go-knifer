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
- [ ] I ran `make quick-check` or documented why it is not applicable.
- [ ] I ran `go test -race -shuffle=on -coverprofile=coverage.out ./...` when the change is non-trivial.
- [ ] I ran `bash bin/check_coverage.sh coverage.out` when a fresh coverage profile was generated.
- [ ] I ran `bash bin/check_arch.sh`.
- [ ] I ran `golangci-lint run ./...`.
- [ ] I updated `CHANGELOG.md` for user-visible changes.

## Validation

- Commands run:
  -
- Commands intentionally skipped and reason:
  -

## API and coverage impact

- Public API changed: yes / no
- `docs/api/exports.txt` updated: yes / no / not applicable
- Coverage impact:
  -

## Reviewer focus

-

## Intentionally excluded files

-

## Security review

- [ ] This change does not touch security-sensitive code.
- [ ] This change touches security-sensitive code and includes regression tests.
- [ ] Any security linter suppression is narrow and documented.

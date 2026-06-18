# go-knifer — AI Agent Guide

> Go utility library (48 public `v*` facade packages + `internal/*` implementations).
> **Do not modify production code unless the task explicitly requires it.**

## Project layout

```
./
├── doc.go              # root package: domain-grouped navigation, no business APIs
├── errors.go           # root error contract: ErrCode, Error, CodeCarrier, CodeOf
├── v<domain>/          # public facade — import these (thin: forwards to internal/*)
│   └── doc.go          # required: package doc with full domain name on first line
├── internal/<domain>/  # implementation — not importable by external modules
├── bin/                # validation scripts (api_snapshot, check_arch, check_coverage)
├── docs/
│   ├── api/exports.txt # exported API snapshot (19k+ lines, CI-enforced)
│   └── doc/            # 48 per-package quickstart docs (01-vbean.md .. 48-vzip.md)
└── .github/workflows/  # CI (go.yml) + release (release.yml) automation
```

- **Public facades never import each other.** Shared logic goes into `internal/*`.
- **Internal packages never import public facades.** Direction: `v* → internal/*`.
- **Production `panic` is exceptional.** Only in `MustXxx`/`PanicXxx` APIs or documented cases.
- **Facades are thin.** No business logic, loops, `panic`, or type assertions in `v*`.

## Package catalog (48 v* packages)

| Package | Domain | internal |
|---------|--------|----------|
| vstr | string/text | internal/str |
| vslice | slices | internal/slice |
| vmap | maps | internal/maps |
| vcache | caches (FIFO/LRU/LFU/TTL) | internal/cache |
| vcrypto | cryptography/digests | internal/crypto |
| vhttp | HTTP (stdlib) | internal/httpx |
| vresty | HTTP (Resty) | internal/httpx/resty |
| vjson | JSON | internal/json |
| vconf | configuration | internal/conf |
| vid | UUID/Snowflake/NanoId | internal/id |
| vjwt | JWT sign/verify | internal/jwt |
| vrand | randomness | internal/rand |
| vlog | logging | internal/log |
| verr | error handling/panic recovery | internal/errx |
| vfile | file/IO | internal/file |
| vurl | URL/URI | internal/url |
| vmask | data masking | internal/mask |
| vcodec | Base64/Hex | internal/codec |
| vhash | non-crypto hashes | internal/hash |
| vzip | archive/compression | internal/zip |
| vset | sets | internal/sets |
| vobj | object helpers | internal/obj |
| vconv | type conversion | internal/conv |
| vdate | date/time | internal/date |
| vref | reflection | internal/ref |
| vregex | regular expressions | internal/regex |
| vtpl | templates | internal/template |
| vbool | booleans | internal/boolean |
| vnum | numeric helpers | internal/num |
| vbean | struct/map mapping | internal/bean |
| vblf | bloom filter | internal/bloomfilter |
| vcsv | CSV | internal/csvx |
| vimg | images/captchas | internal/imgx |
| vxml | XML | internal/xml |
| vmail | email/SMTP | internal/mail |
| vskt | sockets | internal/socket |
| vnet | IP/port/network | internal/net |
| vdb | database/SQL | internal/db |
| vcron | cron scheduling | internal/cron |
| vjob | job orchestration | internal/job |
| vsem | semaphores | internal/semaphore |
| vdfa | DFA word-tree matching | internal/dfa |
| vpass | password strength | internal/pass |
| vident | identity numbers | internal/identity |
| vform | form/input validation | internal/validator |
| vver | version comparison | internal/version |
| vsys | system information | internal/system |
| vpoi | office documents (Excel) | internal/poi |
| vyaml | YAML | internal/yaml |

## Validation commands

| Command | Scope |
|---------|-------|
| `make quick-check` | Fast local: mod-verify → vet → arch → test → api-check → diff-whitespace |
| `make security-check` | Lint + govulncheck |
| `make full-check COVERAGE_FILE=/tmp/coverage.out` | Full pre-push: quick-check + race coverage + coverage gate + lint + vuln |
| `make ci-test` | CI test-job gate |
| `make check` | Alias for `full-check` |
| `UPDATE_API=1 make api-check` | Refresh API snapshot after intentional public API changes |
| `make bench-core` | Core benchmark baselines (internal/slice, maps, str, num) |
| `make bench-facade` | Facade benchmark baselines (vslice, vmap, vstr, vnum) |
| `make govulncheck` | Vulnerability scan |

## Change workflow (end-to-end)

1. **Inspect** existing code and docs before modifying.
2. **Apply** only the requested logical change; no unrelated files.
3. **Format** touched Go files with `gofmt -w`.
4. **Validate** focused tests first, then broaden:
   - `go test -v -gcflags="all=-l -N" ./<changed-package>`
   - `go test -v -gcflags="all=-l -N" ./...`
   - `go vet ./...`
   - `bash bin/check_arch.sh` (after production changes)
   - `bash bin/check_api_compat.sh` (after public API changes)
   - `golangci-lint run ./...`
   - `go test -race -shuffle=on -coverprofile=/tmp/cov.out ./... && bash bin/check_coverage.sh /tmp/cov.out`
   - `git diff --check`
   - Prefer `make quick-check` / `make full-check` aggregate targets.
5. **Commit** with conventional commit message (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`).
6. **Push** to remote.

## API design rules

- Prefer `(result, error)` over `panic` for recoverable failures.
- IO/network functions take `context.Context` as first parameter.
- No synonym aliases for the same function.
- Renaming a public symbol is a breaking change (SemVer).
- Security-sensitive code: `vhttp`, `vresty`, `vurl`, `vconf`, `vzip`, `vfile`, `vcrypto`, `vjwt`, `vrand`, `vid`, `vdb`. See `SECURITY.md`.

## Governance constraints

- **Coverage**: Keep total coverage above the threshold in `bin/check_coverage.sh`.
- **Architecture**: 8 rules enforced by `bin/check_arch.sh` — doc.go existence, no v*-to-v* imports, per-file internal/ imports, no internal→v* imports, package comments, panic policy, facade boundary policy, dependency allowlist.
- **API snapshot**: `docs/api/exports.txt` is CI-enforced. Run `UPDATE_API=1 make api-check` after intentional public API changes.
- **Panic**: Production code must not introduce new `panic()` calls unless in a `MustXxx`/`PanicXxx` function.
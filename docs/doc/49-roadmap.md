# go-knifer Next-Phase Capability Roadmap

This roadmap tracks the next development phase after the tool catalog quality sprints. It prioritizes scenario mindshare, deeper high-value modules, documentation and benchmark trust, and explicit ecosystem adapter lanes.

## Baseline

| Metric | Value |
| --- | ---: |
| Public facade packages | 48 |
| Public functions | 2528 |
| Functions with executable examples | 241 |
| Context-aware functions | 20 |
| Empty function synopses | 0 |

## Strategic themes

1. **Scenario mindshare** — help users choose the right package for a concrete task before adding more APIs.
2. **Deep modules** — improve high-value existing packages where users already expect depth.
3. **Documentation and benchmark trust** — make claims testable through examples, benchmarks, release signals, and clear boundaries.
4. **Ecosystem adapters** — schedule AI, FTP, SSH/SFTP, pinyin, tokenization, multi-template engines, and CLI utilities as explicit development lanes.

## Capability matrix

| Area | Current package | Gap | Priority | First deliverable |
| --- | --- | --- | --- | --- |
| Collections | `vslice`, `vmap` | lo-style error variants, window/sliding helpers, partitioning, zip/unzip, benchmark comparisons | P1 | Add collection gap tests, examples, and benchmark baselines |
| Bean mapping | `vbean`, `vmap` | copy/decode/merge semantics, deep copy, decode hooks, unused-key metadata, default merge | P1 | Define copy/decode/merge contract docs and focused tests |
| Conversion | `vconv` | consistent zero-value/error/default API naming and weak-input docs | P1 | Add conversion API contract docs and examples |
| Validation | `vform` | struct tag validation or explicit recommendation to use `go-playground/validator` | P1 | Record validation direction and update docs/tests accordingly |
| Benchmarks | existing benchmark files | public benchmark narrative and competitor-neutral baselines | P1 | Add benchmark trust section and stable benchmark commands |
| Examples | all large facades | low function-level example ratio in large packages | P1 | Raise reader-facing examples in `vhttp`, `vnet`, `vnum`, `vresty`, and `vzip` |
| CLI utilities | none | command args, environment helpers, process execution, terminal IO helpers | P1 | Design or implement a `vcli` MVP with context-aware execution |
| AI adapters | none | provider abstraction for chat, embeddings, streaming, and tool calls | P1 | Design `vai` with fake-provider tests and redacted logging examples |
| Multi-template engines | `vtpl` | engine abstraction beyond `html/template` | P2 | Define adapter interface and preserve the standard template baseline |
| FTP | none | client helpers, upload/download/list, context support, provider injection | P2 | Design `vftp` with bounded transfer and fake-provider tests |
| SSH/SFTP | none | SSH command execution and SFTP file transfer helpers | P2 | Design `vssh` with explicit host-key verification and no secret leakage |
| Pinyin | none | Chinese transliteration helpers | P2 | Design `vpinyin` with deterministic dictionary/provider injection |
| Tokenization | none | Chinese text segmentation adapters | P2 | Design `vtokenize` with deterministic examples and no network dependency |
| Database | `vdb` | context-first APIs, dialect depth, batch/upsert/scan helpers | P2 | Create a `vdb` deepening backlog and focused tests |
| Crypto | `vcrypto`, `vjwt`, `vrand` | TOTP/HOTP, password hashing, JWK/JWKS, optional national algorithms | P2 | Create a safe crypto extension plan that keeps insecure algorithms out of defaults |
| Office | `vpoi` | streaming Excel, styles, formulas, images, Word/OFD scope decision | P2 | Decide the `vpoi` scope before adding broad Office dependencies |
| Image | `vimg` | GIF, color, draw, watermark, deeper image processing | P3 | Add image backlog with deterministic fixtures and benchmarks |

## Sprint order

| Sprint | Name | Outcome |
| --- | --- | --- |
| 9 | Capability Matrix and Trust Roadmap | Publish this roadmap and link it from the documentation hub. |
| 10 | Collection Mindshare | Deepen `vslice` and `vmap` around lo-style scenarios, examples, and benchmarks. |
| 11 | Bean Copy/Decode/Merge Semantics | Split `vbean` behavior into documented copy, decode, and merge lanes. |
| 12 | Validation Direction | Either add a lightweight `vvalidate` package or explicitly recommend `go-playground/validator` for struct tag validation. |
| 13 | Benchmark and Example Trust | Add benchmark narrative and raise examples in the largest public facades. |
| 14 | Developer Experience Adapters | Plan or implement `vcli`, `venv`, `vdump`, `vtest`, `vretry`, `vctx`, and `vwatch` lanes. |
| 15 | Ecosystem Adapter Lane 1 | Plan or implement AI and multi-template adapter foundations. |
| 16 | Ecosystem Adapter Lane 2 | Plan or implement FTP, SSH/SFTP, pinyin, and tokenization foundations. |
| 17 | Deep Business Modules | Deepen `vdb`, `vcrypto`, `vpoi`, and `vimg` after trust foundations are in place. |

## Scenario guidance

| Scenario | Use now | Planned lane |
| --- | --- | --- |
| Transform slices/maps with type-safe helpers | `vslice`, `vmap` | lo-style error variants, advanced grouping/window helpers, and benchmark comparisons |
| Copy, decode, or merge struct/map data | `vbean`, `vmap` | copy/decode/merge semantic split with deep-copy and metadata options |
| Validate common strings and identity formats | `vform`, `vident` | validation decision: lightweight `vvalidate` or documented `validator` recommendation |
| Build safe HTTP clients or open untrusted URLs | `vhttp`, `vresty`, `vurl` | more examples and benchmarked helper paths |
| Build CLI or terminal utilities | Standard library today | planned `vcli`, `venv`, `vdump`, and `vtest` lanes |
| Call AI model providers | Provider SDKs today | planned `vai` provider abstraction lane |
| Transfer files over FTP or SSH/SFTP | Dedicated clients today | planned `vftp` and `vssh` lanes |
| Convert Chinese text to pinyin or tokenize text | Dedicated NLP libraries today | planned `vpinyin` and `vtokenize` lanes |
| Render templates beyond `html/template` | `vtpl` for standard templates | planned multi-template adapter lane |

## Engineering constraints

- Planned lanes are not public APIs until packages are implemented, tested, documented, and added to the API snapshot.
- New public facade packages must follow the existing `internal/<domain>` plus `v<domain>` architecture.
- New public APIs must include godoc comments, deterministic tests, examples where reader-facing, API snapshot updates, and generated tool catalog updates.
- Security-sensitive lanes such as AI, FTP, SSH/SFTP, HTTP, crypto, JWT, archive, file, URL, config, random, ID, and DB helpers must use explicit errors and provider injection for tests.
- Ecosystem adapters should isolate optional dependencies behind narrow interfaces so unrelated users do not pay dependency or attack-surface costs.
- Benchmark documentation must describe baselines and commands, not claim universal performance wins.

## Validation gates for roadmap-driven work

Run focused tests first, then the standard governance gates for the touched area:

```bash
go test ./internal/<domain> ./v<domain>
UPDATE_API=1 make api-check
make docs-gen
make docs-check
make tools-check
make agent-check
make agent-security-check
```

Use `make bench-smoke` for benchmark-suite health and package-specific `go test -bench=. -benchmem -run=^$ ./<packages>` for benchmark baselines.

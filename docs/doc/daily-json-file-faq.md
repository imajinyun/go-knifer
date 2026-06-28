# Daily JSON/File FAQ

Use this FAQ when an application needs everyday JSON parsing, formatting, path access, file I/O, bounded reads, or file-type routing. Start with `vjson` for small in-memory JSON workflows and `vfile` for explicit filesystem side effects. Use Go's `encoding/json`, `os`, `io`, or `fs` directly when the caller needs lower-level control.

## Decision Matrix

| Task | First package | Use directly when | Guardrail |
| --- | --- | --- | --- |
| Encode or format small JSON values | `vjson` | `encoding/json.Encoder` settings, stream output, or custom numeric behavior must be visible. | Formatting is not validation; redact secrets before pretty-printing. |
| Decode into structs | `vjson` | `encoding/json.Decoder` options such as `DisallowUnknownFields` or `UseNumber` are required. | Validate size and business schema before trusting boundary input. |
| Read or update nested JSON fields | `vjson` | Typed request structs make the schema clearer. | Treat path strings as schema-dependent contracts and handle missing fields explicitly. |
| Bridge small XML payloads to JSON-shaped access | `vjson` | XML namespaces, attributes, streaming, or entity policy matters. | Validate XML inputs separately when they cross a trust boundary. |
| Read or write small trusted files | `vfile` | Platform flags, file descriptors, `fs.FS`, memory mapping, or syscall behavior matters. | Check every error and keep overwrite behavior visible. |
| Read large or untrusted content | `vfile` | A custom streaming pipeline already owns limits and backpressure. | Prefer bounded reads, chunked reads, and explicit size policy. |
| Route files by content family | `vfile` | A downstream parser or scanner owns validation. | Magic-number detection identifies type; it is not a safety boundary. |

## FAQ

### When should I use vjson instead of encoding/json?

Use `vjson` when the workflow is a common whole-value operation: compact encoding, pretty formatting, parsing into `JSONObject` or `JSONArray`, path reads and writes, or a small XML/JSON bridge. It keeps routine application code concise and gives examples a stable facade entry point.

Use `encoding/json` directly when the caller needs streaming, token inspection, multiple JSON values from one stream, `DisallowUnknownFields`, `UseNumber`, custom decoder state, or tight allocation control. In those cases, the decoder settings are part of the contract and should stay visible at the call site.

### When should I use vfile instead of os, io, or fs?

Use `vfile` when the operation is everyday file I/O with explicit errors: read a small text file, write generated content, append text, create parent directories, copy, delete, inspect names, read chunks, or detect a common file type from magic-number bytes.

Use `os`, `io`, or `fs` directly when the caller needs platform-specific flags, file descriptors, `fs.FS` traversal, memory mapping, custom buffering, or precise syscall behavior. `vfile` is a convenience facade, not a replacement for low-level filesystem control.

### How should I handle untrusted JSON input?

Bound the payload size before parsing data from networks, uploads, queues, partner callbacks, or user-provided files. If the input must obey a strict schema, decode into typed structs and validate required fields, allowed values, and unknown-field policy at the boundary.

Use `encoding/json.Decoder` directly when strict decoder behavior is part of the security or compatibility contract. Use `vjson` after the input is small enough and the workflow benefits from object, array, formatting, path, or XML bridge helpers.

### How should I handle untrusted file paths?

Keep user-provided names relative to a trusted base directory. Reject path traversal and do not let arbitrary absolute paths control writes, deletes, or overwrite behavior. For archive entries, use archive-specific safe extraction helpers instead of general string or file helpers.

Use `vfile` helpers once the path policy is clear. Path helpers and filename helpers do not authorize access by themselves; they make the operation easier to read after the trust boundary has been handled.

### When should I choose bounded reads or explicit errors?

Use bounded reads or chunked reads when content size is not controlled by your process. Whole-file helpers are appropriate for small configuration files, test fixtures, generated text, and trusted local inputs.

Prefer explicit error-returning helpers at boundaries. A failed parse, partial file write, cleanup failure, oversized input, or missing file should be classified by the caller rather than hidden behind a default value.

### Can JSON formatting or magic-number detection prove input is safe?

No. JSON formatting shows shape for humans but does not prove schema validity, authorization, or data sensitivity. Magic-number detection identifies common file families but does not prove that a file is benign, well-formed, small enough, or safe for a downstream parser.

Use formatting and detection as review and routing aids. Keep validation, redaction, size limits, parser limits, storage policy, and content scanning as separate boundary decisions.

### How should tests stay deterministic?

Use `vjson` provider options such as `WithParseDecoderFactory` when parsing behavior needs to be injected. Use temporary directories and provider-backed `vfile` options such as `WithOpen`, `WithOpenFile`, `WithStat`, `WithMkdirAll`, or `WithRemoveAll` when filesystem behavior should be hermetic.

Prefer `t.TempDir` or a narrow injected provider over shared package-level paths. Tests should prove behavior without relying on the developer's machine layout.

## Validation

Run these checks after changing Daily JSON/file FAQ guidance or governance metadata:

```bash
go test ./vjson ./vfile ./internal/json ./internal/file
make docs-check
make ai-context-check
make governance-maturity-check
```

The FAQ is governed by `daily_json_file_faq_governance` in `ai-context.json`; `make governance-maturity-check` verifies the document path, covered packages, required FAQ topics, scorecard FAQ status, and Sprint 25 roadmap state.

---
name: Feature request
about: Propose a new helper, package, or behavior change
title: "feat: "
labels: enhancement
assignees: ""
---

## Use case

Describe the user-facing problem this solves.

## Proposed package

Which existing `v*` package owns this capability?

If this needs a new package, explain why no existing package is a clear owner.

## Proposed API

```go
// Sketch the intended public API.
```

## Compatibility

- [ ] This does not rename or remove existing exported identifiers.
- [ ] This keeps the root `knifer` package free of business helpers.
- [ ] This keeps implementation logic in `internal/*`.

## Tests and docs

- [ ] I can describe at least one black-box facade test.
- [ ] I can describe at least one invalid-input or error-path test.
- [ ] I can provide a runnable example for new public API.

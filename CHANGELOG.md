# Changelog

All notable changes to this project are documented in this file.

This project follows [Semantic Versioning](https://semver.org/). Public
subpackage APIs are treated as the compatibility boundary.

## Unreleased

### Governance

- Added a repository security policy for private vulnerability reporting,
  supported versions, and security-sensitive package review areas.
- Added a coverage gate script so CI can enforce a measurable test baseline.
- Added facade coverage tests for `vurl`, `vzip`, `vdb`, `vjwt`, `vhttp`,
  `vnum`, `vtpl`, `vmap`, `vref`, `vobj`, `vmask`, `vresty`, and `vconf`.
- Added package-level coverage gates for security-sensitive facade packages:
  `vhttp`, `vresty`, `vconf`, and `vzip`.
- Documented release notes in a changelog so user-visible changes can be
  reviewed before tagging.

### Quality targets

- Current coverage gate baseline: 72.2%.
- Current security-sensitive package gates: `vhttp` >= 75%, `vresty` >= 65%,
  `vconf` >= 75%, and `vzip` >= 80%.
- Near-term target: 75% total statement coverage.
- Longer-term target: 80% total statement coverage, with priority on public
  facade packages and security-sensitive packages.

### Review focus

- Prioritize tests for `vhttp`, `vresty`, `vurl`, `vconf`, `vjwt`, `vzip`,
  `vcrypto`, `vdb`, and other packages that process untrusted input.
- Keep `v*` facade packages thin and preserve the `v* -> internal/*`
  dependency direction.

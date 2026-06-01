// Package vresty provides convenient HTTP client wrappers backed by resty.
//
// Prefer construction-time RequestOption values such as WithTimeout,
// WithHeader, WithFollowRedirects, WithRestyClient, WithUserAgent, and
// WithCookieDisabled for request-specific behavior. Global defaults are still
// available for compatibility, but per-call options avoid cross-request state
// coupling.
//
// This package only acts as a facade. Concrete implementations live in
// internal/resty.
package vresty

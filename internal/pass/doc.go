// Package pass provides password strength analysis helpers.
//
// The package evaluates passwords with a deterministic, local rule set. It
// does not call external services or maintain a large leaked-password corpus;
// callers that need breach checks should combine these helpers with their own
// allowlist or blocklist.
//
// Public package vpass re-exports this implementation for application code.
package pass

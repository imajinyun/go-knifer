// Package resty provides the internal implementation for the vresty package.
//
// This package builds chainable HTTP client utilities on top of resty.dev/v3.
// Keep lightweight standard-library HTTP helpers in internal/http and expose
// them through vhttp; use this package only for Resty-specific behavior.
// External modules should use the vresty facade instead.
package resty

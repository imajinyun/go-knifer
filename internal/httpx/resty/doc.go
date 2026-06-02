// Package resty provides the internal implementation for the vresty package.
//
// This package builds chainable HTTP client utilities on top of resty.dev/v3.
// Keep lightweight standard-library HTTP helpers in internal/http and expose
// them through vhttp; use this package only for Resty-specific behavior.
// External modules should use the vresty facade instead.
//
// Request defaults can be overridden per call with construction options, for
// example:
//
//	resp := resty.Get("https://example.com",
//	    resty.WithTimeout(3*time.Second),
//	    resty.WithHeader("X-Client", "go-knifer"),
//	).Execute()
package resty

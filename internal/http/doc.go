// Package http is aligned with hutool-http and provides HTTP client, download,
// Cookie, UserAgent, SimpleServer, and related utilities.
//
// This package is the standard-library based HTTP implementation for vhttp. Use
// internal/resty through vresty when a Resty-based chainable client is desired.
//
// Unlike hutool-http, this package wraps Go's standard net/http library and
// provides a chainable API:
//
//	body := http.Get("https://example.com").Execute().Body()
//	resp := http.NewRequest(http.MethodPost, url).
//	            Form(map[string]any{"a": 1}).
//	            Timeout(5 * time.Second).
//	            Execute()
package http

// Package gkhttp is aligned with hutool-http and provides HTTP client,
// download, Cookie, UserAgent, SimpleServer, and related utilities.
//
// Unlike hutool-http, this package wraps Go's standard net/http library and provides a chainable API:
//
//	body := gkhttp.Get("https://example.com").Execute().Body()
//	resp := gkhttp.NewRequest(gkhttp.MethodPost, url).
//	            Form(map[string]any{"a": 1}).
//	            Timeout(5 * time.Second).
//	            Execute()
package http

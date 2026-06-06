package vhttp

import (
	"time"

	httpx "github.com/imajinyun/go-knifer/internal/httpx/http"
)

// Get creates a GET request.
func Get(rawURL string, opts ...RequestOption) *Request { return httpx.Get(rawURL, opts...) }

// Post creates a POST request.
func Post(rawURL string, opts ...RequestOption) *Request { return httpx.Post(rawURL, opts...) }

// Put creates a PUT request.
func Put(rawURL string, opts ...RequestOption) *Request { return httpx.Put(rawURL, opts...) }

// Delete creates a DELETE request.
func Delete(rawURL string, opts ...RequestOption) *Request { return httpx.Delete(rawURL, opts...) }

// Patch creates a PATCH request.
func Patch(rawURL string, opts ...RequestOption) *Request { return httpx.Patch(rawURL, opts...) }

// Head creates a HEAD request.
func Head(rawURL string, opts ...RequestOption) *Request { return httpx.Head(rawURL, opts...) }

// Options delegates to the internal httpx implementation.
func Options(rawURL string, opts ...RequestOption) *Request {
	return httpx.Options(rawURL, opts...)
}

// NewRequest creates a request by method.
func NewRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewRequest(method, rawURL, opts...)
}

// NewIsolatedRequest creates a request without reading package-level global defaults.
func NewIsolatedRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewIsolatedRequest(method, rawURL, opts...)
}

// NewRequestWithConfig creates a request from an explicit global configuration snapshot.
func NewRequestWithConfig(method Method, rawURL string, cfg GlobalConfig, opts ...RequestOption) *Request {
	return httpx.NewRequestWithConfig(method, rawURL, cfg, opts...)
}

// CreateRequest delegates to the internal httpx implementation.
func CreateRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return CreateRequestWithOptions(method, rawURL, opts...)
}

// CreateRequestWithOptions delegates to the internal httpx implementation.
func CreateRequestWithOptions(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.CreateRequest(method, rawURL, opts...)
}

// CreateGet delegates to the internal httpx implementation.
func CreateGet(rawURL string, followRedirects bool) *Request {
	return httpx.CreateGet(rawURL, followRedirects)
}

// CreatePost delegates to the internal httpx implementation.
func CreatePost(rawURL string) *Request {
	return httpx.CreatePost(rawURL)
}

// GetString sends a GET request and returns response body as string.
func GetString(rawURL string) string { return GetStringWithOptions(rawURL) }

// GetStringWithOptions sends a GET request with options and returns response body as string.
func GetStringWithOptions(rawURL string, opts ...RequestOption) string {
	return httpx.GetStringWithOptions(rawURL, opts...)
}

// GetWithTimeout delegates to the internal httpx implementation.
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return GetWithTimeoutWithOptions(rawURL, timeout)
}

// GetWithTimeoutWithOptions delegates to the internal httpx implementation with options.
func GetWithTimeoutWithOptions(rawURL string, timeout time.Duration, opts ...RequestOption) string {
	return httpx.GetWithTimeoutWithOptions(rawURL, timeout, opts...)
}

// GetWithParams delegates to the internal httpx implementation.
func GetWithParams(rawURL string, params map[string]any) string {
	return GetWithParamsWithOptions(rawURL, params)
}

// GetWithParamsWithOptions delegates to the internal httpx implementation with options.
func GetWithParamsWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	return httpx.GetWithParamsWithOptions(rawURL, params, opts...)
}

// PostForm posts form parameters and returns response body as string.
func PostForm(rawURL string, params map[string]any) string {
	return PostFormWithOptions(rawURL, params)
}

// PostFormWithOptions posts form parameters with options and returns response body as string.
func PostFormWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	return httpx.PostFormWithOptions(rawURL, params, opts...)
}

// PostJSON posts JSON body and returns response body as string.
func PostJSON(rawURL, jsonStr string) string { return PostJSONWithOptions(rawURL, jsonStr) }

// PostJSONWithOptions posts JSON body with options and returns response body as string.
func PostJSONWithOptions(rawURL, jsonStr string, opts ...RequestOption) string {
	return httpx.PostJSONWithOptions(rawURL, jsonStr, opts...)
}

// PostString delegates to the internal httpx implementation.
func PostString(rawURL, body string) string {
	return PostStringWithOptions(rawURL, body)
}

// PostStringWithOptions delegates to the internal httpx implementation with options.
func PostStringWithOptions(rawURL, body string, opts ...RequestOption) string {
	return httpx.PostStringWithOptions(rawURL, body, opts...)
}

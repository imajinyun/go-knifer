package vhttp

import (
	"time"

	httpx "github.com/imajinyun/go-knifer/internal/httpx/http"
)

// Get creates a GET request.
//
// Security: Get is for trusted URLs. Use vresty.GetSafe or a custom transport
// policy when the URL is untrusted.
func Get(rawURL string, opts ...RequestOption) *Request { return httpx.Get(rawURL, opts...) }

// Post creates a POST request.
//
// Security: Post is for trusted URLs. Use vresty.PostSafe or a custom transport
// policy when the URL is untrusted.
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
//
// Security: NewRequest is for trusted URLs. Use vresty.NewSafeRequest or a
// custom transport policy when the URL is untrusted.
func NewRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewRequest(method, rawURL, opts...)
}

// NewIsolatedRequest creates a request without reading package-level global defaults.
func NewIsolatedRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewIsolatedRequest(method, rawURL, opts...)
}

// NewRequestWithConfig creates a request from an explicit global configuration snapshot.
//
// Security: NewRequestWithConfig is for trusted URLs. Use vresty.NewSafeRequest
// or a custom transport policy when the URL is untrusted.
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
	return CreateGetWithOptions(rawURL, followRedirects)
}

// CreateGetWithOptions delegates to the internal httpx implementation with options.
func CreateGetWithOptions(rawURL string, followRedirects bool, opts ...RequestOption) *Request {
	return httpx.CreateGetWithOptions(rawURL, followRedirects, opts...)
}

// CreatePost delegates to the internal httpx implementation.
func CreatePost(rawURL string) *Request {
	return CreatePostWithOptions(rawURL)
}

// CreatePostWithOptions delegates to the internal httpx implementation with options.
func CreatePostWithOptions(rawURL string, opts ...RequestOption) *Request {
	return httpx.CreatePostWithOptions(rawURL, opts...)
}

// GetString sends a GET request and returns response body as string.
func GetString(rawURL string) string { return GetStringWithOptions(rawURL) }

// GetStringWithOptions sends a GET request with options and returns response body as string.
func GetStringWithOptions(rawURL string, opts ...RequestOption) string {
	return httpx.GetStringWithOptions(rawURL, opts...)
}

// GetStringE sends a GET request and returns response body as string or an error.
func GetStringE(rawURL string) (string, error) { return GetStringEWithOptions(rawURL) }

// GetStringEWithOptions sends a GET request with options and returns response body as string or an error.
func GetStringEWithOptions(rawURL string, opts ...RequestOption) (string, error) {
	return httpx.GetStringEWithOptions(rawURL, opts...)
}

// GetWithTimeout delegates to the internal httpx implementation.
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return GetWithTimeoutWithOptions(rawURL, timeout)
}

// GetWithTimeoutWithOptions delegates to the internal httpx implementation with options.
func GetWithTimeoutWithOptions(rawURL string, timeout time.Duration, opts ...RequestOption) string {
	return httpx.GetWithTimeoutWithOptions(rawURL, timeout, opts...)
}

// GetWithTimeoutE sends a GET request with a timeout and returns response body or an error.
func GetWithTimeoutE(rawURL string, timeout time.Duration) (string, error) {
	return GetWithTimeoutEWithOptions(rawURL, timeout)
}

// GetWithTimeoutEWithOptions sends a GET request with a timeout and custom options, returning body or error.
func GetWithTimeoutEWithOptions(rawURL string, timeout time.Duration, opts ...RequestOption) (string, error) {
	return httpx.GetWithTimeoutEWithOptions(rawURL, timeout, opts...)
}

// GetWithParams delegates to the internal httpx implementation.
func GetWithParams(rawURL string, params map[string]any) string {
	return GetWithParamsWithOptions(rawURL, params)
}

// GetWithParamsWithOptions delegates to the internal httpx implementation with options.
func GetWithParamsWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	return httpx.GetWithParamsWithOptions(rawURL, params, opts...)
}

// GetWithParamsE sends a GET request with form parameters and returns response body or an error.
func GetWithParamsE(rawURL string, params map[string]any) (string, error) {
	return GetWithParamsEWithOptions(rawURL, params)
}

// GetWithParamsEWithOptions sends a GET request with form parameters and custom options, returning body or error.
func GetWithParamsEWithOptions(rawURL string, params map[string]any, opts ...RequestOption) (string, error) {
	return httpx.GetWithParamsEWithOptions(rawURL, params, opts...)
}

// PostForm posts form parameters and returns response body as string.
func PostForm(rawURL string, params map[string]any) string {
	return PostFormWithOptions(rawURL, params)
}

// PostFormWithOptions posts form parameters with options and returns response body as string.
func PostFormWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	return httpx.PostFormWithOptions(rawURL, params, opts...)
}

// PostFormE posts form parameters and returns response body or an error.
func PostFormE(rawURL string, params map[string]any) (string, error) {
	return PostFormEWithOptions(rawURL, params)
}

// PostFormEWithOptions posts form parameters with options and returns response body or an error.
func PostFormEWithOptions(rawURL string, params map[string]any, opts ...RequestOption) (string, error) {
	return httpx.PostFormEWithOptions(rawURL, params, opts...)
}

// PostJSON posts JSON body and returns response body as string.
func PostJSON(rawURL, jsonStr string) string { return PostJSONWithOptions(rawURL, jsonStr) }

// PostJSONWithOptions posts JSON body with options and returns response body as string.
func PostJSONWithOptions(rawURL, jsonStr string, opts ...RequestOption) string {
	return httpx.PostJSONWithOptions(rawURL, jsonStr, opts...)
}

// PostJSONE posts JSON body and returns response body or an error.
func PostJSONE(rawURL, jsonStr string) (string, error) { return PostJSONEWithOptions(rawURL, jsonStr) }

// PostJSONEWithOptions posts JSON body with options and returns response body or an error.
func PostJSONEWithOptions(rawURL, jsonStr string, opts ...RequestOption) (string, error) {
	return httpx.PostJSONEWithOptions(rawURL, jsonStr, opts...)
}

// PostString delegates to the internal httpx implementation.
func PostString(rawURL, body string) string {
	return PostStringWithOptions(rawURL, body)
}

// PostStringWithOptions delegates to the internal httpx implementation with options.
func PostStringWithOptions(rawURL, body string, opts ...RequestOption) string {
	return httpx.PostStringWithOptions(rawURL, body, opts...)
}

// PostStringE posts a string body and returns response body or an error.
func PostStringE(rawURL, body string) (string, error) { return PostStringEWithOptions(rawURL, body) }

// PostStringEWithOptions posts a string body with options and returns response body or an error.
func PostStringEWithOptions(rawURL, body string, opts ...RequestOption) (string, error) {
	return httpx.PostStringEWithOptions(rawURL, body, opts...)
}

package vhttp

import (
	"time"

	httpx "github.com/imajinyun/go-knifer/internal/httpx/http"
)

// Get creates a GET request.
//
// Security: Get is for trusted URLs. Use GetSafe when the URL is untrusted.
func Get(rawURL string, opts ...RequestOption) *Request { return httpx.Get(rawURL, opts...) }

// NewClient creates a request factory using the current global configuration snapshot.
func NewClient(opts ...ClientOption) *Client { return httpx.NewClient(opts...) }

// NewIsolatedClient creates a request factory without reading package-level global defaults.
func NewIsolatedClient(opts ...ClientOption) *Client { return httpx.NewIsolatedClient(opts...) }

// NewClientWithConfig creates a request factory from an explicit configuration snapshot.
func NewClientWithConfig(cfg GlobalConfig, opts ...RequestOption) *Client {
	return httpx.NewClientWithConfig(cfg, opts...)
}

// WithClientGlobalConfig sets the configuration snapshot used by a Client.
func WithClientGlobalConfig(cfg GlobalConfig) ClientOption { return httpx.WithClientGlobalConfig(cfg) }

// WithClientRequestOptions sets request options applied to every request created by a Client.
func WithClientRequestOptions(opts ...RequestOption) ClientOption {
	return httpx.WithClientRequestOptions(opts...)
}

// GetSafe creates a GET request with SSRF-oriented safety checks enabled.
func GetSafe(rawURL string, opts ...RequestOption) *Request { return httpx.GetSafe(rawURL, opts...) }

// Post creates a POST request.
//
// Security: Post is for trusted URLs. Use PostSafe when the URL is untrusted.
func Post(rawURL string, opts ...RequestOption) *Request { return httpx.Post(rawURL, opts...) }

// PostSafe creates a POST request with SSRF-oriented safety checks enabled.
func PostSafe(rawURL string, opts ...RequestOption) *Request { return httpx.PostSafe(rawURL, opts...) }

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
// Security: NewRequest is for trusted URLs. Use NewSafeRequest when the URL is untrusted.
func NewRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewRequest(method, rawURL, opts...)
}

// NewSafeRequest creates a request with SSRF-oriented safety checks enabled.
func NewSafeRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewSafeRequest(method, rawURL, opts...)
}

// NewIsolatedRequest creates a request without reading package-level global defaults.
func NewIsolatedRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewIsolatedRequest(method, rawURL, opts...)
}

// NewRequestWithConfig creates a request from an explicit global configuration snapshot.
//
// Security: NewRequestWithConfig is for trusted URLs. Use NewSafeRequest when the URL is untrusted.
func NewRequestWithConfig(method Method, rawURL string, cfg GlobalConfig, opts ...RequestOption) *Request {
	return httpx.NewRequestWithConfig(method, rawURL, cfg, opts...)
}

// CreateRequest delegates to the internal httpx implementation.
//
// Deprecated: use NewRequest for trusted URLs or NewSafeRequest for untrusted URLs.
func CreateRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return CreateRequestWithOptions(method, rawURL, opts...)
}

// CreateRequestWithOptions delegates to the internal httpx implementation.
//
// Deprecated: use NewRequest for trusted URLs or NewSafeRequest for untrusted URLs.
func CreateRequestWithOptions(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.CreateRequest(method, rawURL, opts...)
}

// CreateSafeRequest delegates to the internal httpx implementation with SSRF-oriented safety checks enabled.
//
// Deprecated: use NewSafeRequest.
func CreateSafeRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.CreateSafeRequest(method, rawURL, opts...)
}

// CreateGet delegates to the internal httpx implementation.
//
// Deprecated: use Get with WithFollowRedirects.
func CreateGet(rawURL string, followRedirects bool) *Request {
	return CreateGetWithOptions(rawURL, followRedirects)
}

// CreateGetWithOptions delegates to the internal httpx implementation with options.
//
// Deprecated: use Get with WithFollowRedirects.
func CreateGetWithOptions(rawURL string, followRedirects bool, opts ...RequestOption) *Request {
	return httpx.CreateGetWithOptions(rawURL, followRedirects, opts...)
}

// CreateGetSafe creates a GET request with SSRF-oriented safety checks enabled and sets whether redirects are followed.
//
// Deprecated: use GetSafe with WithFollowRedirects.
func CreateGetSafe(rawURL string, followRedirects bool, opts ...RequestOption) *Request {
	return httpx.CreateGetSafe(rawURL, followRedirects, opts...)
}

// CreatePost delegates to the internal httpx implementation.
//
// Deprecated: use Post for trusted URLs or PostSafe for untrusted URLs.
func CreatePost(rawURL string) *Request {
	return CreatePostWithOptions(rawURL)
}

// CreatePostWithOptions delegates to the internal httpx implementation with options.
//
// Deprecated: use Post.
func CreatePostWithOptions(rawURL string, opts ...RequestOption) *Request {
	return httpx.CreatePostWithOptions(rawURL, opts...)
}

// CreatePostSafe creates a POST request with SSRF-oriented safety checks enabled.
//
// Deprecated: use PostSafe.
func CreatePostSafe(rawURL string, opts ...RequestOption) *Request {
	return httpx.CreatePostSafe(rawURL, opts...)
}

// GetString sends a GET request and returns response body as string.
//
// Deprecated: use GetStringE to handle request and read errors explicitly.
func GetString(rawURL string) string { return GetStringWithOptions(rawURL) }

// GetStringWithOptions sends a GET request with options and returns response body as string.
//
// Deprecated: use GetStringEWithOptions to handle request and read errors explicitly.
func GetStringWithOptions(rawURL string, opts ...RequestOption) string {
	return httpx.GetStringWithOptions(rawURL, opts...)
}

// GetStringE sends a GET request and returns response body as string or an error.
func GetStringE(rawURL string) (string, error) { return GetStringEWithOptions(rawURL) }

// GetStringEWithOptions sends a GET request with options and returns response body as string or an error.
func GetStringEWithOptions(rawURL string, opts ...RequestOption) (string, error) {
	return httpx.GetStringEWithOptions(rawURL, opts...)
}

// GetStringSafeE sends a safe GET request and returns response body as string or an error.
func GetStringSafeE(rawURL string, opts ...RequestOption) (string, error) {
	return httpx.GetStringSafeE(rawURL, opts...)
}

// GetWithTimeout delegates to the internal httpx implementation.
//
// Deprecated: use GetWithTimeoutE to handle request and read errors explicitly.
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return GetWithTimeoutWithOptions(rawURL, timeout)
}

// GetWithTimeoutWithOptions delegates to the internal httpx implementation with options.
//
// Deprecated: use GetWithTimeoutEWithOptions to handle request and read errors explicitly.
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
//
// Deprecated: use GetWithParamsE to handle request and read errors explicitly.
func GetWithParams(rawURL string, params map[string]any) string {
	return GetWithParamsWithOptions(rawURL, params)
}

// GetWithParamsWithOptions delegates to the internal httpx implementation with options.
//
// Deprecated: use GetWithParamsEWithOptions to handle request and read errors explicitly.
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
//
// Deprecated: use PostFormE to handle request and read errors explicitly.
func PostForm(rawURL string, params map[string]any) string {
	return PostFormWithOptions(rawURL, params)
}

// PostFormWithOptions posts form parameters with options and returns response body as string.
//
// Deprecated: use PostFormEWithOptions to handle request and read errors explicitly.
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

// PostFormSafeE posts form parameters with SSRF-oriented safety checks enabled.
func PostFormSafeE(rawURL string, params map[string]any, opts ...RequestOption) (string, error) {
	return httpx.PostFormSafeE(rawURL, params, opts...)
}

// PostJSON posts JSON body and returns response body as string.
//
// Deprecated: use PostJSONE to handle request and read errors explicitly.
func PostJSON(rawURL, jsonStr string) string { return PostJSONWithOptions(rawURL, jsonStr) }

// PostJSONWithOptions posts JSON body with options and returns response body as string.
//
// Deprecated: use PostJSONEWithOptions to handle request and read errors explicitly.
func PostJSONWithOptions(rawURL, jsonStr string, opts ...RequestOption) string {
	return httpx.PostJSONWithOptions(rawURL, jsonStr, opts...)
}

// PostJSONE posts JSON body and returns response body or an error.
func PostJSONE(rawURL, jsonStr string) (string, error) { return PostJSONEWithOptions(rawURL, jsonStr) }

// PostJSONEWithOptions posts JSON body with options and returns response body or an error.
func PostJSONEWithOptions(rawURL, jsonStr string, opts ...RequestOption) (string, error) {
	return httpx.PostJSONEWithOptions(rawURL, jsonStr, opts...)
}

// PostJSONSafeE posts JSON body with SSRF-oriented safety checks enabled.
func PostJSONSafeE(rawURL, jsonStr string, opts ...RequestOption) (string, error) {
	return httpx.PostJSONSafeE(rawURL, jsonStr, opts...)
}

// PostString delegates to the internal httpx implementation.
//
// Deprecated: use PostStringE to handle request and read errors explicitly.
func PostString(rawURL, body string) string {
	return PostStringWithOptions(rawURL, body)
}

// PostStringWithOptions delegates to the internal httpx implementation with options.
//
// Deprecated: use PostStringEWithOptions to handle request and read errors explicitly.
func PostStringWithOptions(rawURL, body string, opts ...RequestOption) string {
	return httpx.PostStringWithOptions(rawURL, body, opts...)
}

// PostStringE posts a string body and returns response body or an error.
func PostStringE(rawURL, body string) (string, error) { return PostStringEWithOptions(rawURL, body) }

// PostStringEWithOptions posts a string body with options and returns response body or an error.
func PostStringEWithOptions(rawURL, body string, opts ...RequestOption) (string, error) {
	return httpx.PostStringEWithOptions(rawURL, body, opts...)
}

// PostStringSafeE posts a string body with SSRF-oriented safety checks enabled.
func PostStringSafeE(rawURL, body string, opts ...RequestOption) (string, error) {
	return httpx.PostStringSafeE(rawURL, body, opts...)
}

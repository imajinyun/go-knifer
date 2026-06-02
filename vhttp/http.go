package vhttp

import (
	"crypto/tls"
	"io"
	"io/fs"
	"net/http"
	"time"

	httpx "github.com/imajinyun/go-knifer/internal/http"
)

// Request is a chainable HTTP request builder.
type Request = httpx.HTTPRequest

// RequestOption customizes one HTTP request at construction time.
type RequestOption = httpx.RequestOption

// Response wraps an HTTP response.
type Response = httpx.HTTPResponse

// SaveOption customizes response file saving.
type SaveOption = httpx.SaveOption

// Method represents an HTTP method.
type Method = httpx.Method

// Header represents an HTTP header name.
type Header = httpx.Header

// ContentType represents an HTTP content type.
type ContentType = httpx.ContentType

// Error is the HTTP module error type.
type Error = httpx.HTTPError

// SimpleServer is a small HTTP server helper.
type SimpleServer = httpx.SimpleServer

// UserAgent describes parsed User-Agent information.
type UserAgent = httpx.UserAgent

const (
	// MethodGet is GET.
	MethodGet Method = httpx.MethodGet
	// MethodPost is POST.
	MethodPost Method = httpx.MethodPost
	// MethodPut is PUT.
	MethodPut Method = httpx.MethodPut
	// MethodDelete is DELETE.
	MethodDelete Method = httpx.MethodDelete
	// MethodPatch is PATCH.
	MethodPatch Method = httpx.MethodPatch
	// MethodHead is HEAD.
	MethodHead Method = httpx.MethodHead
	// MethodOptions is OPTIONS.
	MethodOptions Method = httpx.MethodOptions
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

// NewRequest creates a request by method.
func NewRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return httpx.NewRequest(method, rawURL, opts...)
}

// WithTimeout sets a per-request timeout.
func WithTimeout(d time.Duration) RequestOption { return httpx.WithTimeout(d) }

// WithHeader sets one per-request header.
func WithHeader(name, value string) RequestOption { return httpx.WithHeader(name, value) }

// WithHeaders sets per-request headers in batch.
func WithHeaders(headers map[string]string) RequestOption { return httpx.WithHeaders(headers) }

// WithFollowRedirects sets per-request redirect behavior.
func WithFollowRedirects(b bool) RequestOption { return httpx.WithFollowRedirects(b) }

// WithMaxRedirects sets the per-request redirect limit.
func WithMaxRedirects(n int) RequestOption { return httpx.WithMaxRedirects(n) }

// WithSkipTLSVerify sets per-request TLS verification behavior.
func WithSkipTLSVerify(b bool) RequestOption { return httpx.WithSkipTLSVerify(b) }

// WithTLSConfig sets a per-request TLS config.
func WithTLSConfig(cfg *tls.Config) RequestOption { return httpx.WithTLSConfig(cfg) }

// WithTransport sets a per-request RoundTripper.
func WithTransport(t http.RoundTripper) RequestOption { return httpx.WithTransport(t) }

// WithClient sets a per-request HTTP client.
func WithClient(c *http.Client) RequestOption { return httpx.WithClient(c) }

// WithCookieJar sets a per-request CookieJar. nil disables cookie management for this request.
func WithCookieJar(jar http.CookieJar) RequestOption { return httpx.WithCookieJar(jar) }

// WithUserAgent sets a per-request User-Agent.
func WithUserAgent(ua string) RequestOption { return httpx.WithUserAgent(ua) }

// WithContentType sets a per-request Content-Type at construction time.
func WithContentType(ct string) RequestOption { return httpx.WithContentType(ct) }

// WithCharset sets a per-request charset at construction time.
func WithCharset(charset string) RequestOption { return httpx.WithCharset(charset) }

// WithSaveFilePerm sets the file permission used when creating the destination file.
func WithSaveFilePerm(perm fs.FileMode) SaveOption { return httpx.WithSaveFilePerm(perm) }

// WithSaveDirPerm sets the directory permission used when creating parent directories.
func WithSaveDirPerm(perm fs.FileMode) SaveOption { return httpx.WithSaveDirPerm(perm) }

// WithSaveOverwrite controls whether an existing destination file may be replaced.
func WithSaveOverwrite(overwrite bool) SaveOption { return httpx.WithSaveOverwrite(overwrite) }

// WithSaveCreateParents controls whether parent directories are created automatically.
func WithSaveCreateParents(create bool) SaveOption { return httpx.WithSaveCreateParents(create) }

// WithSaveDefaultFilename sets the fallback file name used when dest is a directory.
func WithSaveDefaultFilename(name string) SaveOption { return httpx.WithSaveDefaultFilename(name) }

// GetString sends a GET request and returns response body as string.
func GetString(rawURL string) string { return httpx.GetString(rawURL) }

// PostForm posts form parameters and returns response body as string.
func PostForm(rawURL string, params map[string]any) string { return httpx.PostForm(rawURL, params) }

// PostJSON posts JSON body and returns response body as string.
func PostJSON(rawURL, jsonStr string) string { return httpx.PostJSON(rawURL, jsonStr) }

// Download downloads rawURL into w.
func Download(rawURL string, w io.Writer) (int64, error) { return httpx.Download(rawURL, w) }

// DownloadFile downloads rawURL to dest.
func DownloadFile(rawURL, dest string, opts ...SaveOption) (int64, error) {
	return httpx.DownloadFile(rawURL, dest, opts...)
}

// SetGlobalTimeout sets the global HTTP timeout.
func SetGlobalTimeout(d time.Duration) { httpx.SetGlobalTimeout(d) }

// GetGlobalTimeout returns the global HTTP timeout.
func GetGlobalTimeout() time.Duration { return httpx.GetGlobalTimeout() }

// SetGlobalHeader sets a global HTTP header.
func SetGlobalHeader(name, value string) { httpx.SetGlobalHeader(name, value) }

// AddGlobalHeader adds a global HTTP header value.
func AddGlobalHeader(name, value string) { httpx.AddGlobalHeader(name, value) }

// RemoveGlobalHeader removes a global HTTP header.
func RemoveGlobalHeader(name string) { httpx.RemoveGlobalHeader(name) }

// CloneGlobalHeaders returns cloned global headers.
func CloneGlobalHeaders() http.Header { return httpx.CloneGlobalHeaders() }

// BuildBasicAuth builds a Basic authorization value.
func BuildBasicAuth(user, pass string) string { return httpx.BuildBasicAuth(user, pass) }

// ToParams converts a map to query parameters.
func ToParams(m map[string]any) string { return httpx.ToParams(m) }

// ParseUserAgent parses a User-Agent string.
func ParseUserAgent(ua string) *UserAgent { return httpx.ParseUserAgent(ua) }

// NewSimpleServer creates a simple HTTP server on port.
func NewSimpleServer(port int) *SimpleServer { return httpx.NewSimpleServer(port) }

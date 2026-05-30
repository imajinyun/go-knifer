package vresty

import (
	"io"
	"time"

	restyimpl "github.com/imajinyun/go-knifer/internal/resty"
)

// Request is a chainable HTTP request builder backed by resty.
type Request = restyimpl.HTTPRequest

// Response wraps an HTTP response.
type Response = restyimpl.HTTPResponse

// Method represents an HTTP method.
type Method = restyimpl.Method

// Header represents an HTTP header name.
type Header = restyimpl.Header

// ContentType represents an HTTP content type.
type ContentType = restyimpl.ContentType

// HeaderValues stores HTTP header values.
type HeaderValues = restyimpl.HeaderValues

// Cookie contains a response cookie name and value.
type Cookie = restyimpl.Cookie

// Error is the HTTP module error type.
type Error = restyimpl.HTTPError

const (
	// MethodGet is GET.
	MethodGet Method = restyimpl.MethodGet
	// MethodPost is POST.
	MethodPost Method = restyimpl.MethodPost
	// MethodPut is PUT.
	MethodPut Method = restyimpl.MethodPut
	// MethodDelete is DELETE.
	MethodDelete Method = restyimpl.MethodDelete
	// MethodPatch is PATCH.
	MethodPatch Method = restyimpl.MethodPatch
	// MethodHead is HEAD.
	MethodHead Method = restyimpl.MethodHead
	// MethodOptions is OPTIONS.
	MethodOptions Method = restyimpl.MethodOptions
	// MethodTrace is TRACE.
	MethodTrace Method = restyimpl.MethodTrace
	// MethodConnect is CONNECT.
	MethodConnect Method = restyimpl.MethodConnect
)

// Get creates a GET request.
func Get(rawURL string) *Request { return restyimpl.Get(rawURL) }

// Post creates a POST request.
func Post(rawURL string) *Request { return restyimpl.Post(rawURL) }

// Put creates a PUT request.
func Put(rawURL string) *Request { return restyimpl.Put(rawURL) }

// Delete creates a DELETE request.
func Delete(rawURL string) *Request { return restyimpl.Delete(rawURL) }

// Patch creates a PATCH request.
func Patch(rawURL string) *Request { return restyimpl.Patch(rawURL) }

// Head creates a HEAD request.
func Head(rawURL string) *Request { return restyimpl.Head(rawURL) }

// Options creates an OPTIONS request.
func Options(rawURL string) *Request { return restyimpl.Options(rawURL) }

// NewRequest creates a request by method.
func NewRequest(method Method, rawURL string) *Request { return restyimpl.NewRequest(method, rawURL) }

// GetString sends a GET request and returns response body as string.
func GetString(rawURL string) string { return restyimpl.GetString(rawURL) }

// PostForm posts form parameters and returns response body as string.
func PostForm(rawURL string, params map[string]any) string { return restyimpl.PostForm(rawURL, params) }

// PostJSON posts JSON body and returns response body as string.
func PostJSON(rawURL, jsonStr string) string { return restyimpl.PostJSON(rawURL, jsonStr) }

// Download downloads rawURL into w.
func Download(rawURL string, w io.Writer) (int64, error) { return restyimpl.Download(rawURL, w) }

// DownloadFile downloads rawURL to dest.
func DownloadFile(rawURL, dest string) (int64, error) { return restyimpl.DownloadFile(rawURL, dest) }

// DownloadBytes downloads and returns bytes.
func DownloadBytes(rawURL string) []byte { return restyimpl.DownloadBytes(rawURL) }

// SetGlobalTimeout sets the global HTTP timeout.
func SetGlobalTimeout(d time.Duration) { restyimpl.SetGlobalTimeout(d) }

// GetGlobalTimeout returns the global HTTP timeout.
func GetGlobalTimeout() time.Duration { return restyimpl.GetGlobalTimeout() }

// SetGlobalHeader sets a global HTTP header.
func SetGlobalHeader(name, value string) { restyimpl.SetGlobalHeader(name, value) }

// AddGlobalHeader adds a global HTTP header value.
func AddGlobalHeader(name, value string) { restyimpl.AddGlobalHeader(name, value) }

// RemoveGlobalHeader removes a global HTTP header.
func RemoveGlobalHeader(name string) { restyimpl.RemoveGlobalHeader(name) }

// CloneGlobalHeaders returns cloned global headers.
func CloneGlobalHeaders() HeaderValues { return restyimpl.CloneGlobalHeaders() }

// CloseCookie disables global cookie management.
func CloseCookie() { restyimpl.CloseCookie() }

// BuildBasicAuth builds a Basic authorization value.
func BuildBasicAuth(user, pass string) string { return restyimpl.BuildBasicAuth(user, pass) }

// ToParams converts a map to query parameters.
func ToParams(m map[string]any) string { return restyimpl.ToParams(m) }

// IsHTTPS reports whether u uses HTTPS.
func IsHTTPS(u string) bool { return restyimpl.IsHTTPS(u) }

// IsHTTP reports whether u uses HTTP.
func IsHTTP(u string) bool { return restyimpl.IsHTTP(u) }

// BuildContentType builds a Content-Type string with charset.
func BuildContentType(contentType, charset string) string {
	return restyimpl.BuildContentType(contentType, charset)
}

// GuessContentType guesses Content-Type from the body.
func GuessContentType(body string) ContentType { return restyimpl.GuessContentType(body) }

// GetCharsetFromContentType extracts charset from Content-Type.
func GetCharsetFromContentType(ct string) string { return restyimpl.GetCharsetFromContentType(ct) }

// GetCharsetFromHTML extracts charset from HTML meta tags.
func GetCharsetFromHTML(html string) string { return restyimpl.GetCharsetFromHTML(html) }

// GetMimeType returns the MIME type by file extension.
func GetMimeType(filename string) string { return restyimpl.GetMimeType(filename) }

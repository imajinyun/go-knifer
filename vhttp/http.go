package vhttp

import (
	"io"
	"net/http"
	"time"

	httpx "github.com/imajinyun/go-knifer/internal/http"
)

// HTTPRequest is a chainable HTTP request builder.
type HTTPRequest = httpx.HTTPRequest

// HTTPResponse wraps an HTTP response.
type HTTPResponse = httpx.HTTPResponse

// HTTPMethod represents an HTTP method.
type HTTPMethod = httpx.Method

// Method represents an HTTP method.
type Method = httpx.Method

// HTTPHeader represents an HTTP header name.
type HTTPHeader = httpx.Header

// Header represents an HTTP header name.
type Header = httpx.Header

// HTTPContentType represents an HTTP content type.
type HTTPContentType = httpx.ContentType

// ContentType represents an HTTP content type.
type ContentType = httpx.ContentType

// HTTPError is the HTTP module error type.
type HTTPError = httpx.HTTPError

// SimpleServer is a small HTTP server helper.
type SimpleServer = httpx.SimpleServer

// UserAgent describes parsed User-Agent information.
type UserAgent = httpx.UserAgent

const (
	// HTTPMethodGet is GET.
	HTTPMethodGet HTTPMethod = httpx.MethodGet
	// HTTPMethodPost is POST.
	HTTPMethodPost HTTPMethod = httpx.MethodPost
	// HTTPMethodPut is PUT.
	HTTPMethodPut HTTPMethod = httpx.MethodPut
	// HTTPMethodDelete is DELETE.
	HTTPMethodDelete HTTPMethod = httpx.MethodDelete
	// HTTPMethodPatch is PATCH.
	HTTPMethodPatch HTTPMethod = httpx.MethodPatch
	// HTTPMethodHead is HEAD.
	HTTPMethodHead HTTPMethod = httpx.MethodHead
	// HTTPMethodOptions is OPTIONS.
	HTTPMethodOptions HTTPMethod = httpx.MethodOptions
)

// HTTPGet creates a GET request.
func HTTPGet(rawURL string) *HTTPRequest { return httpx.Get(rawURL) }

// HTTPPost creates a POST request.
func HTTPPost(rawURL string) *HTTPRequest { return httpx.Post(rawURL) }

// HTTPPut creates a PUT request.
func HTTPPut(rawURL string) *HTTPRequest { return httpx.Put(rawURL) }

// HTTPDelete creates a DELETE request.
func HTTPDelete(rawURL string) *HTTPRequest { return httpx.Delete(rawURL) }

// HTTPPatch creates a PATCH request.
func HTTPPatch(rawURL string) *HTTPRequest { return httpx.Patch(rawURL) }

// HTTPHead creates a HEAD request.
func HTTPHead(rawURL string) *HTTPRequest { return httpx.Head(rawURL) }

// HTTPNewRequest creates a request by method.
func HTTPNewRequest(method HTTPMethod, rawURL string) *HTTPRequest {
	return httpx.NewRequest(method, rawURL)
}

// HTTPGetString sends a GET request and returns response body as string.
func HTTPGetString(rawURL string) string { return httpx.GetString(rawURL) }

// HTTPPostForm posts form parameters and returns response body as string.
func HTTPPostForm(rawURL string, params map[string]any) string { return httpx.PostForm(rawURL, params) }

// HTTPPostJSON posts JSON body and returns response body as string.
func HTTPPostJSON(rawURL, jsonStr string) string { return httpx.PostJSON(rawURL, jsonStr) }

// HTTPDownload downloads rawURL into w.
func HTTPDownload(rawURL string, w io.Writer) (int64, error) { return httpx.Download(rawURL, w) }

// HTTPDownloadFile downloads rawURL to dest.
func HTTPDownloadFile(rawURL, dest string) (int64, error) { return httpx.DownloadFile(rawURL, dest) }

// HTTPSetGlobalTimeout sets the global HTTP timeout.
func HTTPSetGlobalTimeout(d time.Duration) { httpx.SetGlobalTimeout(d) }

// HTTPGetGlobalTimeout returns the global HTTP timeout.
func HTTPGetGlobalTimeout() time.Duration { return httpx.GetGlobalTimeout() }

// HTTPSetGlobalHeader sets a global HTTP header.
func HTTPSetGlobalHeader(name, value string) { httpx.SetGlobalHeader(name, value) }

// HTTPAddGlobalHeader adds a global HTTP header value.
func HTTPAddGlobalHeader(name, value string) { httpx.AddGlobalHeader(name, value) }

// HTTPRemoveGlobalHeader removes a global HTTP header.
func HTTPRemoveGlobalHeader(name string) { httpx.RemoveGlobalHeader(name) }

// HTTPCloneGlobalHeaders returns cloned global headers.
func HTTPCloneGlobalHeaders() http.Header { return httpx.CloneGlobalHeaders() }

// HTTPBuildBasicAuth builds a Basic authorization value.
func HTTPBuildBasicAuth(user, pass string) string { return httpx.BuildBasicAuth(user, pass) }

// HTTPToParams converts a map to query parameters.
func HTTPToParams(m map[string]any) string { return httpx.ToParams(m) }

// HTTPParseUserAgent parses a User-Agent string.
func HTTPParseUserAgent(ua string) *UserAgent { return httpx.ParseUserAgent(ua) }

// HTTPNewSimpleServer creates a simple HTTP server on port.
func HTTPNewSimpleServer(port int) *SimpleServer { return httpx.NewSimpleServer(port) }

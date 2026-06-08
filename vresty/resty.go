package vresty

import (
	"context"
	"crypto/tls"
	"io"
	"io/fs"
	"net"
	"os"
	"regexp"
	"time"

	restyimpl "github.com/imajinyun/go-knifer/internal/httpx/resty"
	grestry "resty.dev/v3"
)

// Request is a chainable HTTP request builder backed by resty.
type Request = restyimpl.HTTPRequest

// RequestOption customizes one HTTP request at construction time.
type RequestOption = restyimpl.RequestOption

// Response wraps an HTTP response.
type Response = restyimpl.HTTPResponse

// SaveOption customizes response file saving.
type SaveOption = restyimpl.SaveOption

// Method represents an HTTP method.
type Method = restyimpl.Method

// Header represents an HTTP header name.
type Header = restyimpl.Header

// ContentType represents an HTTP content type.
type ContentType = restyimpl.ContentType

// HeaderValues stores HTTP header values.
type HeaderValues = restyimpl.HeaderValues

// GlobalConfig captures resty package-level defaults for explicit request construction.
type GlobalConfig = restyimpl.GlobalConfig

// URLPolicy controls SSRF-oriented request validation for untrusted URLs.
type URLPolicy = restyimpl.URLPolicy

// CharsetOption customizes charset extraction helpers per call.
type CharsetOption = restyimpl.CharsetOption

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

const (
	HeaderAuthorization      Header = restyimpl.HeaderAuthorization
	HeaderProxyAuthorization Header = restyimpl.HeaderProxyAuthorization
	HeaderDate               Header = restyimpl.HeaderDate
	HeaderConnection         Header = restyimpl.HeaderConnection
	HeaderMimeVersion        Header = restyimpl.HeaderMimeVersion
	HeaderTrailer            Header = restyimpl.HeaderTrailer
	HeaderTransferEncoding   Header = restyimpl.HeaderTransferEncoding
	HeaderUpgrade            Header = restyimpl.HeaderUpgrade
	HeaderVia                Header = restyimpl.HeaderVia
	HeaderCacheControl       Header = restyimpl.HeaderCacheControl
	HeaderPragma             Header = restyimpl.HeaderPragma
	HeaderContentType        Header = restyimpl.HeaderContentType
	HeaderHost               Header = restyimpl.HeaderHost
	HeaderReferer            Header = restyimpl.HeaderReferer
	HeaderOrigin             Header = restyimpl.HeaderOrigin
	HeaderUserAgent          Header = restyimpl.HeaderUserAgent
	HeaderAccept             Header = restyimpl.HeaderAccept
	HeaderAcceptLanguage     Header = restyimpl.HeaderAcceptLanguage
	HeaderAcceptEncoding     Header = restyimpl.HeaderAcceptEncoding
	HeaderAcceptCharset      Header = restyimpl.HeaderAcceptCharset
	HeaderCookie             Header = restyimpl.HeaderCookie
	HeaderContentLength      Header = restyimpl.HeaderContentLength
	HeaderWWWAuthenticate    Header = restyimpl.HeaderWWWAuthenticate
	HeaderSetCookie          Header = restyimpl.HeaderSetCookie
	HeaderContentEncoding    Header = restyimpl.HeaderContentEncoding
	HeaderContentDisposition Header = restyimpl.HeaderContentDisposition
	HeaderETag               Header = restyimpl.HeaderETag
	HeaderLocation           Header = restyimpl.HeaderLocation
)

const (
	ContentTypeFormURLEncoded ContentType = restyimpl.ContentTypeFormURLEncoded
	ContentTypeMultipart      ContentType = restyimpl.ContentTypeMultipart
	ContentTypeJSON           ContentType = restyimpl.ContentTypeJSON
	ContentTypeXML            ContentType = restyimpl.ContentTypeXML
	ContentTypeTextPlain      ContentType = restyimpl.ContentTypeTextPlain
	ContentTypeTextXML        ContentType = restyimpl.ContentTypeTextXML
	ContentTypeTextHTML       ContentType = restyimpl.ContentTypeTextHTML
	ContentTypeOctetStream    ContentType = restyimpl.ContentTypeOctetStream
	ContentTypeEventStream    ContentType = restyimpl.ContentTypeEventStream
)

// Get creates a GET request.
//
// Security: Get is for trusted URLs. Use GetSafe for untrusted URLs.
func Get(rawURL string, opts ...RequestOption) *Request { return restyimpl.Get(rawURL, opts...) }

// GetSafe creates a GET request with SSRF-oriented safety checks enabled.
func GetSafe(rawURL string, opts ...RequestOption) *Request {
	return restyimpl.GetSafe(rawURL, opts...)
}

// Post creates a POST request.
//
// Security: Post is for trusted URLs. Use PostSafe for untrusted URLs.
func Post(rawURL string, opts ...RequestOption) *Request { return restyimpl.Post(rawURL, opts...) }

// PostSafe creates a POST request with SSRF-oriented safety checks enabled.
func PostSafe(rawURL string, opts ...RequestOption) *Request {
	return restyimpl.PostSafe(rawURL, opts...)
}

// Put creates a PUT request.
func Put(rawURL string, opts ...RequestOption) *Request { return restyimpl.Put(rawURL, opts...) }

// Delete creates a DELETE request.
func Delete(rawURL string, opts ...RequestOption) *Request { return restyimpl.Delete(rawURL, opts...) }

// Patch creates a PATCH request.
func Patch(rawURL string, opts ...RequestOption) *Request { return restyimpl.Patch(rawURL, opts...) }

// Head creates a HEAD request.
func Head(rawURL string, opts ...RequestOption) *Request { return restyimpl.Head(rawURL, opts...) }

// Options creates an OPTIONS request.
func Options(rawURL string, opts ...RequestOption) *Request {
	return restyimpl.Options(rawURL, opts...)
}

// NewRequest creates a request by method.
//
// Security: NewRequest is for trusted URLs unless callers provide WithURLPolicy
// with RejectPrivate enabled. Use NewSafeRequest for untrusted URLs.
func NewRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return restyimpl.NewRequest(method, rawURL, opts...)
}

// NewSafeRequest creates a request with SSRF-oriented safety checks enabled.
func NewSafeRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return restyimpl.NewSafeRequest(method, rawURL, opts...)
}

// NewIsolatedRequest creates a request without reading package-level global defaults.
func NewIsolatedRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return restyimpl.NewIsolatedRequest(method, rawURL, opts...)
}

// NewRequestWithConfig creates a request from an explicit global configuration snapshot.
//
// Security: NewRequestWithConfig is for trusted URLs unless callers provide
// WithURLPolicy with RejectPrivate enabled. Use NewSafeRequest for untrusted
// URLs.
func NewRequestWithConfig(method Method, rawURL string, cfg GlobalConfig, opts ...RequestOption) *Request {
	return restyimpl.NewRequestWithConfig(method, rawURL, cfg, opts...)
}

// CreateRequest creates a request by method.
func CreateRequest(method Method, rawURL string, opts ...RequestOption) *Request {
	return CreateRequestWithOptions(method, rawURL, opts...)
}

// CreateRequestWithOptions creates a request by method with per-call options.
func CreateRequestWithOptions(method Method, rawURL string, opts ...RequestOption) *Request {
	return restyimpl.CreateRequest(method, rawURL, opts...)
}

// CreateGet creates a GET request and sets whether redirects are followed.
func CreateGet(rawURL string, followRedirects bool) *Request {
	return CreateGetWithOptions(rawURL, followRedirects)
}

// CreateGetWithOptions creates a GET request with options and sets whether redirects are followed.
func CreateGetWithOptions(rawURL string, followRedirects bool, opts ...RequestOption) *Request {
	return restyimpl.CreateGetWithOptions(rawURL, followRedirects, opts...)
}

// CreatePost creates a POST request.
func CreatePost(rawURL string) *Request { return CreatePostWithOptions(rawURL) }

// CreatePostWithOptions creates a POST request with options.
func CreatePostWithOptions(rawURL string, opts ...RequestOption) *Request {
	return restyimpl.CreatePostWithOptions(rawURL, opts...)
}

// WithGlobalConfig initializes request defaults from a captured global configuration snapshot.
func WithGlobalConfig(cfg GlobalConfig) RequestOption { return restyimpl.WithGlobalConfig(cfg) }

// WithTimeout sets a per-request timeout.
func WithTimeout(d time.Duration) RequestOption { return restyimpl.WithTimeout(d) }

// WithHeader sets one per-request header.
func WithHeader(name, value string) RequestOption { return restyimpl.WithHeader(name, value) }

// WithHeaders sets per-request headers in batch.
func WithHeaders(headers map[string]string) RequestOption { return restyimpl.WithHeaders(headers) }

// WithFollowRedirects sets per-request redirect behavior.
func WithFollowRedirects(b bool) RequestOption { return restyimpl.WithFollowRedirects(b) }

// WithMaxRedirects sets the per-request redirect limit.
func WithMaxRedirects(n int) RequestOption { return restyimpl.WithMaxRedirects(n) }

// WithTLSConfig sets a per-request TLS config. It is ignored when WithRestyClient is set.
func WithTLSConfig(cfg *tls.Config) RequestOption { return restyimpl.WithTLSConfig(cfg) }

// WithRestyClient sets a per-request resty client and takes precedence over WithTLSConfig.
func WithRestyClient(c *grestry.Client) RequestOption { return restyimpl.WithRestyClient(c) }

// WithRestyClientFactory sets a per-request resty client factory.
func WithRestyClientFactory(factory func() *grestry.Client) RequestOption {
	return restyimpl.WithRestyClientFactory(factory)
}

// ConfigureDefaultRestyClientProvider sets the provider used to create resty clients when no per-request client is set.
func ConfigureDefaultRestyClientProvider(provider func() *grestry.Client) {
	restyimpl.ConfigureDefaultRestyClientProvider(provider)
}

// ResetDefaultRestyClientProvider restores resty.New as the default client provider.
func ResetDefaultRestyClientProvider() { restyimpl.ResetDefaultRestyClientProvider() }

// WithUserAgent sets a per-request User-Agent.
func WithUserAgent(ua string) RequestOption { return restyimpl.WithUserAgent(ua) }

// WithCookieDisabled sets per-request cookie management behavior.
func WithCookieDisabled(disabled bool) RequestOption { return restyimpl.WithCookieDisabled(disabled) }

// WithContentType sets a per-request Content-Type at construction time.
func WithContentType(ct string) RequestOption { return restyimpl.WithContentType(ct) }

// WithCharset sets a per-request charset at construction time.
func WithCharset(charset string) RequestOption { return restyimpl.WithCharset(charset) }

// WithJSONMarshalFunc sets the JSON marshal provider used by request body encoding.
func WithJSONMarshalFunc(marshal func(any) ([]byte, error)) RequestOption {
	return restyimpl.WithJSONMarshalFunc(marshal)
}

// WithJSONUnmarshalFunc sets the JSON unmarshal provider used by response decoding.
func WithJSONUnmarshalFunc(unmarshal func([]byte, any) error) RequestOption {
	return restyimpl.WithJSONUnmarshalFunc(unmarshal)
}

// WithJSONDecodeReadAllFunc sets the reader used before custom JSON unmarshalling.
func WithJSONDecodeReadAllFunc(readAll func(io.Reader) ([]byte, error)) RequestOption {
	return restyimpl.WithJSONDecodeReadAllFunc(readAll)
}

// WithMaxDecodeBytes limits bytes read before custom JSON unmarshalling. Non-positive means unlimited.
func WithMaxDecodeBytes(maxBytes int64) RequestOption {
	return restyimpl.WithMaxDecodeBytes(maxBytes)
}

// WithURLPolicy sets SSRF-oriented validation for the request URL and redirect targets.
func WithURLPolicy(policy URLPolicy) RequestOption { return restyimpl.WithURLPolicy(policy) }

// WithAllowedHosts restricts Safe requests to the provided host names.
func WithAllowedHosts(hosts ...string) RequestOption { return restyimpl.WithAllowedHosts(hosts...) }

// WithLookupIP sets the host resolver used by SSRF-oriented URL validation.
func WithLookupIP(lookupIP func(context.Context, string) ([]net.IP, error)) RequestOption {
	return restyimpl.WithLookupIP(lookupIP)
}

// WithSaveFilePerm sets the file permission used when creating the destination file.
func WithSaveFilePerm(perm fs.FileMode) SaveOption { return restyimpl.WithSaveFilePerm(perm) }

// WithSaveDirPerm sets the directory permission used when creating parent directories.
func WithSaveDirPerm(perm fs.FileMode) SaveOption { return restyimpl.WithSaveDirPerm(perm) }

// WithSaveOverwrite controls whether an existing destination file may be replaced.
func WithSaveOverwrite(overwrite bool) SaveOption { return restyimpl.WithSaveOverwrite(overwrite) }

// WithSaveCreateParents controls whether parent directories are created automatically.
func WithSaveCreateParents(create bool) SaveOption { return restyimpl.WithSaveCreateParents(create) }

// WithSaveDefaultFilename sets the fallback file name used when dest is a directory.
func WithSaveDefaultFilename(name string) SaveOption { return restyimpl.WithSaveDefaultFilename(name) }

// WithSaveStat sets the stat provider used to resolve directory destinations.
func WithSaveStat(stat func(string) (os.FileInfo, error)) SaveOption {
	return restyimpl.WithSaveStat(stat)
}

// WithSaveMkdirAll sets the directory creator used when saving responses.
func WithSaveMkdirAll(mkdirAll func(string, fs.FileMode) error) SaveOption {
	return restyimpl.WithSaveMkdirAll(mkdirAll)
}

// WithSaveOpenFile sets the file opener used when saving responses.
func WithSaveOpenFile(openFile func(string, int, fs.FileMode) (io.WriteCloser, error)) SaveOption {
	return restyimpl.WithSaveOpenFile(openFile)
}

// GetString sends a GET request and returns response body as string.
func GetString(rawURL string) string { return GetStringWithOptions(rawURL) }

// GetStringWithOptions sends a GET request with options and returns response body as string.
func GetStringWithOptions(rawURL string, opts ...RequestOption) string {
	return restyimpl.GetStringWithOptions(rawURL, opts...)
}

// GetWithTimeout sends a GET request with a timeout.
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return GetWithTimeoutWithOptions(rawURL, timeout)
}

// GetWithTimeoutWithOptions sends a GET request with a timeout and custom options.
func GetWithTimeoutWithOptions(rawURL string, timeout time.Duration, opts ...RequestOption) string {
	return restyimpl.GetWithTimeoutWithOptions(rawURL, timeout, opts...)
}

// GetWithParams sends a GET request with form parameters.
func GetWithParams(rawURL string, params map[string]any) string {
	return GetWithParamsWithOptions(rawURL, params)
}

// GetWithParamsWithOptions sends a GET request with form parameters and custom options.
func GetWithParamsWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	return restyimpl.GetWithParamsWithOptions(rawURL, params, opts...)
}

// PostString sends a POST request with a string body.
func PostString(rawURL, body string) string { return PostStringWithOptions(rawURL, body) }

// PostStringWithOptions sends a POST request with a string body and custom options.
func PostStringWithOptions(rawURL, body string, opts ...RequestOption) string {
	return restyimpl.PostStringWithOptions(rawURL, body, opts...)
}

// PostForm posts form parameters and returns response body as string.
func PostForm(rawURL string, params map[string]any) string {
	return PostFormWithOptions(rawURL, params)
}

// PostFormWithOptions posts form parameters with options and returns response body as string.
func PostFormWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	return restyimpl.PostFormWithOptions(rawURL, params, opts...)
}

// PostJSON posts JSON body and returns response body as string.
func PostJSON(rawURL, jsonStr string) string { return PostJSONWithOptions(rawURL, jsonStr) }

// PostJSONWithOptions posts JSON body with options and returns response body as string.
func PostJSONWithOptions(rawURL, jsonStr string, opts ...RequestOption) string {
	return restyimpl.PostJSONWithOptions(rawURL, jsonStr, opts...)
}

// Download downloads rawURL into w.
func Download(rawURL string, w io.Writer) (int64, error) { return DownloadWithOptions(rawURL, w) }

// DownloadWithOptions downloads rawURL into w with per-request options.
func DownloadWithOptions(rawURL string, w io.Writer, opts ...RequestOption) (int64, error) {
	return restyimpl.DownloadWithOptions(rawURL, w, opts...)
}

// DownloadFile downloads rawURL to dest.
func DownloadFile(rawURL, dest string, opts ...SaveOption) (int64, error) {
	return DownloadFileWithOptions(rawURL, dest, nil, opts...)
}

// DownloadFileWithOptions downloads rawURL to dest with per-request and per-save options.
func DownloadFileWithOptions(rawURL, dest string, requestOpts []RequestOption, saveOpts ...SaveOption) (int64, error) {
	return restyimpl.DownloadFileWithOptions(rawURL, dest, requestOpts, saveOpts...)
}

// DownloadBytes downloads and returns bytes.
func DownloadBytes(rawURL string) []byte { return DownloadBytesWithOptions(rawURL) }

// DownloadBytesWithOptions downloads and returns bytes with per-request options.
func DownloadBytesWithOptions(rawURL string, opts ...RequestOption) []byte {
	return restyimpl.DownloadBytesWithOptions(rawURL, opts...)
}

// DownloadString downloads remote text.
func DownloadString(rawURL, customCharset string) string {
	return DownloadStringWithOptions(rawURL, customCharset)
}

// DownloadStringWithOptions downloads remote text with per-request options.
func DownloadStringWithOptions(rawURL, customCharset string, opts ...RequestOption) string {
	return restyimpl.DownloadStringWithOptions(rawURL, customCharset, opts...)
}

// SetGlobalTimeout sets the global HTTP timeout.
func SetGlobalTimeout(d time.Duration) { restyimpl.SetGlobalTimeout(d) }

// GetGlobalTimeout returns the global HTTP timeout.
func GetGlobalTimeout() time.Duration { return restyimpl.GetGlobalTimeout() }

// SetGlobalMaxRedirects sets the global maximum redirect count.
func SetGlobalMaxRedirects(n int) { restyimpl.SetGlobalMaxRedirects(n) }

// GetGlobalMaxRedirects returns the global maximum redirect count.
func GetGlobalMaxRedirects() int { return restyimpl.GetGlobalMaxRedirects() }

// SetGlobalFollowRedirects sets whether redirects are followed globally.
func SetGlobalFollowRedirects(b bool) { restyimpl.SetGlobalFollowRedirects(b) }

// GetGlobalFollowRedirects reports whether redirects are followed globally.
func GetGlobalFollowRedirects() bool { return restyimpl.GetGlobalFollowRedirects() }

// SetGlobalUserAgent sets the global default User-Agent.
func SetGlobalUserAgent(ua string) { restyimpl.SetGlobalUserAgent(ua) }

// GetGlobalUserAgent returns the global default User-Agent.
func GetGlobalUserAgent() string { return restyimpl.GetGlobalUserAgent() }

// SnapshotGlobalConfig returns a copy of the package-level resty defaults.
func SnapshotGlobalConfig() GlobalConfig { return restyimpl.SnapshotGlobalConfig() }

// ResetGlobalConfig restores package-level resty defaults, including headers and cookies.
func ResetGlobalConfig() { restyimpl.ResetGlobalConfig() }

// ConfigureGlobalConfig replaces package-level resty defaults with cfg.
func ConfigureGlobalConfig(cfg GlobalConfig) { restyimpl.ConfigureGlobalConfig(cfg) }

// WithScopedGlobalConfig runs fn with cfg installed as package-level resty defaults,
// then restores the previous defaults.
func WithScopedGlobalConfig(cfg GlobalConfig, fn func()) { restyimpl.WithScopedGlobalConfig(cfg, fn) }

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

// WithCharsetRegexp sets the regexp used by GetCharsetFromContentTypeWithOptions.
func WithCharsetRegexp(re *regexp.Regexp) CharsetOption { return restyimpl.WithCharsetRegexp(re) }

// WithMetaCharsetRegexp sets the regexp used by GetCharsetFromHTMLWithOptions.
func WithMetaCharsetRegexp(re *regexp.Regexp) CharsetOption {
	return restyimpl.WithMetaCharsetRegexp(re)
}

// IsHTTPS reports whether the given URL is https.
func IsHTTPS(rawURL string) bool { return restyimpl.IsHTTPS(rawURL) }

// IsHTTP reports whether the given URL is http.
func IsHTTP(rawURL string) bool { return restyimpl.IsHTTP(rawURL) }

// ToParams converts a map to a URL query string.
func ToParams(m map[string]any) string { return restyimpl.ToParams(m) }

// URLWithForm appends form values to a URL.
func URLWithForm(rawURL string, form map[string]any) string {
	return restyimpl.URLWithForm(rawURL, form)
}

// BuildContentType builds a Content-Type string with charset.
func BuildContentType(contentType, charset string) string {
	return restyimpl.BuildContentType(contentType, charset)
}

// GuessContentType guesses Content-Type from the body.
func GuessContentType(body string) ContentType { return restyimpl.GuessContentType(body) }

// IsDefaultContentType reports whether the value is a default Content-Type.
func IsDefaultContentType(contentType string) bool {
	return restyimpl.IsDefaultContentType(contentType)
}

// IsFormURLEncoded reports whether the value is application/x-www-form-urlencoded.
func IsFormURLEncoded(contentType string) bool { return restyimpl.IsFormURLEncoded(contentType) }

// NewHTTPError creates an HTTP error.
func NewHTTPError(msg string, cause error) *Error { return restyimpl.NewHTTPError(msg, cause) }

// HTTPErrorf creates an HTTP error with a formatted message.
func HTTPErrorf(format string, args ...any) *Error { return restyimpl.HTTPErrorf(format, args...) }

// GetCharsetFromContentType extracts charset from Content-Type.
func GetCharsetFromContentType(ct string) string { return restyimpl.GetCharsetFromContentType(ct) }

// GetCharsetFromContentTypeWithOptions extracts charset from Content-Type with options.
func GetCharsetFromContentTypeWithOptions(ct string, opts ...CharsetOption) string {
	return restyimpl.GetCharsetFromContentTypeWithOptions(ct, opts...)
}

// GetCharsetFromHTML extracts charset from HTML meta tags.
func GetCharsetFromHTML(html string) string { return restyimpl.GetCharsetFromHTML(html) }

// GetCharsetFromHTMLWithOptions extracts charset from HTML meta tags with options.
func GetCharsetFromHTMLWithOptions(html string, opts ...CharsetOption) string {
	return restyimpl.GetCharsetFromHTMLWithOptions(html, opts...)
}

// GetMimeType returns the MIME type by file extension.
func GetMimeType(filename string) string { return restyimpl.GetMimeType(filename) }

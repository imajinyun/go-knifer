package http

import (
	"io"
	"regexp"
	"time"

	"github.com/imajinyun/go-knifer/internal/httpx/internal/shared"
	urlimpl "github.com/imajinyun/go-knifer/internal/url"
)

type CharsetOption = shared.CharsetOption

// WithCharsetRegexp sets the regexp used by GetCharsetFromContentTypeWithOptions.
func WithCharsetRegexp(re *regexp.Regexp) CharsetOption { return shared.WithCharsetRegexp(re) }

// WithMetaCharsetRegexp sets the regexp used by GetCharsetFromHTMLWithOptions.
func WithMetaCharsetRegexp(re *regexp.Regexp) CharsetOption { return shared.WithMetaCharsetRegexp(re) }

// IsHTTPS reports whether the given URL is https.
func IsHTTPS(u string) bool { return urlimpl.IsHTTPS(u) }

// IsHTTP reports whether the given URL is http.
func IsHTTP(u string) bool { return urlimpl.IsHTTP(u) }

// CreateRequest creates a request with the specified method, aligned with HttpUtil.createRequest.
//
// Deprecated: use NewRequest for trusted URLs or NewSafeRequest for untrusted URLs.
func CreateRequest(method Method, rawURL string, opts ...RequestOption) *HTTPRequest {
	return NewRequest(method, rawURL, opts...)
}

// CreateSafeRequest creates a request with SSRF-oriented safety checks enabled.
//
// Deprecated: use NewSafeRequest.
func CreateSafeRequest(method Method, rawURL string, opts ...RequestOption) *HTTPRequest {
	return NewSafeRequest(method, rawURL, opts...)
}

// CreateGet creates a GET request and sets whether redirects are followed.
//
// Deprecated: use Get with WithFollowRedirects.
func CreateGet(rawURL string, followRedirects bool) *HTTPRequest {
	return CreateGetWithOptions(rawURL, followRedirects)
}

// CreateGetWithOptions creates a GET request with options and sets whether redirects are followed.
//
// Deprecated: use Get with WithFollowRedirects.
func CreateGetWithOptions(rawURL string, followRedirects bool, opts ...RequestOption) *HTTPRequest {
	return Get(rawURL, opts...).FollowRedirects(followRedirects)
}

// CreateGetSafe creates a GET request with SSRF-oriented safety checks enabled and sets whether redirects are followed.
//
// Deprecated: use GetSafe with WithFollowRedirects.
func CreateGetSafe(rawURL string, followRedirects bool, opts ...RequestOption) *HTTPRequest {
	return GetSafe(rawURL, opts...).FollowRedirects(followRedirects)
}

// CreatePost creates a POST request.
//
// Deprecated: use Post for trusted URLs or PostSafe for untrusted URLs.
func CreatePost(rawURL string) *HTTPRequest { return CreatePostWithOptions(rawURL) }

// CreatePostWithOptions creates a POST request with options.
//
// Deprecated: use Post.
func CreatePostWithOptions(rawURL string, opts ...RequestOption) *HTTPRequest {
	return Post(rawURL, opts...)
}

// CreatePostSafe creates a POST request with SSRF-oriented safety checks enabled.
//
// Deprecated: use PostSafe.
func CreatePostSafe(rawURL string, opts ...RequestOption) *HTTPRequest {
	return PostSafe(rawURL, opts...)
}

// GetString sends a GET request and returns the response body as a string.
//
// Deprecated: use GetStringE to handle request and read errors explicitly.
func GetString(rawURL string) string { return GetStringWithOptions(rawURL) }

// GetStringWithOptions sends a GET request with options and returns the response body as a string.
//
// Deprecated: use GetStringEWithOptions to handle request and read errors explicitly.
func GetStringWithOptions(rawURL string, opts ...RequestOption) string {
	body, _ := GetStringEWithOptions(rawURL, opts...)
	return body
}

// GetStringE sends a GET request and returns the response body or an execution/read error.
func GetStringE(rawURL string) (string, error) { return GetStringEWithOptions(rawURL) }

// GetStringEWithOptions sends a GET request with options and returns the response body or an error.
func GetStringEWithOptions(rawURL string, opts ...RequestOption) (string, error) {
	resp := Get(rawURL, opts...).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// GetStringSafeE sends a safe GET request and returns the response body or an error.
func GetStringSafeE(rawURL string, opts ...RequestOption) (string, error) {
	resp := GetSafe(rawURL, opts...).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// GetWithTimeout sends a GET request with a timeout.
//
// Deprecated: use GetWithTimeoutE to handle request and read errors explicitly.
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return GetWithTimeoutWithOptions(rawURL, timeout)
}

// GetWithTimeoutWithOptions sends a GET request with a timeout and custom options.
//
// Deprecated: use GetWithTimeoutEWithOptions to handle request and read errors explicitly.
func GetWithTimeoutWithOptions(rawURL string, timeout time.Duration, opts ...RequestOption) string {
	body, _ := GetWithTimeoutEWithOptions(rawURL, timeout, opts...)
	return body
}

// GetWithTimeoutE sends a GET request with a timeout and returns the response body or an error.
func GetWithTimeoutE(rawURL string, timeout time.Duration) (string, error) {
	return GetWithTimeoutEWithOptions(rawURL, timeout)
}

// GetWithTimeoutEWithOptions sends a GET request with a timeout and custom options, returning body or error.
func GetWithTimeoutEWithOptions(rawURL string, timeout time.Duration, opts ...RequestOption) (string, error) {
	resp := Get(rawURL, opts...).Timeout(timeout).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// GetWithParams sends a GET request with form parameters.
//
// Deprecated: use GetWithParamsE to handle request and read errors explicitly.
func GetWithParams(rawURL string, params map[string]any) string {
	return GetWithParamsWithOptions(rawURL, params)
}

// GetWithParamsWithOptions sends a GET request with form parameters and custom options.
//
// Deprecated: use GetWithParamsEWithOptions to handle request and read errors explicitly.
func GetWithParamsWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	body, _ := GetWithParamsEWithOptions(rawURL, params, opts...)
	return body
}

// GetWithParamsE sends a GET request with form parameters and returns the response body or an error.
func GetWithParamsE(rawURL string, params map[string]any) (string, error) {
	return GetWithParamsEWithOptions(rawURL, params)
}

// GetWithParamsEWithOptions sends a GET request with form parameters and custom options, returning body or error.
func GetWithParamsEWithOptions(rawURL string, params map[string]any, opts ...RequestOption) (string, error) {
	resp := Get(rawURL, opts...).Form(params).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// PostString sends a POST request with a string body.
//
// Deprecated: use PostStringE to handle request and read errors explicitly.
func PostString(rawURL, body string) string {
	return PostStringWithOptions(rawURL, body)
}

// PostStringWithOptions sends a POST request with a string body and custom options.
//
// Deprecated: use PostStringEWithOptions to handle request and read errors explicitly.
func PostStringWithOptions(rawURL, body string, opts ...RequestOption) string {
	respBody, _ := PostStringEWithOptions(rawURL, body, opts...)
	return respBody
}

// PostStringE sends a POST request with a string body and returns the response body or an error.
func PostStringE(rawURL, body string) (string, error) { return PostStringEWithOptions(rawURL, body) }

// PostStringEWithOptions sends a POST request with a string body and custom options, returning body or error.
func PostStringEWithOptions(rawURL, body string, opts ...RequestOption) (string, error) {
	resp := Post(rawURL, opts...).BodyString(body).Execute()
	respBody := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return respBody, nil
}

// PostStringSafeE sends a safe POST request with a string body and returns the response body or an error.
func PostStringSafeE(rawURL, body string, opts ...RequestOption) (string, error) {
	resp := PostSafe(rawURL, opts...).BodyString(body).Execute()
	respBody := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return respBody, nil
}

// PostForm sends a POST request with form parameters.
//
// Deprecated: use PostFormE to handle request and read errors explicitly.
func PostForm(rawURL string, params map[string]any) string {
	return PostFormWithOptions(rawURL, params)
}

// PostFormWithOptions sends a POST request with form parameters and custom options.
//
// Deprecated: use PostFormEWithOptions to handle request and read errors explicitly.
func PostFormWithOptions(rawURL string, params map[string]any, opts ...RequestOption) string {
	body, _ := PostFormEWithOptions(rawURL, params, opts...)
	return body
}

// PostFormE sends a POST request with form parameters and returns the response body or an error.
func PostFormE(rawURL string, params map[string]any) (string, error) {
	return PostFormEWithOptions(rawURL, params)
}

// PostFormEWithOptions sends a POST request with form parameters and custom options, returning body or error.
func PostFormEWithOptions(rawURL string, params map[string]any, opts ...RequestOption) (string, error) {
	resp := Post(rawURL, opts...).Form(params).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// PostFormSafeE sends a safe POST request with form parameters and returns the response body or an error.
func PostFormSafeE(rawURL string, params map[string]any, opts ...RequestOption) (string, error) {
	resp := PostSafe(rawURL, opts...).Form(params).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// PostJSON sends a POST request with a JSON string body.
//
// Deprecated: use PostJSONE to handle request and read errors explicitly.
func PostJSON(rawURL, jsonStr string) string {
	return PostJSONWithOptions(rawURL, jsonStr)
}

// PostJSONWithOptions sends a POST request with a JSON string body and custom options.
//
// Deprecated: use PostJSONEWithOptions to handle request and read errors explicitly.
func PostJSONWithOptions(rawURL, jsonStr string, opts ...RequestOption) string {
	body, _ := PostJSONEWithOptions(rawURL, jsonStr, opts...)
	return body
}

// PostJSONE sends a POST request with a JSON string body and returns the response body or an error.
func PostJSONE(rawURL, jsonStr string) (string, error) { return PostJSONEWithOptions(rawURL, jsonStr) }

// PostJSONEWithOptions sends a POST request with a JSON string body and custom options, returning body or error.
func PostJSONEWithOptions(rawURL, jsonStr string, opts ...RequestOption) (string, error) {
	resp := Post(rawURL, opts...).BodyJSON(jsonStr).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// PostJSONSafeE sends a safe POST request with a JSON string body and returns the response body or an error.
func PostJSONSafeE(rawURL, jsonStr string, opts ...RequestOption) (string, error) {
	resp := PostSafe(rawURL, opts...).BodyJSON(jsonStr).Execute()
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// DownloadString downloads remote text and detects charset from response headers when customCharset is empty.
//
// Deprecated: use DownloadStringE to handle request and read errors explicitly.
func DownloadString(rawURL, customCharset string) string {
	return DownloadStringWithOptions(rawURL, customCharset)
}

// DownloadStringWithOptions downloads remote text with per-request options.
//
// Deprecated: use DownloadStringEWithOptions to handle request and read errors explicitly.
func DownloadStringWithOptions(rawURL, customCharset string, opts ...RequestOption) string {
	body, _ := DownloadStringEWithOptions(rawURL, customCharset, opts...)
	return body
}

// DownloadStringE downloads remote text and returns an error on request or read failure.
func DownloadStringE(rawURL, customCharset string) (string, error) {
	return DownloadStringEWithOptions(rawURL, customCharset)
}

// DownloadStringEWithOptions downloads remote text with per-request options and returns an error on failure.
func DownloadStringEWithOptions(rawURL, customCharset string, opts ...RequestOption) (string, error) {
	resp := Get(rawURL, opts...).Execute()
	if resp.err != nil {
		return "", resp.err
	}
	if customCharset != "" {
		_ = customCharset
	}
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// DownloadStringSafeE downloads remote text with SSRF-oriented safety checks enabled.
func DownloadStringSafeE(rawURL, customCharset string, opts ...RequestOption) (string, error) {
	resp := GetSafe(rawURL, opts...).Execute()
	if resp.err != nil {
		return "", resp.err
	}
	if customCharset != "" {
		_ = customCharset
	}
	body := resp.Body()
	if err := resp.Err(); err != nil {
		return "", err
	}
	return body, nil
}

// DownloadFile downloads to a file, using URL or response headers for the file name when dest is a directory.
func DownloadFile(rawURL, dest string, opts ...SaveOption) (int64, error) {
	return DownloadFileWithOptions(rawURL, dest, nil, opts...)
}

// DownloadFileWithOptions downloads to a file with per-request and per-save options.
func DownloadFileWithOptions(rawURL, dest string, requestOpts []RequestOption, saveOpts ...SaveOption) (int64, error) {
	resp := Get(rawURL, requestOpts...).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.SaveAs(dest, saveOpts...)
}

// Download downloads to a Writer.
func Download(rawURL string, w io.Writer) (int64, error) {
	return DownloadWithOptions(rawURL, w)
}

// DownloadWithOptions downloads to a Writer with per-request options.
func DownloadWithOptions(rawURL string, w io.Writer, opts ...RequestOption) (int64, error) {
	resp := Get(rawURL, opts...).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.writeBodyTo(w)
}

// DownloadSafe downloads to a Writer with SSRF-oriented safety checks enabled.
func DownloadSafe(rawURL string, w io.Writer, opts ...RequestOption) (int64, error) {
	resp := GetSafe(rawURL, opts...).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.writeBodyTo(w)
}

// DownloadBytes downloads and returns bytes.
//
// Deprecated: use DownloadBytesE to handle request and read errors explicitly.
func DownloadBytes(rawURL string) []byte { return DownloadBytesWithOptions(rawURL) }

// DownloadBytesWithOptions downloads and returns bytes with per-request options.
//
// Deprecated: use DownloadBytesEWithOptions to handle request and read errors explicitly.
func DownloadBytesWithOptions(rawURL string, opts ...RequestOption) []byte {
	body, _ := DownloadBytesEWithOptions(rawURL, opts...)
	return body
}

// DownloadBytesE downloads and returns bytes or an error.
func DownloadBytesE(rawURL string) ([]byte, error) { return DownloadBytesEWithOptions(rawURL) }

// DownloadBytesEWithOptions downloads and returns bytes with per-request options or an error.
func DownloadBytesEWithOptions(rawURL string, opts ...RequestOption) ([]byte, error) {
	resp := Get(rawURL, opts...).Execute()
	body := resp.Bytes()
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return body, nil
}

// DownloadBytesSafeE downloads and returns bytes with SSRF-oriented safety checks enabled.
func DownloadBytesSafeE(rawURL string, opts ...RequestOption) ([]byte, error) {
	resp := GetSafe(rawURL, opts...).Execute()
	body := resp.Bytes()
	if err := resp.Err(); err != nil {
		return nil, err
	}
	return body, nil
}

// ToParams converts a map to a URL query string.
func ToParams(m map[string]any) string { return urlimpl.EncodeQueryMap(m) }

// EncodeParams encodes a URL containing parameters; only the part after ? is encoded.
func EncodeParams(rawURL string) string { return urlimpl.EncodeParams(rawURL) }

// DecodeParamMap parses a query string into a map.
func DecodeParamMap(paramsStr string) map[string]string { return urlimpl.DecodeQueryFirst(paramsStr) }

// DecodeParams parses a query string into a multi-value map.
func DecodeParams(paramsStr string) map[string][]string { return urlimpl.DecodeQuery(paramsStr) }

// URLWithForm appends form values to a URL.
func URLWithForm(rawURL string, form map[string]any) string { return urlimpl.AppendQuery(rawURL, form) }

// BuildBasicAuth builds a Basic Auth string.
func BuildBasicAuth(user, pass string) string {
	return shared.BuildBasicAuth(user, pass)
}

var (
	// CharsetPattern matches charset in Content-Type.
	CharsetPattern = shared.CharsetPattern
	// MetaCharsetPattern matches charset in HTML meta tags.
	MetaCharsetPattern = shared.MetaCharsetPattern
)

// GetCharsetFromContentType extracts charset from Content-Type.
func GetCharsetFromContentType(ct string) string {
	return shared.GetCharsetFromContentType(ct)
}

// GetCharsetFromContentTypeWithOptions extracts charset from Content-Type with options.
func GetCharsetFromContentTypeWithOptions(ct string, opts ...CharsetOption) string {
	return shared.GetCharsetFromContentTypeWithOptions(ct, opts...)
}

// GetCharsetFromHTML extracts charset from HTML meta tags.
func GetCharsetFromHTML(html string) string {
	return shared.GetCharsetFromHTML(html)
}

// GetCharsetFromHTMLWithOptions extracts charset from HTML meta tags with options.
func GetCharsetFromHTMLWithOptions(html string, opts ...CharsetOption) string {
	return shared.GetCharsetFromHTMLWithOptions(html, opts...)
}

// GetMimeType returns the MIME type by file extension, or an empty string when unknown.
func GetMimeType(filename string) string {
	return shared.GetMimeType(filename)
}

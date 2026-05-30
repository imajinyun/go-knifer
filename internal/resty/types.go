package resty

import (
	"fmt"
	"strings"
)

// Method represents an HTTP request method.
type Method string

const (
	MethodGet     Method = "GET"
	MethodPost    Method = "POST"
	MethodHead    Method = "HEAD"
	MethodOptions Method = "OPTIONS"
	MethodPut     Method = "PUT"
	MethodDelete  Method = "DELETE"
	MethodTrace   Method = "TRACE"
	MethodConnect Method = "CONNECT"
	MethodPatch   Method = "PATCH"
)

// String returns the method string.
func (m Method) String() string { return string(m) }

// Header defines common HTTP header names.
type Header string

const (
	HeaderAuthorization      Header = "Authorization"
	HeaderProxyAuthorization Header = "Proxy-Authorization"
	HeaderDate               Header = "Date"
	HeaderConnection         Header = "Connection"
	HeaderMimeVersion        Header = "MIME-Version"
	HeaderTrailer            Header = "Trailer"
	HeaderTransferEncoding   Header = "Transfer-Encoding"
	HeaderUpgrade            Header = "Upgrade"
	HeaderVia                Header = "Via"
	HeaderCacheControl       Header = "Cache-Control"
	HeaderPragma             Header = "Pragma"
	HeaderContentType        Header = "Content-Type"
	HeaderHost               Header = "Host"
	HeaderReferer            Header = "Referer"
	HeaderOrigin             Header = "Origin"
	HeaderUserAgent          Header = "User-Agent"
	HeaderAccept             Header = "Accept"
	HeaderAcceptLanguage     Header = "Accept-Language"
	HeaderAcceptEncoding     Header = "Accept-Encoding"
	HeaderAcceptCharset      Header = "Accept-Charset"
	HeaderCookie             Header = "Cookie"
	HeaderContentLength      Header = "Content-Length"
	HeaderWWWAuthenticate    Header = "WWW-Authenticate"
	HeaderSetCookie          Header = "Set-Cookie"
	HeaderContentEncoding    Header = "Content-Encoding"
	HeaderContentDisposition Header = "Content-Disposition"
	HeaderETag               Header = "ETag"
	HeaderLocation           Header = "Location"
)

// String returns the header name.
func (h Header) String() string { return string(h) }

// ContentType defines common Content-Type values.
type ContentType string

const (
	ContentTypeFormURLEncoded ContentType = "application/x-www-form-urlencoded"
	ContentTypeMultipart      ContentType = "multipart/form-data"
	ContentTypeJSON           ContentType = "application/json"
	ContentTypeXML            ContentType = "application/xml"
	ContentTypeTextPlain      ContentType = "text/plain"
	ContentTypeTextXML        ContentType = "text/xml"
	ContentTypeTextHTML       ContentType = "text/html"
	ContentTypeOctetStream    ContentType = "application/octet-stream"
	ContentTypeEventStream    ContentType = "text/event-stream"
)

// String returns the Content-Type string.
func (c ContentType) String() string { return string(c) }

// WithCharset returns the Content-Type string with charset.
func (c ContentType) WithCharset(charset string) string {
	return BuildContentType(string(c), charset)
}

// BuildContentType builds a Content-Type string with charset.
func BuildContentType(contentType, charset string) string {
	if charset == "" {
		return contentType
	}
	return fmt.Sprintf("%s;charset=%s", contentType, charset)
}

// IsDefaultContentType reports whether the value is a default Content-Type.
func IsDefaultContentType(contentType string) bool {
	return contentType == "" || IsFormURLEncoded(contentType)
}

// IsFormURLEncoded reports whether the value is application/x-www-form-urlencoded.
func IsFormURLEncoded(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), string(ContentTypeFormURLEncoded))
}

// GuessContentType guesses Content-Type from the first body character.
func GuessContentType(body string) ContentType {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	switch body[0] {
	case '{', '[':
		return ContentTypeJSON
	case '<':
		return ContentTypeXML
	}
	return ""
}

// HTTPError represents an error during HTTP operations.
type HTTPError struct {
	Msg   string
	Cause error
}

// Error returns the error message.
func (e *HTTPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Cause)
	}
	return e.Msg
}

// Unwrap returns the underlying error.
func (e *HTTPError) Unwrap() error { return e.Cause }

// NewHTTPError creates an HTTP error.
func NewHTTPError(msg string, cause error) *HTTPError { return &HTTPError{Msg: msg, Cause: cause} }

// HTTPErrorf creates an HTTP error with a formatted message.
func HTTPErrorf(format string, args ...any) *HTTPError {
	return &HTTPError{Msg: fmt.Sprintf(format, args...)}
}

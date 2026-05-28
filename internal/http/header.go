package http

// Header 常见 HTTP 头域名称（对应 hutool-http Header 枚举）。
type Header string

const (
	// 通用头域
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

	// 请求头域
	HeaderHost           Header = "Host"
	HeaderReferer        Header = "Referer"
	HeaderOrigin         Header = "Origin"
	HeaderUserAgent      Header = "User-Agent"
	HeaderAccept         Header = "Accept"
	HeaderAcceptLanguage Header = "Accept-Language"
	HeaderAcceptEncoding Header = "Accept-Encoding"
	HeaderAcceptCharset  Header = "Accept-Charset"
	HeaderCookie         Header = "Cookie"
	HeaderContentLength  Header = "Content-Length"

	// 响应头域
	HeaderWWWAuthenticate    Header = "WWW-Authenticate"
	HeaderSetCookie          Header = "Set-Cookie"
	HeaderContentEncoding    Header = "Content-Encoding"
	HeaderContentDisposition Header = "Content-Disposition"
	HeaderETag               Header = "ETag"
	HeaderLocation           Header = "Location"
)

// String 返回头域字符串。
func (h Header) String() string { return string(h) }

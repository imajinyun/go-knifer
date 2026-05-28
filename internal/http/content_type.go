package http

import (
	"fmt"
	"strings"
)

// ContentType 常用 Content-Type 类型（对应 hutool-http ContentType 枚举）。
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

// String 返回 Content-Type 字符串。
func (c ContentType) String() string { return string(c) }

// WithCharset 输出 Content-Type 字符串，附带编码信息，如 "application/json;charset=UTF-8"。
func (c ContentType) WithCharset(charset string) string {
	return BuildContentType(string(c), charset)
}

// BuildContentType 构造 Content-Type 字符串，附带编码。
func BuildContentType(contentType, charset string) string {
	if charset == "" {
		return contentType
	}
	return fmt.Sprintf("%s;charset=%s", contentType, charset)
}

// IsDefaultContentType 判断是否为默认 Content-Type，包括 "" 和 form-urlencoded。
func IsDefaultContentType(contentType string) bool {
	return contentType == "" || IsFormURLEncoded(contentType)
}

// IsFormURLEncoded 判断是否为 application/x-www-form-urlencoded。
func IsFormURLEncoded(contentType string) bool {
	return strings.HasPrefix(strings.ToLower(contentType), string(ContentTypeFormURLEncoded))
}

// GuessContentType 从请求体首字符猜测 Content-Type，仅支持 JSON/XML。
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

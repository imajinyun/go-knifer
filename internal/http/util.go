package http

import (
	"encoding/base64"
	"io"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// IsHTTPS 判断给定 URL 是否为 https。
func IsHTTPS(u string) bool { return strings.HasPrefix(strings.ToLower(u), "https:") }

// IsHTTP 判断给定 URL 是否为 http。
func IsHTTP(u string) bool { return strings.HasPrefix(strings.ToLower(u), "http:") }

// CreateRequest 创建指定方法的请求（对应 HttpUtil.createRequest）。
func CreateRequest(method Method, rawURL string) *HTTPRequest { return NewRequest(method, rawURL) }

// CreateGet 创建 GET 请求并设置是否跟随重定向。
func CreateGet(rawURL string, followRedirects bool) *HTTPRequest {
	return Get(rawURL).FollowRedirects(followRedirects)
}

// CreatePost 创建 POST 请求。
func CreatePost(rawURL string) *HTTPRequest { return Post(rawURL) }

// GetString 直接 GET 并返回响应字符串。
func GetString(rawURL string) string { return Get(rawURL).Execute().Body() }

// GetWithTimeout 带超时的 GET。
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return Get(rawURL).Timeout(timeout).Execute().Body()
}

// GetWithParams 带 form 参数的 GET。
func GetWithParams(rawURL string, params map[string]any) string {
	return Get(rawURL).Form(params).Execute().Body()
}

// PostString 以字符串 body 发起 POST。
func PostString(rawURL, body string) string {
	return Post(rawURL).BodyString(body).Execute().Body()
}

// PostForm 以表单方式发起 POST。
func PostForm(rawURL string, params map[string]any) string {
	return Post(rawURL).Form(params).Execute().Body()
}

// PostJSON 以 JSON 字符串发起 POST。
func PostJSON(rawURL, jsonStr string) string {
	return Post(rawURL).BodyJSON(jsonStr).Execute().Body()
}

// DownloadString 下载远程文本，customCharset 为空时按响应头识别。
func DownloadString(rawURL, customCharset string) string {
	resp := Get(rawURL).Execute()
	if resp.err != nil {
		return ""
	}
	if customCharset != "" {
		// Go 不内置编码转换，只直接按字节返回；调用方按需转码
		_ = customCharset
	}
	return resp.Body()
}

// DownloadFile 下载到文件，dest 为目录时取 URL 或响应头中的文件名。
func DownloadFile(rawURL, dest string) (int64, error) {
	resp := Get(rawURL).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.SaveAs(dest)
}

// Download 下载到 Writer。
func Download(rawURL string, w io.Writer) (int64, error) {
	resp := Get(rawURL).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.WriteTo(w)
}

// DownloadBytes 下载并返回字节数据。
func DownloadBytes(rawURL string) []byte { return Get(rawURL).Execute().Bytes() }

// ToParams 将 map 转为 URL Query 字符串。
func ToParams(m map[string]any) string {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, toString(v))
	}
	return values.Encode()
}

// EncodeParams 对包含参数的 URL 进行编码（仅对 ?后的部分）。
func EncodeParams(rawURL string) string {
	idx := strings.Index(rawURL, "?")
	if idx < 0 {
		return rawURL
	}
	pre := rawURL[:idx]
	q := rawURL[idx+1:]
	values, err := url.ParseQuery(q)
	if err != nil {
		return rawURL
	}
	return pre + "?" + values.Encode()
}

// DecodeParamMap 将 query 字符串解析为 map。
func DecodeParamMap(paramsStr string) map[string]string {
	out := map[string]string{}
	values, err := url.ParseQuery(paramsStr)
	if err != nil {
		return out
	}
	for k, vs := range values {
		if len(vs) > 0 {
			out[k] = vs[0]
		}
	}
	return out
}

// DecodeParams 将 query 字符串解析为多值 map。
func DecodeParams(paramsStr string) map[string][]string {
	values, err := url.ParseQuery(paramsStr)
	if err != nil {
		return map[string][]string{}
	}
	out := map[string][]string{}
	for k, v := range values {
		out[k] = v
	}
	return out
}

// URLWithForm 将 form 拼接到 URL 中。
func URLWithForm(rawURL string, form map[string]any) string {
	encoded := ToParams(form)
	if encoded == "" {
		return rawURL
	}
	if strings.Contains(rawURL, "?") {
		if strings.HasSuffix(rawURL, "&") || strings.HasSuffix(rawURL, "?") {
			return rawURL + encoded
		}
		return rawURL + "&" + encoded
	}
	return rawURL + "?" + encoded
}

// BuildBasicAuth 构造 Basic Auth 字符串。
func BuildBasicAuth(user, pass string) string {
	token := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	return "Basic " + token
}

var (
	// CharsetPattern Content-Type 中匹配 charset 的正则。
	CharsetPattern = regexp.MustCompile(`(?i)charset\s*=\s*([a-z0-9-]+)`)
	// MetaCharsetPattern HTML meta 标签中匹配 charset 的正则。
	MetaCharsetPattern = regexp.MustCompile(`(?i)<meta[^>]*?charset\s*=\s*['"]?([a-z0-9-]+)`)
)

// GetCharsetFromContentType 从 Content-Type 中提取字符集。
func GetCharsetFromContentType(ct string) string {
	m := CharsetPattern.FindStringSubmatch(ct)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// GetCharsetFromHTML 从 HTML meta 中提取 charset。
func GetCharsetFromHTML(html string) string {
	m := MetaCharsetPattern.FindStringSubmatch(html)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// GetMimeType 根据文件扩展名返回 MIME 类型，未知返回空字符串。
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return ""
	}
	if v, ok := mimeTypes[ext]; ok {
		return v
	}
	return ""
}

// 简化版 MIME 表
var mimeTypes = map[string]string{
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".xml":  "application/xml",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".pdf":  "application/pdf",
	".zip":  "application/zip",
	".gz":   "application/gzip",
	".tar":  "application/x-tar",
	".txt":  "text/plain",
	".csv":  "text/csv",
	".mp4":  "video/mp4",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
}

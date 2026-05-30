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

// IsHTTPS reports whether the given URL is https.
func IsHTTPS(u string) bool { return strings.HasPrefix(strings.ToLower(u), "https:") }

// IsHTTP reports whether the given URL is http.
func IsHTTP(u string) bool { return strings.HasPrefix(strings.ToLower(u), "http:") }

// CreateRequest creates a request with the specified method, aligned with HttpUtil.createRequest.
func CreateRequest(method Method, rawURL string) *HTTPRequest { return NewRequest(method, rawURL) }

// CreateGet creates a GET request and sets whether redirects are followed.
func CreateGet(rawURL string, followRedirects bool) *HTTPRequest {
	return Get(rawURL).FollowRedirects(followRedirects)
}

// CreatePost creates a POST request.
func CreatePost(rawURL string) *HTTPRequest { return Post(rawURL) }

// GetString sends a GET request and returns the response body as a string.
func GetString(rawURL string) string { return Get(rawURL).Execute().Body() }

// GetWithTimeout sends a GET request with a timeout.
func GetWithTimeout(rawURL string, timeout time.Duration) string {
	return Get(rawURL).Timeout(timeout).Execute().Body()
}

// GetWithParams sends a GET request with form parameters.
func GetWithParams(rawURL string, params map[string]any) string {
	return Get(rawURL).Form(params).Execute().Body()
}

// PostString sends a POST request with a string body.
func PostString(rawURL, body string) string {
	return Post(rawURL).BodyString(body).Execute().Body()
}

// PostForm sends a POST request with form parameters.
func PostForm(rawURL string, params map[string]any) string {
	return Post(rawURL).Form(params).Execute().Body()
}

// PostJSON sends a POST request with a JSON string body.
func PostJSON(rawURL, jsonStr string) string {
	return Post(rawURL).BodyJSON(jsonStr).Execute().Body()
}

// DownloadString downloads remote text and detects charset from response headers when customCharset is empty.
func DownloadString(rawURL, customCharset string) string {
	resp := Get(rawURL).Execute()
	if resp.err != nil {
		return ""
	}
	if customCharset != "" {
		// Go does not provide built-in charset conversion; return bytes directly and let callers convert if needed.
		_ = customCharset
	}
	return resp.Body()
}

// DownloadFile downloads to a file, using URL or response headers for the file name when dest is a directory.
func DownloadFile(rawURL, dest string) (int64, error) {
	resp := Get(rawURL).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.SaveAs(dest)
}

// Download downloads to a Writer.
func Download(rawURL string, w io.Writer) (int64, error) {
	resp := Get(rawURL).Execute()
	if resp.err != nil {
		return 0, resp.err
	}
	return resp.writeBodyTo(w)
}

// DownloadBytes downloads and returns bytes.
func DownloadBytes(rawURL string) []byte { return Get(rawURL).Execute().Bytes() }

// ToParams converts a map to a URL query string.
func ToParams(m map[string]any) string {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, toString(v))
	}
	return values.Encode()
}

// EncodeParams encodes a URL containing parameters; only the part after ? is encoded.
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

// DecodeParamMap parses a query string into a map.
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

// DecodeParams parses a query string into a multi-value map.
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

// URLWithForm appends form values to a URL.
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

// BuildBasicAuth builds a Basic Auth string.
func BuildBasicAuth(user, pass string) string {
	token := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	return "Basic " + token
}

var (
	// CharsetPattern matches charset in Content-Type.
	CharsetPattern = regexp.MustCompile(`(?i)charset\s*=\s*([a-z0-9-]+)`)
	// MetaCharsetPattern matches charset in HTML meta tags.
	MetaCharsetPattern = regexp.MustCompile(`(?i)<meta[^>]*?charset\s*=\s*['"]?([a-z0-9-]+)`)
)

// GetCharsetFromContentType extracts charset from Content-Type.
func GetCharsetFromContentType(ct string) string {
	m := CharsetPattern.FindStringSubmatch(ct)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// GetCharsetFromHTML extracts charset from HTML meta tags.
func GetCharsetFromHTML(html string) string {
	m := MetaCharsetPattern.FindStringSubmatch(html)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// GetMimeType returns the MIME type by file extension, or an empty string when unknown.
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

// Simplified MIME table.
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

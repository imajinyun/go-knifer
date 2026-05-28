package http

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// HTTPResponse 包装 http.Response 提供便捷读取（对应 hutool-http HttpResponse）。
type HTTPResponse struct {
	resp *http.Response
	body []byte
	once sync.Once
	err  error
}

func wrapResponse(r *http.Response) *HTTPResponse { return &HTTPResponse{resp: r} }

// Err 返回执行过程中的错误。
func (r *HTTPResponse) Err() error { return r.err }

// Status 返回 HTTP 状态码（出错时为 0）。
func (r *HTTPResponse) Status() int {
	if r.resp == nil {
		return 0
	}
	return r.resp.StatusCode
}

// IsOK 是否 2xx 成功。
func (r *HTTPResponse) IsOK() bool {
	return r.Status() >= 200 && r.Status() < 300
}

// Header 返回响应头。
func (r *HTTPResponse) Header(name string) string {
	if r.resp == nil {
		return ""
	}
	return r.resp.Header.Get(name)
}

// Headers 返回完整响应头。
func (r *HTTPResponse) Headers() http.Header {
	if r.resp == nil {
		return nil
	}
	return r.resp.Header
}

// Cookies 返回响应中的 Cookie。
func (r *HTTPResponse) Cookies() []*http.Cookie {
	if r.resp == nil {
		return nil
	}
	return r.resp.Cookies()
}

// ContentType 返回响应 Content-Type。
func (r *HTTPResponse) ContentType() string { return r.Header(string(HeaderContentType)) }

// ContentLength 返回响应 Content-Length。
func (r *HTTPResponse) ContentLength() int64 {
	if r.resp == nil {
		return -1
	}
	return r.resp.ContentLength
}

// Charset 解析 Content-Type 中的字符集，未指定时返回 UTF-8。
func (r *HTTPResponse) Charset() string {
	if cs := charsetFromContentType(r.ContentType()); cs != "" {
		return cs
	}
	return "UTF-8"
}

// Bytes 读取并返回响应体的字节内容。
func (r *HTTPResponse) Bytes() []byte {
	r.once.Do(func() {
		if r.resp == nil || r.resp.Body == nil {
			return
		}
		defer func() {
			if err := r.resp.Body.Close(); err != nil && r.err == nil {
				r.err = NewHTTPError("close response body failed", err)
			}
		}()
		reader, err := decodedBody(r.resp)
		if err != nil {
			r.err = err
			return
		}
		data, err := io.ReadAll(reader)
		if err != nil && (!IsIgnoreEOFError() || err != io.ErrUnexpectedEOF) {
			r.err = NewHTTPError("read response body failed", err)
			return
		}
		r.body = data
	})
	return r.body
}

// Body 读取响应体并以字符串返回。
func (r *HTTPResponse) Body() string { return string(r.Bytes()) }

// WriteTo 将响应体写入 Writer，返回写入字节数。
func (r *HTTPResponse) WriteTo(w io.Writer) (int64, error) {
	data := r.Bytes()
	if r.err != nil {
		return 0, r.err
	}
	n, err := w.Write(data)
	return int64(n), err
}

// SaveAs 将响应体保存到文件，返回写入字节数。
//
// 当 dest 是目录时，自动从 URL 或 Content-Disposition 中提取文件名。
func (r *HTTPResponse) SaveAs(dest string) (n int64, err error) {
	if r.resp == nil {
		return 0, HTTPErrorf("no response")
	}
	target := dest
	if info, err := os.Stat(dest); err == nil && info.IsDir() {
		fileName := r.fileName()
		if fileName == "" {
			fileName = "download.bin"
		}
		target = filepath.Join(dest, fileName)
	}
	f, err := os.Create(target)
	if err != nil {
		return 0, NewHTTPError("create file failed", err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return r.WriteTo(f)
}

// Close 关闭底层响应体（仅在未读取时需要）。
func (r *HTTPResponse) Close() error {
	if r.resp != nil && r.resp.Body != nil {
		return r.resp.Body.Close()
	}
	return nil
}

// Raw 返回原始 *http.Response（可用于流式处理，注意手动关闭 Body）。
func (r *HTTPResponse) Raw() *http.Response { return r.resp }

func (r *HTTPResponse) fileName() string {
	if cd := r.Header(string(HeaderContentDisposition)); cd != "" {
		if i := strings.Index(strings.ToLower(cd), "filename="); i >= 0 {
			name := strings.TrimSpace(cd[i+len("filename="):])
			name = strings.Trim(name, `"`)
			if idx := strings.Index(name, ";"); idx >= 0 {
				name = name[:idx]
			}
			if name != "" {
				return name
			}
		}
	}
	if r.resp != nil && r.resp.Request != nil && r.resp.Request.URL != nil {
		_, name := filepath.Split(r.resp.Request.URL.Path)
		return name
	}
	return ""
}

func decodedBody(resp *http.Response) (io.Reader, error) {
	enc := strings.ToLower(resp.Header.Get(string(HeaderContentEncoding)))
	switch enc {
	case "gzip":
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			// 部分服务即便声明 gzip 也可能未压缩，尝试回退
			if err == io.EOF {
				return bytes.NewReader(nil), nil
			}
			return nil, NewHTTPError("gzip reader init failed", err)
		}
		return gr, nil
	case "deflate":
		zr, err := zlib.NewReader(resp.Body)
		if err != nil {
			return nil, NewHTTPError("deflate reader init failed", err)
		}
		return zr, nil
	default:
		return resp.Body, nil
	}
}

var charsetRegex = regexp.MustCompile(`(?i)charset\s*=\s*([a-z0-9-]+)`)

// charsetFromContentType 从 Content-Type 中提取字符集。
func charsetFromContentType(ct string) string {
	if ct == "" {
		return ""
	}
	m := charsetRegex.FindStringSubmatch(ct)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

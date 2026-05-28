package http

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// HTTPRequest 链式 HTTP 请求构建器（对应 hutool-http HttpRequest）。
type HTTPRequest struct {
	method       Method
	rawURL       string
	queryParams  url.Values
	headers      http.Header
	cookies      []*http.Cookie
	body         []byte
	bodyReader   io.Reader
	form         map[string]any
	multipart    bool
	multipartFs  []*formFile
	contentType  string
	charset      string
	timeout      time.Duration
	followRedir  *bool
	maxRedirects int
	tlsSkip      bool
	transport    http.RoundTripper
	basicUser    string
	basicPass    string
	hasBasic     bool
	httpClient   *http.Client
}

type formFile struct {
	field    string
	fileName string
	data     []byte
	reader   io.Reader
}

// NewRequest 使用指定方法和 URL 创建请求。
func NewRequest(method Method, rawURL string) *HTTPRequest {
	return &HTTPRequest{
		method:       method,
		rawURL:       rawURL,
		queryParams:  url.Values{},
		headers:      CloneGlobalHeaders(),
		charset:      "UTF-8",
		maxRedirects: GetGlobalMaxRedirects(),
	}
}

// Get 创建 GET 请求。
func Get(rawURL string) *HTTPRequest { return NewRequest(MethodGet, rawURL) }

// Post 创建 POST 请求。
func Post(rawURL string) *HTTPRequest { return NewRequest(MethodPost, rawURL) }

// Put 创建 PUT 请求。
func Put(rawURL string) *HTTPRequest { return NewRequest(MethodPut, rawURL) }

// Delete 创建 DELETE 请求。
func Delete(rawURL string) *HTTPRequest { return NewRequest(MethodDelete, rawURL) }

// Patch 创建 PATCH 请求。
func Patch(rawURL string) *HTTPRequest { return NewRequest(MethodPatch, rawURL) }

// Head 创建 HEAD 请求。
func Head(rawURL string) *HTTPRequest { return NewRequest(MethodHead, rawURL) }

// Options 创建 OPTIONS 请求。
func Options(rawURL string) *HTTPRequest { return NewRequest(MethodOptions, rawURL) }

// Method 设置 HTTP 方法。
func (r *HTTPRequest) Method(m Method) *HTTPRequest { r.method = m; return r }

// URL 设置请求地址。
func (r *HTTPRequest) URL(u string) *HTTPRequest { r.rawURL = u; return r }

// Header 添加单个请求头（覆盖）。
func (r *HTTPRequest) Header(name, value string) *HTTPRequest {
	r.headers.Set(name, value)
	return r
}

// AddHeader 添加单个请求头（追加）。
func (r *HTTPRequest) AddHeader(name, value string) *HTTPRequest {
	r.headers.Add(name, value)
	return r
}

// Headers 批量设置请求头。
func (r *HTTPRequest) Headers(h map[string]string) *HTTPRequest {
	for k, v := range h {
		r.headers.Set(k, v)
	}
	return r
}

// Cookie 添加 Cookie。
func (r *HTTPRequest) Cookie(c *http.Cookie) *HTTPRequest {
	r.cookies = append(r.cookies, c)
	return r
}

// CookieString 通过原始字符串添加 Cookie。
func (r *HTTPRequest) CookieString(s string) *HTTPRequest {
	r.headers.Set(string(HeaderCookie), s)
	return r
}

// ContentType 设置 Content-Type。
func (r *HTTPRequest) ContentType(ct string) *HTTPRequest {
	r.contentType = ct
	return r
}

// Charset 设置请求字符集。
func (r *HTTPRequest) Charset(c string) *HTTPRequest { r.charset = c; return r }

// Timeout 设置请求超时。
func (r *HTTPRequest) Timeout(d time.Duration) *HTTPRequest { r.timeout = d; return r }

// FollowRedirects 设置是否跟随重定向。
func (r *HTTPRequest) FollowRedirects(b bool) *HTTPRequest {
	r.followRedir = &b
	return r
}

// MaxRedirects 设置最大重定向次数。
func (r *HTTPRequest) MaxRedirects(n int) *HTTPRequest { r.maxRedirects = n; return r }

// SkipTLSVerify 跳过 TLS 证书校验。
func (r *HTTPRequest) SkipTLSVerify(b bool) *HTTPRequest { r.tlsSkip = b; return r }

// Transport 自定义 RoundTripper。
func (r *HTTPRequest) Transport(t http.RoundTripper) *HTTPRequest { r.transport = t; return r }

// Client 自定义 *http.Client（设置后会覆盖 Transport / Timeout 等组合）。
func (r *HTTPRequest) Client(c *http.Client) *HTTPRequest { r.httpClient = c; return r }

// BasicAuth 设置 Basic Auth。
func (r *HTTPRequest) BasicAuth(user, pass string) *HTTPRequest {
	r.basicUser = user
	r.basicPass = pass
	r.hasBasic = true
	return r
}

// BearerAuth 设置 Bearer Token。
func (r *HTTPRequest) BearerAuth(token string) *HTTPRequest {
	r.headers.Set(string(HeaderAuthorization), "Bearer "+token)
	return r
}

// Query 添加单个 URL Query 参数。
func (r *HTTPRequest) Query(key string, value any) *HTTPRequest {
	r.queryParams.Add(key, toString(value))
	return r
}

// QueryMap 批量设置 URL Query 参数。
func (r *HTTPRequest) QueryMap(m map[string]any) *HTTPRequest {
	for k, v := range m {
		r.queryParams.Set(k, toString(v))
	}
	return r
}

// Body 设置原始请求体。
func (r *HTTPRequest) Body(body []byte) *HTTPRequest {
	r.body = body
	r.bodyReader = nil
	if r.contentType == "" {
		if ct := GuessContentType(string(body)); ct != "" {
			r.contentType = ct.WithCharset(r.charset)
		}
	}
	return r
}

// BodyString 设置字符串请求体。
func (r *HTTPRequest) BodyString(s string) *HTTPRequest { return r.Body([]byte(s)) }

// BodyJSON 设置 JSON 请求体（调用方需自行序列化或传入 string）。
func (r *HTTPRequest) BodyJSON(s string) *HTTPRequest {
	r.contentType = ContentTypeJSON.WithCharset(r.charset)
	return r.Body([]byte(s))
}

// BodyReader 通过 io.Reader 设置请求体。
func (r *HTTPRequest) BodyReader(reader io.Reader) *HTTPRequest {
	r.bodyReader = reader
	r.body = nil
	return r
}

// Form 设置表单参数（默认 form-urlencoded，存在文件时自动切换为 multipart）。
func (r *HTTPRequest) Form(m map[string]any) *HTTPRequest {
	if r.form == nil {
		r.form = make(map[string]any)
	}
	for k, v := range m {
		r.form[k] = v
	}
	return r
}

// FormFile 添加文件上传字段（自动启用 multipart）。
func (r *HTTPRequest) FormFile(field, fileName string, data []byte) *HTTPRequest {
	r.multipart = true
	r.multipartFs = append(r.multipartFs, &formFile{
		field: field, fileName: fileName, data: data,
	})
	return r
}

// FormFileReader 通过 Reader 添加文件上传字段。
func (r *HTTPRequest) FormFileReader(field, fileName string, reader io.Reader) *HTTPRequest {
	r.multipart = true
	r.multipartFs = append(r.multipartFs, &formFile{
		field: field, fileName: fileName, reader: reader,
	})
	return r
}

// Execute 执行请求并返回响应。
func (r *HTTPRequest) Execute() *HTTPResponse {
	resp, err := r.doExecute()
	if err != nil {
		return &HTTPResponse{err: err}
	}
	return resp
}

// MustExecute 执行请求，失败 panic。
func (r *HTTPRequest) MustExecute() *HTTPResponse {
	resp := r.Execute()
	if resp.err != nil {
		panic(resp.err)
	}
	return resp
}

func (r *HTTPRequest) buildURL() (string, error) {
	u, err := url.Parse(r.rawURL)
	if err != nil {
		return "", NewHTTPError("invalid url", err)
	}
	if len(r.queryParams) > 0 {
		q := u.Query()
		// 保持稳定输出顺序
		keys := make([]string, 0, len(r.queryParams))
		for k := range r.queryParams {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			for _, v := range r.queryParams[k] {
				q.Add(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

func (r *HTTPRequest) prepareBody() (io.Reader, string, error) {
	switch {
	case r.bodyReader != nil:
		return r.bodyReader, r.contentType, nil
	case len(r.body) > 0:
		return strings.NewReader(string(r.body)), r.contentType, nil
	case r.multipart || len(r.multipartFs) > 0:
		reader, ct, err := buildMultipartBody(r.form, r.multipartFs)
		if err != nil {
			return nil, "", err
		}
		return reader, ct, nil
	case len(r.form) > 0 && (r.method == MethodPost || r.method == MethodPut || r.method == MethodPatch):
		values := url.Values{}
		for k, v := range r.form {
			values.Set(k, toString(v))
		}
		ct := r.contentType
		if ct == "" {
			ct = ContentTypeFormURLEncoded.WithCharset(r.charset)
		}
		return strings.NewReader(values.Encode()), ct, nil
	case len(r.form) > 0:
		// GET 等：合并到 query
		for k, v := range r.form {
			r.queryParams.Add(k, toString(v))
		}
		r.form = nil
		return nil, r.contentType, nil
	}
	return nil, r.contentType, nil
}

func (r *HTTPRequest) buildClient() *http.Client {
	if r.httpClient != nil {
		return r.httpClient
	}
	transport := r.transport
	if transport == nil {
		t := &http.Transport{}
		if r.tlsSkip || IsTrustAnyHost() {
			t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		transport = t
	}
	timeout := r.timeout
	if timeout == 0 {
		timeout = GetGlobalTimeout()
	}
	follow := GetGlobalFollowRedirects()
	if r.followRedir != nil {
		follow = *r.followRedir
	}
	max := r.maxRedirects
	c := &http.Client{
		Timeout:   timeout,
		Transport: transport,
		Jar:       GetCookieJar(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !follow {
				return http.ErrUseLastResponse
			}
			if max > 0 && len(via) >= max {
				return fmt.Errorf("stopped after %d redirects", max)
			}
			return nil
		},
	}
	return c
}

func (r *HTTPRequest) doExecute() (*HTTPResponse, error) {
	finalURL, err := r.buildURL()
	if err != nil {
		return nil, err
	}
	bodyReader, ct, err := r.prepareBody()
	if err != nil {
		return nil, err
	}
	// prepareBody 可能修改 query，需要再构造一次
	if r.form != nil {
		finalURL, err = r.buildURL()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(string(r.method), finalURL, bodyReader)
	if err != nil {
		return nil, NewHTTPError("build request failed", err)
	}
	for k, vs := range r.headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if ct != "" {
		req.Header.Set(string(HeaderContentType), ct)
	}
	if ua := GetGlobalUserAgent(); ua != "" && req.Header.Get(string(HeaderUserAgent)) == "" {
		req.Header.Set(string(HeaderUserAgent), ua)
	}
	for _, c := range r.cookies {
		req.AddCookie(c)
	}
	if r.hasBasic {
		token := base64.StdEncoding.EncodeToString([]byte(r.basicUser + ":" + r.basicPass))
		req.Header.Set(string(HeaderAuthorization), "Basic "+token)
	}

	client := r.buildClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewHTTPError("send request failed", err)
	}
	return wrapResponse(resp), nil
}

func toString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

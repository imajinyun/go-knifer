package vhttp_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	knifer "github.com/imajinyun/go-knifer"
	"github.com/imajinyun/go-knifer/vhttp"
	"github.com/imajinyun/go-knifer/vurl"
)

func TestFacadeUsesNamesWithoutHTTPPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != string(vhttp.MethodGet) {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if got := r.URL.Query().Get("lang"); got != "go" {
			t.Fatalf("query lang = %q, want go", got)
		}
		if got := r.Header.Get("X-Client"); got != "go-knifer" {
			t.Fatalf("header X-Client = %q, want go-knifer", got)
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	req := vhttp.Get(server.URL).
		Query("lang", "go").
		Header("X-Client", "go-knifer")

	resp := executeRequest(req)
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if got := resp.Body(); got != "ok" {
		t.Fatalf("Body() = %q, want ok", got)
	}
}

func TestFacadeSharedConstants(t *testing.T) {
	_ = []vhttp.Method{vhttp.MethodTrace, vhttp.MethodConnect}
	_ = []vhttp.Header{vhttp.HeaderContentType, vhttp.HeaderUserAgent, vhttp.HeaderLocation}
	_ = []vhttp.ContentType{vhttp.ContentTypeJSON, vhttp.ContentTypeEventStream}

	if vhttp.MethodTrace.String() != http.MethodTrace {
		t.Fatalf("MethodTrace = %q", vhttp.MethodTrace.String())
	}
	if got := vhttp.ContentTypeJSON.WithCharset("UTF-8"); got != "application/json;charset=UTF-8" {
		t.Fatalf("ContentTypeJSON.WithCharset = %q", got)
	}
}

func TestFacadeRequestFollowRedirectOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Opt") + ":" + r.Header.Get("User-Agent")))
	}))
	defer server.Close()

	resp := vhttp.Get(
		server.URL,
		vhttp.WithHeader("X-Opt", "yes"),
		vhttp.WithUserAgent("vhttp-test/1.0"),
	).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if got := resp.Body(); got != "yes:vhttp-test/1.0" {
		t.Fatalf("Body() = %q, want option headers", got)
	}
}

func TestFacadeRequestCloneAndSingleUse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			b, _ := io.ReadAll(r.Body)
			_, _ = w.Write(b)
			return
		}
		_, _ = w.Write([]byte(r.URL.Query().Get("q") + ":" + r.Header.Get("X-Token")))
	}))
	defer server.Close()

	base := vhttp.Get(server.URL).Query("q", "base").Header("X-Token", "base")
	clone := base.Clone().Header("X-Token", "clone")
	if got := base.Execute().Body(); got != "base:base" {
		t.Fatalf("base Body() = %q", got)
	}
	if got := clone.Execute().Body(); got != "base:clone" {
		t.Fatalf("clone Body() = %q", got)
	}

	req := vhttp.Post(server.URL).BodyReader(strings.NewReader("payload"))
	if got := req.Execute().Body(); got != "payload" {
		t.Fatalf("first reader body = %q", got)
	}
	if resp := req.Execute(); resp.Err() == nil {
		t.Fatal("second Execute() should reject reader-backed body reuse")
	}
}

func TestFacadeTransportProviderOption(t *testing.T) {
	calls := 0
	transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(req.Header.Get("X-Transport"))),
			Header:     http.Header{},
			Request:    req,
		}, nil
	})
	resp := vhttp.Get("https://example.com",
		vhttp.WithHeader("X-Transport", "facade"),
		vhttp.WithTransportProvider(func() http.RoundTripper {
			calls++
			return transport
		}),
	).Execute()
	if resp.Err() != nil {
		t.Fatal(resp.Err())
	}
	if calls != 1 || resp.Body() != "facade" {
		t.Fatalf("transport provider calls=%d body=%q", calls, resp.Body())
	}
}

func TestFacadeDefaultTransportProviderLifecycle(t *testing.T) {
	custom := &http.Transport{MaxIdleConnsPerHost: 5}
	vhttp.ConfigureDefaultTransportProvider(func() *http.Transport { return custom })
	t.Cleanup(vhttp.ResetDefaultTransport)

	providerCalls := 0
	resp := vhttp.Get("https://example.com",
		vhttp.WithTransportProvider(func() http.RoundTripper {
			providerCalls++
			return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("ok")), Header: http.Header{}, Request: req}, nil
			})
		}),
	).Execute()
	if resp.Err() != nil || resp.Body() != "ok" || providerCalls != 1 {
		t.Fatalf("per-request transport provider resp=%q err=%v calls=%d", resp.Body(), resp.Err(), providerCalls)
	}

	vhttp.ResetDefaultTransport()
}

func TestFacadeRequestOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("X-Create")))
	}))
	defer server.Close()

	getResp := vhttp.Get(server.URL+"/redirect", vhttp.WithFollowRedirects(false), vhttp.WithHeader("X-Create", "get")).Execute()
	if getResp.Err() != nil {
		t.Fatal(getResp.Err())
	}
	if got := getResp.Status(); got != http.StatusFound {
		t.Fatalf("Get status = %d, want 302", got)
	}

	postResp := vhttp.Post(server.URL, vhttp.WithHeader("X-Create", "post")).Execute()
	if postResp.Err() != nil {
		t.Fatal(postResp.Err())
	}
	if got := postResp.Body(); got != "POST:post" {
		t.Fatalf("Post body = %q, want POST:post", got)
	}
}

func TestFacadeAdditionalHTTPMethodsAndGlobalAccessors(t *testing.T) {
	var lastMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastMethod = r.Method
		w.Header().Set("X-Method", r.Method)
		if r.Method != http.MethodHead {
			_, _ = w.Write([]byte(r.Method))
		}
	}))
	defer server.Close()

	tests := []struct {
		name   string
		method string
		req    *vhttp.Request
	}{
		{name: "put", method: http.MethodPut, req: vhttp.Put(server.URL)},
		{name: "delete", method: http.MethodDelete, req: vhttp.Delete(server.URL)},
		{name: "patch", method: http.MethodPatch, req: vhttp.Patch(server.URL)},
		{name: "head", method: http.MethodHead, req: vhttp.Head(server.URL)},
		{name: "options", method: http.MethodOptions, req: vhttp.Options(server.URL)},
		{name: "new request", method: http.MethodTrace, req: vhttp.NewRequest(vhttp.MethodTrace, server.URL)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.req.Execute()
			if resp.Err() != nil {
				t.Fatalf("Execute: %v", resp.Err())
			}
			if lastMethod != tt.method {
				t.Fatalf("server method = %q, want %q", lastMethod, tt.method)
			}
			if got := resp.Header("X-Method"); got != tt.method {
				t.Fatalf("response method header = %q, want %q", got, tt.method)
			}
		})
	}

	previous := vhttp.SnapshotGlobalConfig()
	defer vhttp.ConfigureGlobalConfig(previous)
	vhttp.SetGlobalMaxRedirects(3)
	vhttp.SetGlobalMaxResponseBytes(99)
	vhttp.SetGlobalFollowRedirects(false)
	vhttp.SetGlobalUserAgent("vhttp-extra/1.0")
	vhttp.SetIgnoreEOFError(true)
	vhttp.SetGlobalBoundary("boundary-extra")
	vhttp.SetGlobalDecodeURL(true)
	if vhttp.GetGlobalMaxRedirects() != 3 || vhttp.GetGlobalMaxResponseBytes() != 99 || vhttp.GetGlobalFollowRedirects() || vhttp.GetGlobalUserAgent() != "vhttp-extra/1.0" || !vhttp.IsIgnoreEOFError() || vhttp.GetGlobalBoundary() != "boundary-extra" || !vhttp.IsGlobalDecodeURL() {
		t.Fatalf("global accessors snapshot = %#v", vhttp.SnapshotGlobalConfig())
	}
}

func TestFacadeContentHTMLAndUserAgentHelpers(t *testing.T) {
	if got := vhttp.BuildContentType("text/plain", "utf-8"); got != "text/plain;charset=utf-8" {
		t.Fatalf("BuildContentType = %q", got)
	}
	if !vhttp.IsDefaultContentType("") || !vhttp.IsFormURLEncoded("application/x-www-form-urlencoded; charset=utf-8") {
		t.Fatal("content type predicates returned unexpected result")
	}
	if got := vhttp.GuessContentType(`{"ok":true}`); got != vhttp.ContentTypeJSON {
		t.Fatalf("GuessContentType json = %q", got)
	}
	if got := vhttp.GetCharsetFromContentType("text/plain; charset=gbk"); got != "gbk" {
		t.Fatalf("GetCharsetFromContentType = %q", got)
	}
	if got := vhttp.GetCharsetFromContentTypeWithOptions("text/plain; enc=big5", vhttp.WithCharsetRegexp(regexp.MustCompile(`enc=([a-z0-9-]+)`))); got != "big5" {
		t.Fatalf("GetCharsetFromContentTypeWithOptions = %q", got)
	}
	if got := vhttp.GetCharsetFromHTML(`<meta charset="utf-8"><p>x</p>`); got != "utf-8" {
		t.Fatalf("GetCharsetFromHTML = %q", got)
	}
	if got := vhttp.GetCharsetFromHTMLWithOptions(`<meta data-charset="gb2312">`, vhttp.WithMetaCharsetRegexp(regexp.MustCompile(`data-charset="([^"]+)"`))); got != "gb2312" {
		t.Fatalf("GetCharsetFromHTMLWithOptions = %q", got)
	}
	if got := vhttp.GetMimeType("payload.JSON"); got != "application/json" {
		t.Fatalf("GetMimeType = %q", got)
	}

	if got := vhttp.HTMLEscape(`<b>"go"</b>`); got != `&lt;b&gt;&#34;go&#34;&lt;/b&gt;` {
		t.Fatalf("HTMLEscape = %q", got)
	}
	if got := vhttp.HTMLUnescape("&lt;b&gt;go&lt;/b&gt;"); got != "<b>go</b>" {
		t.Fatalf("HTMLUnescape = %q", got)
	}
	if got := vhttp.CleanHTML("<p>Hello</p><!--drop-->"); got != "Hello" {
		t.Fatalf("CleanHTML = %q", got)
	}
	if got := vhttp.CleanHTMLWithOptions("a[drop]b", vhttp.WithHTMLTagRegexp(regexp.MustCompile(`\[.*?\]`)), vhttp.WithHTMLCommentRegexp(regexp.MustCompile(`$^`))); got != "ab" {
		t.Fatalf("CleanHTMLWithOptions = %q", got)
	}
	if got := vhttp.FilterHTMLTag("<div>ok</div><script>x</script>", "script"); got != "<div>ok</div>" {
		t.Fatalf("FilterHTMLTag = %q", got)
	}
	if got := vhttp.FilterHTMLTagWithOptions("<custom>drop</custom><p>keep</p>", []string{"custom"}, vhttp.WithHTMLFilterCompileFunc(regexp.Compile)); got != "<p>keep</p>" {
		t.Fatalf("FilterHTMLTagWithOptions = %q", got)
	}
	if ua := vhttp.ParseUserAgent("Mozilla/5.0 Chrome/120.0"); ua == nil {
		t.Fatal("ParseUserAgent returned nil")
	}
	if !vhttp.IsRedirected(http.StatusTemporaryRedirect) || vhttp.IsRedirected(http.StatusOK) {
		t.Fatal("IsRedirected returned unexpected result")
	}
}

func TestFacadeResponseDecodeOptions(t *testing.T) {
	gzipServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		_, _ = gz.Write([]byte("gzipped"))
		_ = gz.Close()
	}))
	defer gzipServer.Close()

	compressed := vhttp.Get(gzipServer.URL, vhttp.WithAutoDecodeResponse(false)).Execute().Bytes()
	if bytes.Contains(compressed, []byte("gzipped")) || len(compressed) == 0 {
		t.Fatalf("body should remain compressed, got %q", compressed)
	}

	customServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "upper")
		_, _ = w.Write([]byte("hello"))
	}))
	defer customServer.Close()

	decoder := func(r io.Reader) (io.ReadCloser, error) {
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(strings.NewReader(strings.ToUpper(string(data)))), nil
	}
	if got := vhttp.Get(customServer.URL, vhttp.WithContentDecoder("upper", decoder)).Execute().Body(); got != "HELLO" {
		t.Fatalf("custom decoded body = %q", got)
	}
}

func TestFacadeSimpleServerOptions(t *testing.T) {
	server := vhttp.NewSimpleServerAddrWithOptions("127.0.0.1:0",
		vhttp.WithReadHeaderTimeout(time.Second),
		vhttp.WithReadTimeout(time.Second),
		vhttp.WithWriteTimeout(time.Second),
		vhttp.WithIdleTimeout(time.Second),
		vhttp.WithHTTPServer(&http.Server{Addr: "127.0.0.1:0"}),
	)
	if server == nil {
		t.Fatal("NewSimpleServerAddrWithOptions returned nil")
	}
	if err := server.StopWithContext(context.Background()); err != nil {
		t.Fatalf("StopWithContext on idle server = %v", err)
	}
}

func TestFacadeServerStarterLifecycle(t *testing.T) {
	vhttp.ResetServerStarters()
	t.Cleanup(vhttp.ResetServerStarters)

	called := 0
	server := vhttp.NewSimpleServerAddrWithOptions("127.0.0.1:0", vhttp.WithListenAndServeFunc(func(server *http.Server) error {
		called++
		return http.ErrServerClosed
	}))
	if err := server.Start(); err != http.ErrServerClosed {
		t.Fatalf("Start() = %v, want ErrServerClosed", err)
	}
	if called != 1 {
		t.Fatalf("custom starter called %d times, want 1", called)
	}
}

func TestFacadeHelperNamesWithoutHTTPPrefix(t *testing.T) {
	previous := vhttp.SnapshotGlobalConfig()
	defer vhttp.ConfigureGlobalConfig(previous)

	vhttp.SetGlobalTimeout(2 * time.Second)
	if got := vhttp.GetGlobalTimeout(); got != 2*time.Second {
		t.Fatalf("GetGlobalTimeout() = %v, want 2s", got)
	}

	vhttp.SetGlobalHeader("X-Test", "a")
	vhttp.AddGlobalHeader("X-Test", "b")
	if got := vhttp.CloneGlobalHeaders().Values("X-Test"); len(got) != 2 {
		t.Fatalf("CloneGlobalHeaders().Values(X-Test) = %v, want 2 values", got)
	}
	vhttp.RemoveGlobalHeader("X-Test")
	if got := vhttp.CloneGlobalHeaders().Values("X-Test"); len(got) != 0 {
		t.Fatalf("after RemoveGlobalHeader values = %v, want empty", got)
	}

	if got := vhttp.BuildBasicAuth("aladdin", "opensesame"); got != "Basic YWxhZGRpbjpvcGVuc2VzYW1l" {
		t.Fatalf("BuildBasicAuth() = %q", got)
	}
	if got := vurl.EncodeQueryMap(map[string]any{"q": "go", "page": 1}); !strings.Contains(got, "q=go") || !strings.Contains(got, "page=1") {
		t.Fatalf("EncodeQueryMap() = %q", got)
	}
}

func TestFacadeScopedGlobalConfig(t *testing.T) {
	previous := vhttp.SnapshotGlobalConfig()
	defer vhttp.ConfigureGlobalConfig(previous)

	vhttp.ResetGlobalConfig()
	vhttp.WithScopedGlobalConfig(vhttp.GlobalConfig{
		Timeout:          3 * time.Second,
		MaxRedirects:     1,
		MaxResponseBytes: 32,
		IgnoreEOFError:   true,
		FollowRedirects:  false,
		DefaultUserAgent: "facade-scope-agent",
		Boundary:         "facade-boundary",
		Headers:          http.Header{"X-Facade-Scope": []string{"inner"}},
		CookieJar:        nil,
	}, func() {
		cfg := vhttp.SnapshotGlobalConfig()
		if cfg.Timeout != 3*time.Second || cfg.MaxRedirects != 1 || cfg.MaxResponseBytes != 32 || cfg.FollowRedirects || cfg.DefaultUserAgent != "facade-scope-agent" || cfg.Headers.Get("X-Facade-Scope") != "inner" || cfg.CookieJar != nil {
			t.Fatalf("facade scoped config = %#v", cfg)
		}
	})

	cfg := vhttp.SnapshotGlobalConfig()
	if cfg.Timeout != 30*time.Second || cfg.MaxRedirects != 10 || cfg.MaxResponseBytes != 64<<20 || !cfg.FollowRedirects || cfg.Headers.Get("X-Facade-Scope") != "" || cfg.CookieJar == nil {
		t.Fatalf("facade config not restored after scoped helper: %#v", cfg)
	}
}

func TestFacadeErrorNamesWithoutHTTPPrefix(t *testing.T) {
	cause := errors.New("closed")
	err := vhttp.NewError("read failed", cause)
	if !errors.Is(err, cause) {
		t.Fatalf("NewError() does not unwrap cause")
	}
	if !errors.Is(err, knifer.ErrCodeInternal) {
		t.Fatalf("NewError() does not match ErrCodeInternal")
	}
	code, ok := knifer.CodeOf(err)
	if !ok || code != knifer.ErrCodeInternal {
		t.Fatalf("CodeOf(NewError()) = %q, %v; want internal", code, ok)
	}

	formatted := vhttp.Errorf("status %d", 500)
	if got := errorString(formatted); got != "status 500" {
		t.Fatalf("Errorf().Error() = %q, want status 500", got)
	}
}

func TestFacadeSaveProviderOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("vhttp-save"))
	}))
	defer server.Close()

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written bytes.Buffer
	n, err := vhttp.Get(server.URL).Execute().SaveAs("/virtual/out.txt",
		vhttp.WithSaveMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		vhttp.WithSaveOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return nopWriteCloser{Writer: &written}, nil
		}),
		vhttp.WithSaveDirPerm(0o700), vhttp.WithSaveFilePerm(0o600),
	)
	if err != nil || n != int64(len("vhttp-save")) {
		t.Fatalf("SaveAs provider n=%d err=%v", n, err)
	}
	if mkdirPath != "/virtual" || mkdirPerm != 0o700 || openPath != "/virtual/out.txt" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.String() != "vhttp-save" {
		t.Fatalf("providers mkdir=%q/%v open=%q flag=%#x perm=%v content=%q", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.String())
	}
}

func TestFacadeDownloadFileSafe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Mode") != "safe" {
			http.Error(w, "missing option header", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("vhttp-safe-file"))
	}))
	defer server.Close()

	dir := t.TempDir()
	n, err := vhttp.DownloadFileSafeWithOptions(server.URL, dir,
		[]vhttp.RequestOption{
			vhttp.WithHeader("X-Mode", "safe"),
			vhttp.WithURLPolicy(vhttp.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false}),
		},
		vhttp.WithSaveDefaultFilename("safe.txt"),
	)
	if err != nil {
		t.Fatalf("DownloadFileSafeWithOptions() error = %v", err)
	}
	if n != int64(len("vhttp-safe-file")) {
		t.Fatalf("DownloadFileSafeWithOptions() n = %d", n)
	}
	data, err := os.ReadFile(filepath.Join(dir, "safe.txt"))
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}
	if string(data) != "vhttp-safe-file" {
		t.Fatalf("saved file = %q", data)
	}
}

func TestFacadeShortcutRequestHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch r.URL.Path {
		case "/get":
			_, _ = w.Write([]byte(r.Method + ":" + r.URL.Query().Get("q")))
		case "/form":
			_, _ = w.Write([]byte(r.Method + ":" + string(body) + ":" + r.Header.Get("X-Shortcut")))
		case "/json":
			_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("Content-Type") + ":" + string(body)))
		case "/string":
			_, _ = w.Write([]byte(r.Method + ":" + string(body)))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	got, err := vhttp.GetStringE(server.URL + "/get")
	if err != nil || got != "GET:" {
		t.Fatalf("GetStringE = %q, %v", got, err)
	}
	got, err = vhttp.GetStringEWithOptions(server.URL+"/get?q=go", vhttp.WithHeader("X-Shortcut", "yes"))
	if err != nil || got != "GET:go" {
		t.Fatalf("GetStringEWithOptions = %q, %v", got, err)
	}
	got, err = vhttp.GetWithTimeoutE(server.URL+"/get", time.Second)
	if err != nil || got != "GET:" {
		t.Fatalf("GetWithTimeoutE = %q, %v", got, err)
	}
	got, err = vhttp.GetWithTimeoutEWithOptions(server.URL+"/get?q=timeout", time.Second, vhttp.WithHeader("X-Shortcut", "yes"))
	if err != nil || got != "GET:timeout" {
		t.Fatalf("GetWithTimeoutEWithOptions = %q, %v", got, err)
	}
	got, err = vhttp.GetWithParamsE(server.URL+"/get", map[string]any{"q": "params"})
	if err != nil || got != "GET:params" {
		t.Fatalf("GetWithParamsE = %q, %v", got, err)
	}
	got, err = vhttp.GetWithParamsEWithOptions(server.URL+"/get", map[string]any{"q": "params2"}, vhttp.WithHeader("X-Shortcut", "yes"))
	if err != nil || got != "GET:params2" {
		t.Fatalf("GetWithParamsEWithOptions = %q, %v", got, err)
	}
	got, err = vhttp.PostFormE(server.URL+"/form", map[string]any{"name": "alice"})
	if err != nil || !strings.Contains(got, "POST:") || !strings.Contains(got, "name=alice") {
		t.Fatalf("PostFormE = %q, %v", got, err)
	}
	got, err = vhttp.PostFormEWithOptions(server.URL+"/form", map[string]any{"name": "bob"}, vhttp.WithHeader("X-Shortcut", "hdr"))
	if err != nil || !strings.Contains(got, "name=bob") || !strings.HasSuffix(got, ":hdr") {
		t.Fatalf("PostFormEWithOptions = %q, %v", got, err)
	}
	got, err = vhttp.PostJSONE(server.URL+"/json", `{"name":"json"}`)
	if err != nil || !strings.Contains(got, `{"name":"json"}`) || !strings.Contains(got, "application/json") {
		t.Fatalf("PostJSONE = %q, %v", got, err)
	}
	got, err = vhttp.PostJSONEWithOptions(server.URL+"/json", `{"name":"json2"}`, vhttp.WithHeader("X-Shortcut", "hdr"))
	if err != nil || !strings.Contains(got, `{"name":"json2"}`) {
		t.Fatalf("PostJSONEWithOptions = %q, %v", got, err)
	}
	got, err = vhttp.PostStringE(server.URL+"/string", "plain")
	if err != nil || got != "POST:plain" {
		t.Fatalf("PostStringE = %q, %v", got, err)
	}
	got, err = vhttp.PostStringEWithOptions(server.URL+"/string", "plain2", vhttp.WithHeader("X-Shortcut", "hdr"))
	if err != nil || got != "POST:plain2" {
		t.Fatalf("PostStringEWithOptions = %q, %v", got, err)
	}
}

func TestFacadeSafeShortcutHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if r.Method == http.MethodPost {
			_, _ = w.Write([]byte(r.Method + ":" + string(body)))
			return
		}
		_, _ = w.Write([]byte(r.Method))
	}))
	defer server.Close()

	allowLocal := vhttp.WithURLPolicy(vhttp.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	if got, err := vhttp.GetStringSafeE(server.URL, allowLocal); err != nil || got != "GET" {
		t.Fatalf("GetStringSafeE allowed = %q, %v", got, err)
	}
	if _, err := vhttp.GetStringSafeE(server.URL); err == nil {
		t.Fatal("GetStringSafeE(localhost default policy) error = nil")
	}
	if got, err := vhttp.PostFormSafeE(server.URL, map[string]any{"name": "safe"}, allowLocal); err != nil || got != "POST:name=safe" {
		t.Fatalf("PostFormSafeE = %q, %v", got, err)
	}
	if got, err := vhttp.PostJSONSafeE(server.URL, `{"safe":true}`, allowLocal); err != nil || got != `POST:{"safe":true}` {
		t.Fatalf("PostJSONSafeE = %q, %v", got, err)
	}
	if got, err := vhttp.PostStringSafeE(server.URL, "safe-string", allowLocal); err != nil || got != "POST:safe-string" {
		t.Fatalf("PostStringSafeE = %q, %v", got, err)
	}
}

func TestFacadeClientAndAdditionalServerOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method + ":" + r.Header.Get("X-Client-Default")))
	}))
	defer server.Close()

	client := vhttp.NewClient(vhttp.WithClientRequestOptions(vhttp.WithHeader("X-Client-Default", "shared")))
	if got := client.Get(server.URL).Execute().Body(); got != "GET:shared" {
		t.Fatalf("client.Get body = %q", got)
	}
	if got := client.Post(server.URL).Execute().Body(); got != "POST:shared" {
		t.Fatalf("client.Post body = %q", got)
	}
	if got := client.NewRequest(vhttp.MethodPut, server.URL).Execute().Body(); got != "PUT:shared" {
		t.Fatalf("client.NewRequest body = %q", got)
	}

	cfg := vhttp.SnapshotGlobalConfig()
	cfg.Headers.Set("X-Client-Default", "configured")
	if got := vhttp.NewClientWithConfig(cfg).Get(server.URL).Execute().Body(); got != "GET:configured" {
		t.Fatalf("NewClientWithConfig body = %q", got)
	}
	isolated := vhttp.NewIsolatedClient(vhttp.WithClientGlobalConfig(cfg))
	if got := isolated.Get(server.URL).Execute().Body(); got != "GET:configured" {
		t.Fatalf("NewIsolatedClient body = %q", got)
	}
	if resp := client.GetSafe(server.URL).Execute(); resp.Err() == nil {
		t.Fatal("client.GetSafe(localhost default policy) error = nil")
	}
	if resp := client.PostSafe(server.URL, vhttp.WithURLPolicy(vhttp.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})).Execute(); resp.Err() != nil {
		t.Fatalf("client.PostSafe allowed error = %v", resp.Err())
	}

	baseContextCalled := false
	connContextCalled := false
	runnerCalled := false
	logger := log.New(io.Discard, "", 0)
	simple := vhttp.NewSimpleServerWithOptions(0,
		vhttp.WithServerErrorLog(logger),
		vhttp.WithBaseContext(func(net.Listener) context.Context {
			baseContextCalled = true
			return context.Background()
		}),
		vhttp.WithConnContext(func(ctx context.Context, conn net.Conn) context.Context {
			connContextCalled = conn != nil
			return ctx
		}),
		vhttp.WithAsyncRunner(func(run func()) {
			runnerCalled = true
			run()
		}),
		vhttp.WithListenAndServeFunc(func(server *http.Server) error {
			if server.BaseContext != nil {
				_ = server.BaseContext(nil)
			}
			if server.ConnContext != nil {
				_ = server.ConnContext(context.Background(), nil)
			}
			return http.ErrServerClosed
		}),
	)
	errCh := simple.StartAsync()
	if err, ok := <-errCh; ok || err != nil {
		t.Fatalf("StartAsync channel = (%v, %v), want closed", err, ok)
	}
	if !runnerCalled || !baseContextCalled {
		t.Fatalf("server option calls runner=%v base=%v conn=%v", runnerCalled, baseContextCalled, connContextCalled)
	}
	if created := vhttp.CreateServer(0); created == nil {
		t.Fatal("CreateServer returned nil")
	}
	if created := vhttp.CreateServerWithOptions(0, vhttp.WithHTTPServer(&http.Server{Addr: ":0"})); created == nil {
		t.Fatal("CreateServerWithOptions returned nil")
	}
}

func TestFacadeSafeRequestsAndDownloadWrappers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Method", r.Method)
		if r.Method != http.MethodHead {
			_, _ = w.Write([]byte("download-text"))
		}
	}))
	defer server.Close()

	allowLocal := vhttp.WithURLPolicy(vhttp.URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})
	tests := []struct {
		name   string
		method string
		req    *vhttp.Request
	}{
		{name: "get safe", method: http.MethodGet, req: vhttp.GetSafe(server.URL, allowLocal)},
		{name: "post safe", method: http.MethodPost, req: vhttp.PostSafe(server.URL, allowLocal)},
		{name: "put safe", method: http.MethodPut, req: vhttp.PutSafe(server.URL, allowLocal)},
		{name: "delete safe", method: http.MethodDelete, req: vhttp.DeleteSafe(server.URL, allowLocal)},
		{name: "patch safe", method: http.MethodPatch, req: vhttp.PatchSafe(server.URL, allowLocal)},
		{name: "head safe", method: http.MethodHead, req: vhttp.HeadSafe(server.URL, allowLocal)},
		{name: "options safe", method: http.MethodOptions, req: vhttp.OptionsSafe(server.URL, allowLocal)},
		{name: "new safe", method: http.MethodTrace, req: vhttp.NewSafeRequest(vhttp.MethodTrace, server.URL, allowLocal)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.req.Execute()
			if resp.Err() != nil {
				t.Fatalf("Execute: %v", resp.Err())
			}
			if got := resp.Header("X-Method"); got != tt.method {
				t.Fatalf("method header = %q, want %q", got, tt.method)
			}
		})
	}

	var buf bytes.Buffer
	if n, err := vhttp.Download(server.URL, &buf); err != nil || n != int64(len("download-text")) || buf.String() != "download-text" {
		t.Fatalf("Download n=%d body=%q err=%v", n, buf.String(), err)
	}
	buf.Reset()
	if n, err := vhttp.DownloadWithOptions(server.URL, &buf, vhttp.WithMaxResponseBytes(64)); err != nil || n != int64(len("download-text")) || buf.String() != "download-text" {
		t.Fatalf("DownloadWithOptions n=%d body=%q err=%v", n, buf.String(), err)
	}
	buf.Reset()
	if n, err := vhttp.DownloadSafe(server.URL, &buf, allowLocal); err != nil || n != int64(len("download-text")) || buf.String() != "download-text" {
		t.Fatalf("DownloadSafe n=%d body=%q err=%v", n, buf.String(), err)
	}
	if b, err := vhttp.DownloadBytesE(server.URL); err != nil || string(b) != "download-text" {
		t.Fatalf("DownloadBytesE = %q, %v", b, err)
	}
	if b, err := vhttp.DownloadBytesEWithOptions(server.URL, vhttp.WithMaxResponseBytes(64)); err != nil || string(b) != "download-text" {
		t.Fatalf("DownloadBytesEWithOptions = %q, %v", b, err)
	}
	if b, err := vhttp.DownloadBytesSafeE(server.URL, allowLocal); err != nil || string(b) != "download-text" {
		t.Fatalf("DownloadBytesSafeE = %q, %v", b, err)
	}
	if got, err := vhttp.DownloadStringE(server.URL, ""); err != nil || got != "download-text" {
		t.Fatalf("DownloadStringE = %q, %v", got, err)
	}
	if got, err := vhttp.DownloadStringEWithOptions(server.URL, "", vhttp.WithMaxResponseBytes(64)); err != nil || got != "download-text" {
		t.Fatalf("DownloadStringEWithOptions = %q, %v", got, err)
	}
	if got, err := vhttp.DownloadStringSafeE(server.URL, "", allowLocal); err != nil || got != "download-text" {
		t.Fatalf("DownloadStringSafeE = %q, %v", got, err)
	}

	dir := t.TempDir()
	file := filepath.Join(dir, "plain.txt")
	if n, err := vhttp.DownloadFile(server.URL, file); err != nil || n != int64(len("download-text")) {
		t.Fatalf("DownloadFile n=%d err=%v", n, err)
	}
	fileWithOpts := filepath.Join(dir, "with-options.txt")
	if n, err := vhttp.DownloadFileWithOptions(server.URL, fileWithOpts, []vhttp.RequestOption{vhttp.WithMaxResponseBytes(64)}, vhttp.WithSaveOverwrite(true)); err != nil || n != int64(len("download-text")) {
		t.Fatalf("DownloadFileWithOptions n=%d err=%v", n, err)
	}
	if n, err := vhttp.DownloadFileSafe(server.URL, filepath.Join(dir, "blocked.txt")); err == nil || n != 0 {
		t.Fatalf("DownloadFileSafe default policy n=%d err=%v, want private host rejection", n, err)
	}
}

func TestFacadeAdditionalOptionWrappers(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	previous := vhttp.SnapshotGlobalConfig()
	defer vhttp.ConfigureGlobalConfig(previous)
	vhttp.SetCookieJar(jar)
	if vhttp.GetCookieJar() == nil {
		t.Fatal("GetCookieJar returned nil after SetCookieJar")
	}
	vhttp.CloseCookie()
	if vhttp.GetCookieJar() != nil {
		t.Fatal("CloseCookie did not clear global jar")
	}

	requestFactoryCalled := false
	readAllCalled := false
	transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
			Body:       io.NopCloser(strings.NewReader(req.Method + ":" + req.Header.Get("X-A"))),
			Request:    req,
		}, nil
	})
	client := &http.Client{Transport: transport}
	cfg := vhttp.SnapshotGlobalConfig()
	cfg.Headers.Set("X-A", "from-config")
	resp := vhttp.NewIsolatedRequest(vhttp.MethodPost, "https://example.com",
		vhttp.WithGlobalConfig(cfg),
		vhttp.WithTimeout(time.Second),
		vhttp.WithHeaders(map[string]string{"X-A": "one"}),
		vhttp.WithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}),
		vhttp.WithTransport(transport),
		vhttp.WithClient(client),
		vhttp.WithCookieJar(jar),
		vhttp.WithContentType(string(vhttp.ContentTypeTextPlain)),
		vhttp.WithCharset("utf-8"),
		vhttp.WithMaxResponseBytes(64),
		vhttp.WithResponseReadAllFunc(func(r io.Reader) ([]byte, error) {
			readAllCalled = true
			return io.ReadAll(r)
		}),
		vhttp.WithRequestFactory(func(method, rawURL string, body io.Reader) (*http.Request, error) {
			requestFactoryCalled = true
			return http.NewRequest(method, rawURL, body)
		}),
		vhttp.WithMultipartWriterFactory(func(w io.Writer) vhttp.MultipartWriter {
			return multipart.NewWriter(w)
		}),
	).BodyString("payload").Execute()
	if resp.Err() != nil {
		t.Fatalf("NewIsolatedRequest Execute: %v", resp.Err())
	}
	if got := resp.Body(); got != "POST:one" {
		t.Fatalf("option response = %q", got)
	}
	if !requestFactoryCalled || !readAllCalled {
		t.Fatalf("providers called request=%v readAll=%v", requestFactoryCalled, readAllCalled)
	}

	cfg.Headers.Set("X-A", "configured")
	if got := vhttp.NewRequestWithConfig(vhttp.MethodGet, "https://example.com", cfg, vhttp.WithTransport(transport)).Execute().Body(); got != "GET:configured" {
		t.Fatalf("NewRequestWithConfig body = %q", got)
	}
	if resp := vhttp.Get("http://public.example",
		vhttp.WithAllowedHosts("public.example"),
		vhttp.WithLookupIP(func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("8.8.8.8")}, nil
		}),
		vhttp.WithTransport(transport),
	).Execute(); resp.Err() != nil {
		t.Fatalf("allowed host request error = %v", resp.Err())
	}
}

func TestFacadeErrorCodesAndSimpleServerWrappers(t *testing.T) {
	cause := errors.New("bad request")
	err := vhttp.NewErrorWithCode(knifer.ErrCodeInvalidInput, "invalid request", cause)
	if !errors.Is(err, cause) || !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("NewErrorWithCode does not unwrap cause or code: %v", err)
	}
	if got := vhttp.ErrorfWithCode(knifer.ErrCodeInvalidInput, "status %d", http.StatusBadRequest).Error(); got != "status 400" {
		t.Fatalf("ErrorfWithCode = %q", got)
	}
	if vhttp.NewSimpleServer(0) == nil {
		t.Fatal("NewSimpleServer returned nil")
	}
	if vhttp.NewSimpleServerAddr("127.0.0.1:0") == nil {
		t.Fatal("NewSimpleServerAddr returned nil")
	}

	static := vhttp.NewSimpleServerAddrWithOptions("127.0.0.1:0")
	static.SetRootWithOptions(".",
		vhttp.WithStaticFileSystem(http.Dir(".")),
		vhttp.WithStaticFS(os.DirFS(".")),
		vhttp.WithFileServerFactory(func(http.FileSystem) http.Handler {
			return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			})
		}),
		vhttp.WithStaticHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		})),
	)

	listener, listenErr := net.Listen("tcp", "127.0.0.1:0")
	if listenErr != nil {
		t.Fatal(listenErr)
	}
	defer listener.Close()
	if server := vhttp.NewSimpleServerAddrWithOptions("127.0.0.1:0",
		vhttp.WithListener(listener),
		vhttp.WithListenAndServeFunc(func(*http.Server) error {
			return http.ErrServerClosed
		}),
	); server == nil {
		t.Fatal("WithListener server returned nil")
	}
}

type nopWriteCloser struct{ io.Writer }

func (w nopWriteCloser) Close() error { return nil }

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func executeRequest(req *vhttp.Request) *vhttp.Response {
	return req.Execute()
}

func errorString(err *vhttp.Error) string {
	return err.Error()
}

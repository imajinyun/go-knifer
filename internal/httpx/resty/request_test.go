package resty

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	grestry "resty.dev/v3"
)

func TestGetWithQueryAndHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "go" {
			t.Fatalf("query q = %q, want go", r.URL.Query().Get("q"))
		}
		if r.Header.Get("X-Test") != "yes" {
			t.Fatalf("X-Test = %q, want yes", r.Header.Get("X-Test"))
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Query("q", "go").Header("X-Test", "yes").Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if !resp.IsOK() || resp.Body() != "ok" {
		t.Fatalf("status/body = %d/%q, want 2xx/ok", resp.Status(), resp.Body())
	}
}

func TestPostForm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm() error = %v", err)
		}
		if got := r.Form.Get("name"); got != "resty" {
			t.Fatalf("form name = %q, want resty", got)
		}
		_, _ = w.Write([]byte("posted"))
	}))
	defer srv.Close()

	if got := PostForm(srv.URL, map[string]any{"name": "resty"}); got != "posted" {
		t.Fatalf("PostForm() = %q, want posted", got)
	}
}

func TestPostJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), string(ContentTypeJSON)) {
			t.Fatalf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		_, _ = w.Write([]byte("json"))
	}))
	defer srv.Close()

	if got := PostJSON(srv.URL, `{"ok":true}`); got != "json" {
		t.Fatalf("PostJSON() = %q, want json", got)
	}
}

func TestResponseHeadersCookiesAndLength(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Cookie"); !strings.Contains(got, "k=v") {
			t.Fatalf("Cookie = %q, want k=v", got)
		}
		w.Header().Set("X-Test", "yes")
		w.Header().Add("Set-Cookie", "sid=abc; Path=/")
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Cookie("k", "v").Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if got := resp.Headers()["X-Test"]; len(got) != 1 || got[0] != "yes" {
		t.Fatalf("Headers()[X-Test] = %v, want [yes]", got)
	}
	cookies := resp.Cookies()
	if len(cookies) != 1 || cookies[0].Name != "sid" || cookies[0].Value != "abc" {
		t.Fatalf("Cookies() = %+v, want sid=abc", cookies)
	}
	if got := resp.ContentLength(); got != int64(len("hello")) {
		t.Fatalf("ContentLength() = %d, want %d", got, len("hello"))
	}
}

func TestGlobalHeadersArePlainValues(t *testing.T) {
	SetGlobalHeader("X-Resty-Plain", "one")
	AddGlobalHeader("X-Resty-Plain", "two")
	defer RemoveGlobalHeader("X-Resty-Plain")

	headers := CloneGlobalHeaders()
	if got := headers["X-Resty-Plain"]; len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("CloneGlobalHeaders()[X-Resty-Plain] = %v, want [one two]", got)
	}
}

func TestTimeoutReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Timeout(time.Millisecond).Execute()
	if resp.Err() == nil {
		t.Fatal("Execute() error is nil, want timeout error")
	}
}

func TestRequestOptionsOverrideGlobalDefaults(t *testing.T) {
	oldUA := GetGlobalUserAgent()
	oldFollow := GetGlobalFollowRedirects()
	defer SetGlobalUserAgent(oldUA)
	defer SetGlobalFollowRedirects(oldFollow)

	SetGlobalUserAgent("global-resty-agent")
	SetGlobalFollowRedirects(false)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/end", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte(r.Header.Get("X-Req") + ":" + r.Header.Get("User-Agent")))
	}))
	defer srv.Close()

	resp := Get(srv.URL+"/start",
		WithHeader("X-Req", "per-call"),
		WithUserAgent("request-resty-agent"),
		WithFollowRedirects(true),
	).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if got := resp.Body(); got != "per-call:request-resty-agent" {
		t.Fatalf("Body() = %q, want per-call options to override globals", got)
	}
}

func TestRestyClientFactoryProviderLifecycle(t *testing.T) {
	ResetDefaultRestyClientProvider()
	t.Cleanup(ResetDefaultRestyClientProvider)

	defaultCalled := 0
	ConfigureDefaultRestyClientProvider(func() *grestry.Client {
		defaultCalled++
		return grestry.New()
	})
	client := NewIsolatedRequest(MethodGet, "http://example.com").buildClient()
	if client == nil || defaultCalled != 1 {
		t.Fatalf("default provider client=%v called=%d", client, defaultCalled)
	}

	perCallCalled := 0
	client = NewIsolatedRequest(MethodGet, "http://example.com", WithRestyClientFactory(func() *grestry.Client {
		perCallCalled++
		return grestry.New()
	})).buildClient()
	if client == nil || perCallCalled != 1 || defaultCalled != 1 {
		t.Fatalf("per-call factory client=%v perCall=%d default=%d", client, perCallCalled, defaultCalled)
	}

	client = NewIsolatedRequest(MethodGet, "http://example.com", WithRestyClientFactory(func() *grestry.Client { return nil })).buildClient()
	if client == nil || defaultCalled != 2 {
		t.Fatalf("nil per-call factory client=%v default=%d", client, defaultCalled)
	}

	ResetDefaultRestyClientProvider()
	client = NewIsolatedRequest(MethodGet, "http://example.com").buildClient()
	if client == nil {
		t.Fatal("reset default provider should create a client")
	}
}

func TestCreateWithOptionsAppliesRequestOptions(t *testing.T) {
	getReq := CreateGetWithOptions("http://example.com", false, WithHeader("X-Create", "get"), WithUserAgent("create-get-agent"))
	if getReq.followRedir == nil || *getReq.followRedir {
		t.Fatalf("followRedir: %v", getReq.followRedir)
	}
	if got := getReq.headers["X-Create"]; len(got) != 1 || got[0] != "get" {
		t.Fatalf("CreateGetWithOptions header = %q, want get", got)
	}
	if got := getReq.userAgent; got != "create-get-agent" {
		t.Fatalf("CreateGetWithOptions userAgent = %q", got)
	}

	postReq := CreatePostWithOptions("http://example.com", WithHeader("X-Create", "post"))
	if postReq.method != MethodPost {
		t.Fatalf("CreatePostWithOptions method = %v, want POST", postReq.method)
	}
	if got := postReq.headers["X-Create"]; len(got) != 1 || got[0] != "post" {
		t.Fatalf("CreatePostWithOptions header = %q, want post", got)
	}
}

func TestSnapshotGlobalConfigAndExplicitRequestConfig(t *testing.T) {
	oldTimeout := GetGlobalTimeout()
	oldMaxRedirects := GetGlobalMaxRedirects()
	oldFollow := GetGlobalFollowRedirects()
	oldUA := GetGlobalUserAgent()
	oldTrust := IsTrustAnyHost()
	defer SetGlobalTimeout(oldTimeout)
	defer SetGlobalMaxRedirects(oldMaxRedirects)
	defer SetGlobalFollowRedirects(oldFollow)
	defer SetGlobalUserAgent(oldUA)
	defer SetTrustAnyHost(oldTrust)
	defer RemoveGlobalHeader("X-Snapshot")

	SetGlobalTimeout(123 * time.Millisecond)
	SetGlobalMaxRedirects(3)
	SetGlobalFollowRedirects(false)
	SetGlobalUserAgent("snapshot-agent")
	SetTrustAnyHost(true)
	SetGlobalHeader("X-Snapshot", "one")

	cfg := SnapshotGlobalConfig()
	SetGlobalHeader("X-Snapshot", "mutated")
	cfg.Headers["X-Snapshot"][0] = "cfg"

	req := NewRequestWithConfig(MethodGet, "http://example.com", cfg)
	if req.timeout != 123*time.Millisecond || req.maxRedirects != 3 || req.followRedir == nil || *req.followRedir || !req.tlsSkip || req.userAgent != "snapshot-agent" {
		t.Fatalf("request config not applied: timeout=%v max=%d follow=%v tls=%v ua=%q", req.timeout, req.maxRedirects, req.followRedir, req.tlsSkip, req.userAgent)
	}
	if got := req.headers["X-Snapshot"]; len(got) != 1 || got[0] != "cfg" {
		t.Fatalf("explicit config headers = %v, want [cfg]", got)
	}
	if got := CloneGlobalHeaders()["X-Snapshot"]; len(got) != 1 || got[0] != "mutated" {
		t.Fatalf("snapshot should be detached from globals, global header = %v", got)
	}
}

func TestNewIsolatedRequestDoesNotReadGlobals(t *testing.T) {
	oldTimeout := GetGlobalTimeout()
	oldMaxRedirects := GetGlobalMaxRedirects()
	oldFollow := GetGlobalFollowRedirects()
	oldUA := GetGlobalUserAgent()
	oldTrust := IsTrustAnyHost()
	defer SetGlobalTimeout(oldTimeout)
	defer SetGlobalMaxRedirects(oldMaxRedirects)
	defer SetGlobalFollowRedirects(oldFollow)
	defer SetGlobalUserAgent(oldUA)
	defer SetTrustAnyHost(oldTrust)
	defer RemoveGlobalHeader("X-Isolated")

	SetGlobalTimeout(time.Second)
	SetGlobalMaxRedirects(1)
	SetGlobalFollowRedirects(false)
	SetGlobalUserAgent("global-agent")
	SetTrustAnyHost(true)
	SetGlobalHeader("X-Isolated", "global")

	req := NewIsolatedRequest(MethodGet, "http://example.com")
	if req.timeout != 0 || req.maxRedirects != 10 || req.followRedir == nil || !*req.followRedir || req.tlsSkip || req.userAgent != "" {
		t.Fatalf("isolated request leaked globals: timeout=%v max=%d follow=%v tls=%v ua=%q", req.timeout, req.maxRedirects, req.followRedir, req.tlsSkip, req.userAgent)
	}
	if got := req.headers["X-Isolated"]; len(got) != 0 {
		t.Fatalf("isolated request should not include global header: %v", got)
	}
}

func TestWithGlobalConfigOptionOverridesConstructionDefaults(t *testing.T) {
	cfg := GlobalConfig{
		Timeout:          250 * time.Millisecond,
		MaxRedirects:     2,
		FollowRedirects:  false,
		DefaultUserAgent: "option-agent",
		Headers:          HeaderValues{"X-Config": []string{"yes"}},
	}
	req := NewIsolatedRequest(MethodGet, "http://example.com", WithGlobalConfig(cfg), WithHeader("X-Req", "ok"))
	if req.timeout != 250*time.Millisecond || req.maxRedirects != 2 || req.followRedir == nil || *req.followRedir || req.userAgent != "option-agent" {
		t.Fatalf("WithGlobalConfig not applied: timeout=%v max=%d follow=%v ua=%q", req.timeout, req.maxRedirects, req.followRedir, req.userAgent)
	}
	if got := req.headers["X-Config"]; len(got) != 1 || got[0] != "yes" {
		t.Fatalf("config header = %v, want [yes]", got)
	}
	if got := req.headers["X-Req"]; len(got) != 1 || got[0] != "ok" {
		t.Fatalf("request header after config = %v, want [ok]", got)
	}
}

func TestRequestOptionContentTypeAndCharset(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("Content-Type")))
	}))
	defer srv.Close()

	got := Post(srv.URL, WithCharset("GBK"), WithContentType("text/custom")).BodyString("hello").Execute().Body()
	if got != "text/custom" {
		t.Fatalf("Content-Type = %q, want text/custom", got)
	}

	got = Post(srv.URL, WithCharset("GBK")).BodyJSON(`{"ok":true}`).Execute().Body()
	if got != "application/json;charset=GBK" {
		t.Fatalf("JSON Content-Type = %q", got)
	}
}

func TestRequestOptionTLSConfig(t *testing.T) {
	c := Get("https://example.com", WithTLSConfig(&tls.Config{ServerName: "example.com"})).buildClient()
	if c == nil {
		t.Fatal("client is nil")
	}
}

func TestSaveAsOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("resty-save"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "out.txt")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatalf("write old: %v", err)
	}
	if _, err := Get(srv.URL).Execute().SaveAs(target, WithSaveOverwrite(false)); err == nil {
		t.Fatal("SaveAs overwrite false should fail")
	}
	if _, err := DownloadFile(srv.URL, target); err != nil {
		t.Fatalf("DownloadFile overwrite default: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "resty-save" {
		t.Fatalf("content = %q", data)
	}
}

func TestSaveAsProviderOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("resty-provider-save"))
	}))
	defer srv.Close()

	var mkdirPath string
	var mkdirPerm fs.FileMode
	var openPath string
	var openFlag int
	var openPerm fs.FileMode
	var written bytes.Buffer
	n, err := Get(srv.URL).Execute().SaveAs("/virtual/resty.txt",
		WithSaveMkdirAll(func(path string, perm fs.FileMode) error {
			mkdirPath, mkdirPerm = path, perm
			return nil
		}),
		WithSaveOpenFile(func(path string, flag int, perm fs.FileMode) (io.WriteCloser, error) {
			openPath, openFlag, openPerm = path, flag, perm
			return nopWriteCloser{Writer: &written}, nil
		}),
		WithSaveDirPerm(0o700), WithSaveFilePerm(0o600),
	)
	if err != nil || n != int64(len("resty-provider-save")) {
		t.Fatalf("SaveAs provider n=%d err=%v", n, err)
	}
	if mkdirPath != "/virtual" || mkdirPerm != 0o700 || openPath != "/virtual/resty.txt" || openPerm != 0o600 || openFlag&os.O_CREATE == 0 || written.String() != "resty-provider-save" {
		t.Fatalf("providers mkdir=%q/%v open=%q flag=%#x perm=%v content=%q", mkdirPath, mkdirPerm, openPath, openFlag, openPerm, written.String())
	}
}

type nopWriteCloser struct{ io.Writer }

func (w nopWriteCloser) Close() error { return nil }

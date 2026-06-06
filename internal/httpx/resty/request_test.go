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

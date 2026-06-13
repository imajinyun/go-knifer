package http

import (
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAdditionalRequestOptionsAndAccessors(t *testing.T) {
	requestFactoryCalled := false
	readAllCalled := false
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    http.StatusOK,
			ContentLength: int64(len(req.Method + ":" + req.Header.Get("X-A"))),
			Header:        http.Header{"Content-Type": []string{"text/plain"}},
			Body:          io.NopCloser(strings.NewReader(req.Method + ":" + req.Header.Get("X-A"))),
			Request:       req,
		}, nil
	})
	cfg := SnapshotGlobalConfig()
	cfg.Headers.Set("X-A", "cfg")
	req := NewIsolatedRequest(MethodPost, "https://example.com/upload",
		WithGlobalConfig(cfg),
		WithHeaders(map[string]string{"X-A": "one", "X-B": "two"}),
		WithClient(&http.Client{Transport: transport}),
		WithTransport(transport),
		WithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}),
		WithResponseReadAllFunc(func(r io.Reader) ([]byte, error) {
			readAllCalled = true
			return io.ReadAll(r)
		}),
		WithRequestFactory(func(method, rawURL string, body io.Reader) (*http.Request, error) {
			requestFactoryCalled = true
			return http.NewRequest(method, rawURL, body)
		}),
		WithMultipartWriterFactory(func(w io.Writer) MultipartWriter {
			return multipart.NewWriter(w)
		}),
	)
	req.Method(MethodPatch).URL("https://example.com/changed").AddHeader("X-A", "extra").Headers(map[string]string{"X-B": "two"}).CookieString("raw=cookie")
	req.Client(&http.Client{Transport: transport}).URLPolicy(URLPolicy{AllowedSchemes: []string{"https"}, RejectPrivate: false})
	resp := req.FormFileReader("file", "a.txt", strings.NewReader("file-data")).Execute()
	if resp.Err() != nil {
		t.Fatalf("multipart Execute: %v", resp.Err())
	}
	if got := resp.Body(); got != "PATCH:one" {
		t.Fatalf("response body = %q", got)
	}
	if !requestFactoryCalled || !readAllCalled {
		t.Fatalf("providers called request=%v readAll=%v", requestFactoryCalled, readAllCalled)
	}
	if got := NewRequestWithConfig(MethodGet, "https://example.com", cfg, WithTransport(transport)).Execute().Body(); got != "GET:cfg" {
		t.Fatalf("NewRequestWithConfig body = %q", got)
	}
}

func TestRequestTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Timeout(50 * time.Millisecond).Execute()
	if resp.Err() == nil {
		t.Fatal("expected timeout error")
	}
}

func TestRequestNoFollowRedirects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/end", http.StatusFound)
			return
		}
		_, _ = w.Write([]byte("end"))
	}))
	defer srv.Close()

	resp := Get(srv.URL + "/start").FollowRedirects(false).Execute()
	if resp.Status() != http.StatusFound {
		t.Fatalf("expected 302, got %d", resp.Status())
	}

	body := Get(srv.URL + "/start").FollowRedirects(true).Execute().Body()
	if body != "end" {
		t.Fatalf("redirected body: %q", body)
	}
}

func TestRequestOptionsOverrideGlobalDefaults(t *testing.T) {
	oldUA := GetGlobalUserAgent()
	oldFollow := GetGlobalFollowRedirects()
	defer SetGlobalUserAgent(oldUA)
	defer SetGlobalFollowRedirects(oldFollow)

	SetGlobalUserAgent("global-agent")
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
		WithUserAgent("request-agent"),
		WithFollowRedirects(true),
	).Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if got := resp.Body(); got != "per-call:request-agent" {
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
	client := Get("https://example.com", WithTLSConfig(&tls.Config{ServerName: "example.com"})).buildClient()
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("transport type = %T", client.Transport)
	}
	if transport.TLSClientConfig == nil || transport.TLSClientConfig.ServerName != "example.com" {
		t.Fatalf("TLS config = %#v", transport.TLSClientConfig)
	}
}

func TestRequestOptionCookieJar(t *testing.T) {
	oldJar := GetCookieJar()
	CloseCookie()
	defer SetCookieJar(oldJar)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/set" {
			http.SetCookie(w, &http.Cookie{Name: "sid", Value: "abc", Path: "/"})
			_, _ = w.Write([]byte("set"))
			return
		}
		c, err := r.Cookie("sid")
		if err != nil {
			_, _ = w.Write([]byte("missing"))
			return
		}
		_, _ = w.Write([]byte(c.Value))
	}))
	defer srv.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	if resp := Get(srv.URL+"/set", WithCookieJar(jar)).Execute(); resp.Err() != nil {
		t.Fatalf("set cookie request error = %v", resp.Err())
	}
	if got := Get(srv.URL+"/get", WithCookieJar(jar)).Execute().Body(); got != "abc" {
		t.Fatalf("cookie jar body = %q, want abc", got)
	}
}

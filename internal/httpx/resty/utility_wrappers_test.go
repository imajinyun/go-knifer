package resty

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	grestry "resty.dev/v3"
)

func TestStringHelpersReturnErrorsExplicitly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method + ":" + r.URL.Query().Get("k")))
	}))
	defer srv.Close()

	body, err := GetWithParamsE(srv.URL, map[string]any{"k": "v"})
	if err != nil || body != "GET:v" {
		t.Fatalf("GetWithParamsE = %q, %v", body, err)
	}

	if body, err = PostStringE(srv.URL, "payload"); err != nil || body != "POST:" {
		t.Fatalf("PostStringE = %q, %v", body, err)
	}

	if _, err = GetStringE("http://[::1"); err == nil {
		t.Fatal("GetStringE invalid URL error = nil")
	}
	if _, err = DownloadBytesE("http://[::1"); err == nil {
		t.Fatal("DownloadBytesE invalid URL error = nil")
	}
	if _, err = GetStringSafeE(srv.URL); err == nil {
		t.Fatal("GetStringSafeE local URL error = nil, want private address rejection")
	}
}

func TestAdditionalRequestAndUtilityWrappers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write([]byte(r.Method + ":" + r.URL.Query().Get("q") + ":" + r.Header.Get("Authorization") + ":" + string(body)))
	}))
	defer srv.Close()

	req := NewIsolatedRequest(MethodGet, srv.URL).
		Method(MethodPost).
		URL(srv.URL).
		Headers(map[string]string{"X-A": "one"}).
		AddHeader("X-A", "two").
		CookieString("raw=cookie").
		ContentType(string(ContentTypeTextPlain)).
		Charset("utf-8").
		Timeout(time.Second).
		FollowRedirects(false).
		MaxRedirects(1).
		TLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}).
		RestyClient(grestry.New()).
		URLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false}).
		BasicAuth("user", "pass").
		BearerAuth("token").
		Query("q", "one").
		QueryMap(map[string]any{"q": "two"}).
		BodyReader(strings.NewReader("reader-body")).
		ErrorResult(&map[string]any{})
	if req.method != MethodPost || req.rawURL != srv.URL || req.urlPolicy == nil || req.errorResult == nil {
		t.Fatalf("request state method=%s url=%s policy=%#v", req.method, req.rawURL, req.urlPolicy)
	}
	resp := req.Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute: %v", resp.Err())
	}
	if !strings.Contains(resp.Body(), "POST:two:Basic ") || !strings.Contains(resp.Body(), ":reader-body") {
		t.Fatalf("response body = %q", resp.Body())
	}

	if got, err := GetWithTimeoutE(srv.URL, time.Second); err != nil || !strings.HasPrefix(got, "GET:") {
		t.Fatalf("GetWithTimeoutE = %q, %v", got, err)
	}
	if got, err := GetWithTimeoutEWithOptions(srv.URL, time.Second, WithHeader("X-T", "v")); err != nil || !strings.HasPrefix(got, "GET:") {
		t.Fatalf("GetWithTimeoutEWithOptions = %q, %v", got, err)
	}
	if got, err := PostStringSafeE(srv.URL, "safe", WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})); err != nil || !strings.Contains(got, "POST:::safe") {
		t.Fatalf("PostStringSafeE = %q, %v", got, err)
	}
	if got, err := PostFormSafeE(srv.URL, map[string]any{"a": "b"}, WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})); err != nil || !strings.HasPrefix(got, "POST:") {
		t.Fatalf("PostFormSafeE = %q, %v", got, err)
	}
	if got, err := PostJSONSafeE(srv.URL, `{"ok":true}`, WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})); err != nil || !strings.Contains(got, `{"ok":true}`) {
		t.Fatalf("PostJSONSafeE = %q, %v", got, err)
	}

	if !IsHTTP("http://example.com") || !IsHTTPS("https://example.com") {
		t.Fatal("scheme helpers returned false")
	}
	if got := BuildContentType("text/plain", "utf-8"); got != "text/plain;charset=utf-8" {
		t.Fatalf("BuildContentType = %q", got)
	}
	if !IsDefaultContentType("") || !IsFormURLEncoded("application/x-www-form-urlencoded;charset=utf-8") {
		t.Fatal("content type predicates returned unexpected result")
	}
	if got := GetCharsetFromContentTypeWithOptions("text/plain; enc=gbk", WithCharsetRegexp(regexp.MustCompile(`enc=([a-z0-9-]+)`))); got != "gbk" {
		t.Fatalf("GetCharsetFromContentTypeWithOptions = %q", got)
	}
	if got := GetCharsetFromHTMLWithOptions(`<meta data-charset="big5">`, WithMetaCharsetRegexp(regexp.MustCompile(`data-charset="([^"]+)"`))); got != "big5" {
		t.Fatalf("GetCharsetFromHTMLWithOptions = %q", got)
	}
	if got := GetMimeType("payload.JSON"); got != "application/json" {
		t.Fatalf("GetMimeType = %q", got)
	}
	if got := BuildBasicAuth("user", "pass"); !strings.HasPrefix(got, "Basic ") {
		t.Fatalf("BuildBasicAuth = %q", got)
	}
}

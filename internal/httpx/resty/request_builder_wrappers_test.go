package resty

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	grestry "resty.dev/v3"
)

func TestRequestBuilderWrappers(t *testing.T) {
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
}

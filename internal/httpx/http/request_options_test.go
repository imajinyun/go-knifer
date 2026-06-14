package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

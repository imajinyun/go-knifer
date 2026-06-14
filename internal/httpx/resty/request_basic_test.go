package resty

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

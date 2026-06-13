package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestCloneCreatesIndependentBuilder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Query().Get("q") + ":" + r.Header.Get("X-Token")))
	}))
	defer srv.Close()

	base := Get(srv.URL).Query("q", "base").Header("X-Token", "base")
	clone := base.Clone().Query("q", "clone").Header("X-Token", "clone")

	if got := base.Execute().Body(); got != "base:base" {
		t.Fatalf("base Body() = %q", got)
	}
	if got := clone.Execute().Body(); got != "base:clone" {
		t.Fatalf("clone Body() = %q", got)
	}
}

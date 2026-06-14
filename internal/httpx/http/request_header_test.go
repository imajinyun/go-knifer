package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Token")))
	}))
	defer srv.Close()

	body := Get(srv.URL).Header("X-Token", "abc").Execute().Body()
	if body != "abc" {
		t.Fatalf("body: %q", body)
	}
}

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method))
	}))
	defer srv.Close()

	body := Patch(srv.URL).Execute().Body()
	if body != http.MethodPatch {
		t.Fatalf("method: %q", body)
	}
}

func TestRequestDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Method))
	}))
	defer srv.Close()

	body := Delete(srv.URL).Execute().Body()
	if body != http.MethodDelete {
		t.Fatalf("method: %q", body)
	}
}

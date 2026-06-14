package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadString(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("plain"))
	}))
	defer srv.Close()

	got, err := DownloadStringE(srv.URL, "")
	if err != nil {
		t.Fatalf("DownloadStringE() error = %v", err)
	}
	if got != "plain" {
		t.Fatalf("body: %q", got)
	}
}

func TestDownloadStringWithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "secret" {
			http.Error(w, "missing option header", http.StatusTeapot)
			return
		}
		_, _ = w.Write([]byte("with-options"))
	}))
	defer srv.Close()

	got, err := DownloadStringEWithOptions(srv.URL, "", WithHeader("X-Token", "secret"))
	if err != nil {
		t.Fatalf("DownloadStringEWithOptions() error = %v", err)
	}
	if got != "with-options" {
		t.Fatalf("body: %q", got)
	}
}

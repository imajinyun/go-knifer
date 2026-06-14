package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadToWriter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("write-me"))
	}))
	defer srv.Close()

	buf := &bytes.Buffer{}
	n, err := Download(srv.URL, buf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != int64(len("write-me")) || buf.String() != "write-me" {
		t.Fatalf("got %d bytes %q", n, buf.String())
	}
}

func TestDownloadWithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Mode") != "writer" {
			http.Error(w, "missing option header", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("write-options"))
	}))
	defer srv.Close()

	buf := &bytes.Buffer{}
	n, err := DownloadWithOptions(srv.URL, buf, WithHeader("X-Mode", "writer"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != int64(len("write-options")) || buf.String() != "write-options" {
		t.Fatalf("got %d bytes %q", n, buf.String())
	}
}

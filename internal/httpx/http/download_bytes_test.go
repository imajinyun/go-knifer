package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadBytes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte{0x01, 0x02, 0x03})
	}))
	defer srv.Close()

	got, err := DownloadBytesE(srv.URL)
	if err != nil {
		t.Fatalf("DownloadBytesE() error = %v", err)
	}
	if !bytes.Equal(got, []byte{0x01, 0x02, 0x03}) {
		t.Fatalf("bytes: %v", got)
	}
}

func TestDownloadBytesWithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Mode") != "bytes" {
			http.Error(w, "missing option header", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte{0x04, 0x05, 0x06})
	}))
	defer srv.Close()

	got, err := DownloadBytesEWithOptions(srv.URL, WithHeader("X-Mode", "bytes"))
	if err != nil {
		t.Fatalf("DownloadBytesEWithOptions() error = %v", err)
	}
	if !bytes.Equal(got, []byte{0x04, 0x05, 0x06}) {
		t.Fatalf("bytes: %v", got)
	}
}

package url

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestOpenAndContentLengthWithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "secret" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Length", "4")
		_, _ = w.Write([]byte("body"))
	}))
	defer srv.Close()

	r, err := OpenWithOptions(srv.URL, WithHeader("X-Token", "secret"), WithTimeout(time.Second), WithCheckStatus(true))
	if err != nil {
		t.Fatalf("OpenWithOptions: %v", err)
	}
	data, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil || string(data) != "body" {
		t.Fatalf("body = %q, %v", data, err)
	}
	length, err := ContentLengthWithOptions(srv.URL, WithHeader("X-Token", "secret"), WithCheckStatus(true))
	if err != nil || length != 4 {
		t.Fatalf("ContentLengthWithOptions = %d, %v", length, err)
	}
	if _, err := OpenWithOptions(srv.URL, WithCheckStatus(true)); err == nil {
		t.Fatal("OpenWithOptions status check error = nil")
	}
}

func TestResourceProviderOptions(t *testing.T) {
	openedPath := ""
	r, err := OpenWithOptions("file:///virtual/data.txt", WithOpenFile(func(path string) (io.ReadCloser, error) {
		openedPath = path
		return io.NopCloser(strings.NewReader("virtual-body")), nil
	}))
	if err != nil {
		t.Fatalf("OpenWithOptions custom open file: %v", err)
	}
	data, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil || string(data) != "virtual-body" || openedPath != "/virtual/data.txt" {
		t.Fatalf("custom open file data=%q path=%q err=%v", data, openedPath, err)
	}

	statSource := t.TempDir() + "/stat.txt"
	if err := os.WriteFile(statSource, []byte("1234567"), 0o600); err != nil {
		t.Fatal(err)
	}
	statPath := ""
	length, err := ContentLengthWithOptions("/virtual/stat.txt", WithStat(func(path string) (os.FileInfo, error) {
		statPath = path
		return os.Stat(statSource)
	}))
	if err != nil || length != 7 || statPath != "/virtual/stat.txt" {
		t.Fatalf("custom stat length=%d path=%q err=%v", length, statPath, err)
	}
}

func TestResourceRequestFactoryOption(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Factory") != "yes" {
			http.Error(w, "factory header missing", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("factory"))
	}))
	defer srv.Close()

	var gotMethod, gotURL string
	factory := func(ctx context.Context, method, raw string) (*http.Request, error) {
		gotMethod, gotURL = method, raw
		req, err := http.NewRequestWithContext(ctx, method, raw, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Factory", "yes")
		return req, nil
	}
	r, err := OpenWithOptions(srv.URL, WithRequestFactory(factory), WithCheckStatus(true))
	if err != nil {
		t.Fatalf("OpenWithOptions with request factory: %v", err)
	}
	data, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil || string(data) != "factory" || gotMethod != http.MethodGet || gotURL != srv.URL {
		t.Fatalf("factory data=%q method=%q url=%q err=%v", data, gotMethod, gotURL, err)
	}
}

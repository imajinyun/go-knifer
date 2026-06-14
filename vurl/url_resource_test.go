package vurl_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/imajinyun/go-knifer/vurl"
)

func TestFacadeResourceOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test"); got != "facade" {
			t.Fatalf("header X-Test = %q, want facade", got)
		}
		w.Header().Set("Content-Length", "5")
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	rc, err := vurl.OpenWithOptions(server.URL, vurl.WithHeader("X-Test", "facade"), vurl.WithCheckStatus(true))
	if err != nil {
		t.Fatalf("OpenWithOptions: %v", err)
	}
	defer func() { _ = rc.Close() }()
	data, err := io.ReadAll(rc)
	if err != nil || string(data) != "hello" {
		t.Fatalf("OpenWithOptions data = %q, %v", data, err)
	}

	length, err := vurl.ContentLengthWithOptions(server.URL, vurl.WithHeader("X-Test", "facade"), vurl.WithCheckStatus(true))
	if err != nil || length != 5 {
		t.Fatalf("ContentLengthWithOptions = %d, %v; want 5, nil", length, err)
	}
	if size, err := vurl.SizeWithOptions(server.URL, vurl.WithHeader("X-Test", "facade")); err != nil || size != 5 {
		t.Fatalf("SizeWithOptions = %d, %v; want 5, nil", size, err)
	}
}

func TestFacadeResourceProviderOptions(t *testing.T) {
	openedPath := ""
	rc, err := vurl.OpenWithOptions("file:///virtual/data.txt", vurl.WithOpenFile(func(path string) (io.ReadCloser, error) {
		openedPath = path
		return io.NopCloser(strings.NewReader("facade-file")), nil
	}))
	if err != nil {
		t.Fatalf("OpenWithOptions custom open: %v", err)
	}
	data, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil || string(data) != "facade-file" || openedPath != "/virtual/data.txt" {
		t.Fatalf("custom open data=%q path=%q err=%v", data, openedPath, err)
	}

	statSource := t.TempDir() + "/stat.txt"
	if err := os.WriteFile(statSource, []byte("12345"), 0o600); err != nil {
		t.Fatal(err)
	}
	statPath := ""
	length, err := vurl.ContentLengthWithOptions("/virtual/stat.txt", vurl.WithStat(func(path string) (os.FileInfo, error) {
		statPath = path
		return os.Stat(statSource)
	}))
	if err != nil || length != 5 || statPath != "/virtual/stat.txt" {
		t.Fatalf("custom stat length=%d path=%q err=%v", length, statPath, err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Factory") != "facade" {
			http.Error(w, "missing factory header", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("factory"))
	}))
	defer server.Close()
	method := ""
	_, err = vurl.ContentLengthWithOptions(server.URL, vurl.WithRequestFactory(func(ctx context.Context, m, raw string) (*http.Request, error) {
		method = m
		req, err := http.NewRequestWithContext(ctx, m, raw, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Factory", "facade")
		return req, nil
	}), vurl.WithCheckStatus(true))
	if err != nil || method != http.MethodHead {
		t.Fatalf("request factory method=%q err=%v", method, err)
	}
}

func TestFacadeTrustedFileResourceWrappers(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "trusted.txt")
	if err := os.WriteFile(file, []byte("trusted"), 0o600); err != nil {
		t.Fatal(err)
	}
	rc, err := vurl.Open(file)
	if err != nil {
		t.Fatalf("Open file: %v", err)
	}
	data, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil || string(data) != "trusted" {
		t.Fatalf("Open file data = %q, %v", data, err)
	}
	if length, err := vurl.ContentLength(file); err != nil || length != int64(len("trusted")) {
		t.Fatalf("ContentLength = %d, %v", length, err)
	}
	if size, err := vurl.Size(file); err != nil || size != int64(len("trusted")) {
		t.Fatalf("Size = %d, %v", size, err)
	}
	if _, err := vurl.OpenWithOptions(file, vurl.WithAllowedSchemes("http")); err == nil {
		t.Fatal("OpenWithOptions disallowed scheme error = nil")
	}
	if _, err := vurl.OpenWithOptions(file, vurl.WithAllowLocalFiles(false)); err == nil {
		t.Fatal("OpenWithOptions local file disabled error = nil")
	}
	if _, err := vurl.ContentLengthWithOptions(file, vurl.WithAllowLocalFiles(false)); err == nil {
		t.Fatal("ContentLengthWithOptions local file disabled error = nil")
	}
}

func TestFacadeResourceMaxBytesAndStatusErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/large":
			w.Header().Set("Content-Length", "10")
			_, _ = w.Write([]byte("0123456789"))
		default:
			http.Error(w, "nope", http.StatusTeapot)
		}
	}))
	defer server.Close()

	if _, err := vurl.OpenWithOptions(server.URL+"/large", vurl.WithMaxBytes(2)); err == nil {
		t.Fatal("OpenWithOptions content-length max bytes error = nil")
	}
	if _, err := vurl.ContentLengthWithOptions(server.URL+"/large", vurl.WithMaxBytes(2)); err == nil {
		t.Fatal("ContentLengthWithOptions max bytes error = nil")
	}
	rc, err := vurl.OpenWithOptions(server.URL+"/large", vurl.WithMaxBytes(4), vurl.WithHTTPClient(&http.Client{Transport: stripLengthTransport{base: http.DefaultTransport}}))
	if err != nil {
		t.Fatalf("OpenWithOptions limited reader: %v", err)
	}
	_, err = io.ReadAll(rc)
	_ = rc.Close()
	if err == nil {
		t.Fatal("limited reader overflow error = nil")
	}
	if _, err := vurl.OpenWithOptions(server.URL+"/status", vurl.WithCheckStatus(true)); err == nil {
		t.Fatal("OpenWithOptions check status error = nil")
	}
}

type stripLengthTransport struct{ base http.RoundTripper }

func (t stripLengthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resp.ContentLength = -1
	resp.Header.Del("Content-Length")
	return resp, nil
}

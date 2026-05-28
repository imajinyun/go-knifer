package http

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 对应 hutool-http DownloadTest

func TestDownloadString(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("plain"))
	}))
	defer srv.Close()

	if got := DownloadString(srv.URL, ""); got != "plain" {
		t.Fatalf("body: %q", got)
	}
}

func TestDownloadBytes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte{0x01, 0x02, 0x03})
	}))
	defer srv.Close()

	got := DownloadBytes(srv.URL)
	if !bytes.Equal(got, []byte{0x01, 0x02, 0x03}) {
		t.Fatalf("bytes: %v", got)
	}
}

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

func TestDownloadFileToFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("file-content"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "out.txt")
	n, err := DownloadFile(srv.URL, target)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if n != int64(len("file-content")) {
		t.Fatalf("size: %d", n)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != "file-content" {
		t.Fatalf("content: %q", string(data))
	}
}

func TestDownloadFileToDirectory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("D"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	// dest 是目录，文件名应来自 URL path
	url := srv.URL + "/foo.bin"
	if _, err := DownloadFile(url, dir); err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "foo.bin")); err != nil {
		t.Fatalf("file should exist: %v", err)
	}
}

func TestDownloadGzipDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		_, _ = gz.Write([]byte("gzipped"))
		_ = gz.Close()
	}))
	defer srv.Close()

	body := Get(srv.URL).Execute().Body()
	if body != "gzipped" {
		t.Fatalf("decoded body: %q", body)
	}
}

func TestSaveAsViaContentDisposition(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", `attachment; filename="real.bin"`)
		_, _ = w.Write([]byte("from-cd"))
	}))
	defer srv.Close()

	dir := t.TempDir()
	resp := Get(srv.URL).Execute()
	if _, err := resp.SaveAs(dir); err != nil {
		t.Fatalf("err: %v", err)
	}
	target := filepath.Join(dir, "real.bin")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("not found: %v", err)
	}
	if !strings.Contains(string(data), "from-cd") {
		t.Fatalf("content: %q", string(data))
	}
}

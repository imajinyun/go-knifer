package vurl_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/imajinyun/go-knifer/vurl"
)

func TestFacadeQueryAndNormalize(t *testing.T) {
	if got := vurl.Normalize("example.com/a b", true, false); got != "http://example.com/a%20b" {
		t.Fatalf("Normalize: %q", got)
	}
	encoded := vurl.URLEncode("a b+c/中文")
	decoded, err := vurl.URLDecode(encoded)
	if err != nil || decoded != "a b+c/中文" {
		t.Fatalf("URL query roundtrip = %q, %v", decoded, err)
	}
	query := vurl.BuildQuery(map[string]any{"a": "1", "b": "x y"})
	if !strings.Contains(query, "a=1") || !strings.Contains(query, "b=x+y") {
		t.Fatalf("BuildQuery: %q", query)
	}
	completed, err := vurl.Complete("example.com/base/", "next")
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if completed != "http://example.com/base/next" {
		t.Fatalf("Complete: %q", completed)
	}
}

func TestFacadeChecksAndDataURI(t *testing.T) {
	if !vurl.IsWebURL("https://example.com") || vurl.IsWebURL("ftp://example.com") {
		t.Fatal("IsWebURL failed")
	}
	if !vurl.IsAbsoluteURL("ftp://example.com") {
		t.Fatal("IsAbsoluteURL failed")
	}
	if got := vurl.DataURI("text/plain", "utf-8", "base64", "aGVsbG8="); got != "data:text/plain;charset=utf-8;base64,aGVsbG8=" {
		t.Fatalf("DataURI: %q", got)
	}
}

func TestFacadeEncodeAndURLBuilder(t *testing.T) {
	if got := vurl.EncodePath("/a b/c+d"); got != "/a%20b/c+d" {
		t.Fatalf("EncodePath = %q", got)
	}
	if got := vurl.EncodePathSegment("a/b"); got != "a%2Fb" {
		t.Fatalf("EncodePathSegment = %q", got)
	}
	if got := vurl.EncodeQuery("a b+c"); got != "a+b%2Bc" {
		t.Fatalf("EncodeQuery = %q", got)
	}
	if got, _ := vurl.DecodeForPath("a+b%2Bc"); got != "a+b+c" {
		t.Fatalf("DecodeForPath = %q", got)
	}
	if got, _ := vurl.DecodeWithOptions("a+b%2Bc", vurl.WithPlusAsSpace(false)); got != "a+b+c" {
		t.Fatalf("DecodeWithOptions = %q", got)
	}
	built := vurl.NewHTTPURLBuilder("example.com").AddPathSegment("a b").AddQuery("q", "go net").SetFragment("top 1").Build()
	if built != "http://example.com/a%20b?q=go+net#top%201" {
		t.Fatalf("URLBuilder = %q", built)
	}
}

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

func TestFacadeNormalizeWithOptions(t *testing.T) {
	got := vurl.NormalizeWithOptions("example.com/a b", true, false, vurl.WithDefaultScheme("https"))
	if got != "https://example.com/a%20b" {
		t.Fatalf("NormalizeWithOptions = %q", got)
	}
	got = vurl.NormalizeUsingOptions("example.com//a b", vurl.WithDefaultScheme("https"), vurl.WithEncodePath(true), vurl.WithReplaceSlash(true))
	if got != "https://example.com/a%20b" {
		t.Fatalf("NormalizeUsingOptions = %q", got)
	}
}

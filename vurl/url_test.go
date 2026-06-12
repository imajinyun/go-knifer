package vurl_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestFacadeAdditionalEncodingAndQueryHelpers(t *testing.T) {
	if got := vurl.EncodeWithOptions("a b", vurl.WithQueryEscapeFunc(func(s string) string { return "escaped:" + s })); got != "escaped:a b" {
		t.Fatalf("EncodeWithOptions = %q", got)
	}
	if got := vurl.URLEncodeWithOptions("a b", vurl.WithQueryEscapeFunc(strings.ToUpper)); got != "A B" {
		t.Fatalf("URLEncodeWithOptions = %q", got)
	}
	if got := vurl.EncodeQueryWithOptions("a b", vurl.WithQueryEscapeFunc(func(s string) string { return "q:" + s })); got != "q:a b" {
		t.Fatalf("EncodeQueryWithOptions = %q", got)
	}
	if got := vurl.EncodePathSegmentWithOptions("a/b", vurl.WithPathEscapeFunc(func(s string) string { return "p:" + s })); got != "p:a/b" {
		t.Fatalf("EncodePathSegmentWithOptions = %q", got)
	}
	if got := vurl.FormURLEncodeWithOptions("a b", vurl.WithQueryEscapeFunc(func(s string) string { return "form:" + s })); got != "form:a b" {
		t.Fatalf("FormURLEncodeWithOptions = %q", got)
	}
	if got := vurl.EncodeAll("a b"); got != "a%20b" {
		t.Fatalf("EncodeAll = %q", got)
	}
	if got := vurl.EncodeFragment("a b#c"); got != "a%20b%23c" {
		t.Fatalf("EncodeFragment = %q", got)
	}

	decoded, err := vurl.Decode("a+b%2Fc")
	if err != nil || decoded != "a b/c" {
		t.Fatalf("Decode = %q, %v", decoded, err)
	}
	decoded, err = vurl.DecodePlus("a+b%2Fc", false)
	if err != nil || decoded != "a+b/c" {
		t.Fatalf("DecodePlus(false) = %q, %v", decoded, err)
	}

	if got := vurl.EncodeParams("https://example.com/search?q=a b&lang=go"); got != "https://example.com/search?lang=go&q=a+b" {
		t.Fatalf("EncodeParams = %q", got)
	}
	if got := vurl.DecodeQueryFirst("a=1&a=2&b=x+y"); got["a"] != "1" || got["b"] != "x y" {
		t.Fatalf("DecodeQueryFirst = %#v", got)
	}
	if got := vurl.DecodeQuery("a=1&a=2"); len(got["a"]) != 2 || got["a"][0] != "1" || got["a"][1] != "2" {
		t.Fatalf("DecodeQuery = %#v", got)
	}
	if got := vurl.AppendQuery("https://example.com/path?x=1", map[string]any{"q": "go net"}); got != "https://example.com/path?x=1&q=go+net" {
		t.Fatalf("AppendQuery = %q", got)
	}
}

func TestFacadePathAndSchemeHelpers(t *testing.T) {
	path, err := vurl.Path("https://example.com/a%20b/file.txt?x=1")
	if err != nil || path != "/a b/file.txt" {
		t.Fatalf("Path = %q, %v", path, err)
	}
	u, err := url.Parse("file:///tmp/demo.jar")
	if err != nil {
		t.Fatal(err)
	}
	if got := vurl.DecodedPath(u); got != "/tmp/demo.jar" {
		t.Fatalf("DecodedPath = %q", got)
	}
	if !vurl.IsFileURL(u) {
		t.Fatal("IsFileURL(file URL) = false")
	}
	if !vurl.IsJarFileURL(u) {
		t.Fatal("IsJarFileURL(file .jar URL) = false")
	}
	jar, err := url.Parse("jar:file:///tmp/demo.jar!/BOOT-INF/classes")
	if err != nil {
		t.Fatal(err)
	}
	if !vurl.IsJarURL(jar) {
		t.Fatal("IsJarURL(jar URL) = false")
	}
	if uri, err := vurl.ToURI("https://example.com/a b", true); err != nil || uri.String() != "https://example.com/a%20b" {
		t.Fatalf("ToURI = %v, %v", uri, err)
	}
	if !vurl.IsHTTP("http://example.com") || !vurl.IsHTTPS("https://example.com") || !vurl.IsHTTPSURL("https://example.com") {
		t.Fatal("HTTP/HTTPS scheme helpers failed")
	}
	if got := vurl.DataURIBase64("text/plain", "aGVsbG8="); got != "data:text/plain;base64,aGVsbG8=" {
		t.Fatalf("DataURIBase64 = %q", got)
	}
}

func TestFacadeSafeResourceHelpersRejectLocalhost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "2")
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	if rc, err := vurl.OpenSafe(server.URL); err == nil {
		_ = rc.Close()
		t.Fatal("OpenSafe(localhost) error = nil, want private host rejection")
	}
	if _, err := vurl.OpenSafeWithOptions(server.URL); err == nil {
		t.Fatal("OpenSafeWithOptions(localhost) error = nil, want private host rejection")
	}
	if _, err := vurl.ContentLengthSafe(server.URL); err == nil {
		t.Fatal("ContentLengthSafe(localhost) error = nil, want private host rejection")
	}
	if _, err := vurl.ContentLengthSafeWithOptions(server.URL); err == nil {
		t.Fatal("ContentLengthSafeWithOptions(localhost) error = nil, want private host rejection")
	}
}

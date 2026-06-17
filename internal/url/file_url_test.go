package url

import (
	neturl "net/url"
	"path/filepath"
	"testing"
)

func TestHostDecodedPathAndJar(t *testing.T) {
	u, _ := neturl.Parse("https://example.com/a%20b?q=1")
	if got := Host(u).String(); got != "https://example.com" {
		t.Fatalf("Host: %q", got)
	}
	if got := DecodedPath(u); got != "/a b" {
		t.Fatalf("DecodedPath: %q", got)
	}
	jar, _ := neturl.Parse("file:///tmp/a.jar")
	if !IsFileURL(jar) || !IsJarFileURL(jar) {
		t.Fatal("jar file checks failed")
	}
}

func TestParseStringURIFileURLAndPathHelpers(t *testing.T) {
	if u, err := Parse(""); err != nil || u != nil {
		t.Fatalf("Parse empty = %#v, %v", u, err)
	}
	u, err := Parse(" https://example.com/a b ")
	if err != nil {
		t.Fatalf("Parse trimmed: %v", err)
	}
	if u.String() != "https://example.com/a%20b" {
		t.Fatalf("Parse = %q", u.String())
	}
	if _, err := Parse("http://[::1"); err == nil {
		t.Fatal("Parse invalid URL error = nil")
	}

	if got := StringURI(""); got != "" {
		t.Fatalf("StringURI empty = %q", got)
	}
	if got := StringURI("payload"); got != "string:///payload" {
		t.Fatalf("StringURI = %q", got)
	}
	if got := StringURI("string:///payload"); got != "string:///payload" {
		t.Fatalf("StringURI idempotent = %q", got)
	}

	dir := t.TempDir()
	fileURL, err := FileURL(filepath.Join(dir, "data.txt"))
	if err != nil {
		t.Fatalf("FileURL: %v", err)
	}
	if fileURL.Scheme != URLProtocolFile || !filepath.IsAbs(filepath.FromSlash(fileURL.Path)) {
		t.Fatalf("FileURL = %s", fileURL.String())
	}
	if _, err := FileURL(""); err == nil {
		t.Fatal("FileURL empty path error = nil")
	}
	urls, err := FileURLs(filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt"))
	if err != nil || len(urls) != 2 {
		t.Fatalf("FileURLs len=%d err=%v", len(urls), err)
	}
	if _, err := FileURLs("ok", ""); err == nil {
		t.Fatal("FileURLs should stop on invalid path")
	}

	if Host(nil) != nil {
		t.Fatal("Host(nil) should return nil")
	}
	if p, err := Path("https://example.com/a%20b?q=1"); err != nil || p != "/a b" {
		t.Fatalf("Path = %q, %v", p, err)
	}
	if p, err := Path(""); err != nil || p != "" {
		t.Fatalf("Path empty = %q, %v", p, err)
	}
	if _, err := Path("http://[::1"); err == nil {
		t.Fatal("Path invalid URL error = nil")
	}
}

func TestToURIAndJarSchemeHelpers(t *testing.T) {
	u, err := ToURI(" https://example.com/a b ", true)
	if err != nil {
		t.Fatalf("ToURI encoded: %v", err)
	}
	if u.EscapedPath() != "/a%20b" {
		t.Fatalf("ToURI encoded path = %q", u.EscapedPath())
	}
	u, err = ToURI("https://example.com/a b", false)
	if err != nil {
		t.Fatalf("ToURI raw: %v", err)
	}
	if u.Path != "/a b" {
		t.Fatalf("ToURI raw path = %q", u.Path)
	}
	if _, err := ToURI("http://[::1", false); err == nil {
		t.Fatal("ToURI invalid URL error = nil")
	}

	for _, raw := range []string{"jar:file:///app.jar!/x", "zip:file:///app.zip!/x", "vfszip:/app.zip", "wsjar:file:///app.jar!/x"} {
		u, err := neturl.Parse(raw)
		if err != nil {
			t.Fatalf("parse %q: %v", raw, err)
		}
		if !IsJarURL(u) {
			t.Fatalf("IsJarURL(%q) = false", raw)
		}
	}
	if IsJarURL(nil) {
		t.Fatal("IsJarURL(nil) = true")
	}
}

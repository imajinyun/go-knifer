package url

import (
	"net/url"
	"strings"
	"testing"
)

func TestEncodeBlankAndParseHTTP(t *testing.T) {
	if got := EncodeBlank("https://example.com/a b"); got != "https://example.com/a%20b" {
		t.Fatalf("EncodeBlank: %q", got)
	}
	u, err := ParseHTTP("https://example.com/a b")
	if err != nil {
		t.Fatalf("ParseHTTP: %v", err)
	}
	if u.EscapedPath() != "/a%20b" {
		t.Fatalf("path: %q", u.EscapedPath())
	}
}

func TestNormalizeAndComplete(t *testing.T) {
	if got := Normalize("\\example.com\\a b", true, true); got != "http://example.com/a%20b" {
		t.Fatalf("Normalize: %q", got)
	}
	got, err := Complete("example.com/dir/", "a.html")
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if got != "http://example.com/dir/a.html" {
		t.Fatalf("Complete got %q", got)
	}
}

func TestQueryHelpers(t *testing.T) {
	queryPart := URLEncode("a b&c=d")
	if queryPart != "a+b%26c%3Dd" {
		t.Fatalf("URLEncode: %q", queryPart)
	}
	decoded, err := URLDecode(queryPart)
	if err != nil || decoded != "a b&c=d" {
		t.Fatalf("URLDecode: %v %q", err, decoded)
	}
	encoded := BuildQuery(map[string]any{"a": "1", "b": "x y", "": "skip"})
	if !strings.Contains(encoded, "a=1") || !strings.Contains(encoded, "b=x+y") || strings.Contains(encoded, "skip") {
		t.Fatalf("BuildQuery: %q", encoded)
	}
	if got := EncodeParams("https://example.com/?q=a b"); got != "https://example.com/?q=a+b" {
		t.Fatalf("EncodeParams: %q", got)
	}
	if got := DecodeQueryFirst("a=1&a=2&b=x+y"); got["a"] == "" || got["b"] != "x y" {
		t.Fatalf("DecodeQueryFirst: %#v", got)
	}
	if got := AppendQuery("https://example.com/path?x=1", map[string]any{"y": 2}); !strings.Contains(got, "x=1") || !strings.Contains(got, "y=2") {
		t.Fatalf("AppendQuery: %q", got)
	}
}

func TestURLChecksAndDataURI(t *testing.T) {
	if !IsHTTP("Http://example.com") || !IsHTTPS("HTTPS://example.com") {
		t.Fatal("scheme prefix checks failed")
	}
	if !IsWebURL("https://example.com/a") || IsWebURL("ftp://example.com/a") {
		t.Fatal("web URL checks failed")
	}
	if !IsAbsoluteURL("ftp://example.com/a") || IsAbsoluteURL("/relative") {
		t.Fatal("absolute URL checks failed")
	}
	if got := DataURIBase64("image/png", "AAAA"); got != "data:image/png;base64,AAAA" {
		t.Fatalf("DataURIBase64: %q", got)
	}
}

func TestHostDecodedPathAndJar(t *testing.T) {
	u, _ := url.Parse("https://example.com/a%20b?q=1")
	if got := Host(u).String(); got != "https://example.com" {
		t.Fatalf("Host: %q", got)
	}
	if got := DecodedPath(u); got != "/a b" {
		t.Fatalf("DecodedPath: %q", got)
	}
	jar, _ := url.Parse("file:///tmp/a.jar")
	if !IsFileURL(jar) || !IsJarFileURL(jar) {
		t.Fatal("jar file checks failed")
	}
}

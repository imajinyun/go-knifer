package vurl_test

import (
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
	built := vurl.NewHTTPURLBuilder("example.com").AddPathSegment("a b").AddQuery("q", "go net").SetFragment("top 1").Build()
	if built != "http://example.com/a%20b?q=go+net#top%201" {
		t.Fatalf("URLBuilder = %q", built)
	}
}

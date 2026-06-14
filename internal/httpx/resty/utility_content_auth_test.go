package resty

import (
	"regexp"
	"strings"
	"testing"
)

func TestContentCharsetMimeAndAuthUtilities(t *testing.T) {
	if !IsHTTP("http://example.com") || !IsHTTPS("https://example.com") {
		t.Fatal("scheme helpers returned false")
	}
	if got := BuildContentType("text/plain", "utf-8"); got != "text/plain;charset=utf-8" {
		t.Fatalf("BuildContentType = %q", got)
	}
	if !IsDefaultContentType("") || !IsFormURLEncoded("application/x-www-form-urlencoded;charset=utf-8") {
		t.Fatal("content type predicates returned unexpected result")
	}
	if got := GetCharsetFromContentTypeWithOptions("text/plain; enc=gbk", WithCharsetRegexp(regexp.MustCompile(`enc=([a-z0-9-]+)`))); got != "gbk" {
		t.Fatalf("GetCharsetFromContentTypeWithOptions = %q", got)
	}
	if got := GetCharsetFromHTMLWithOptions(`<meta data-charset="big5">`, WithMetaCharsetRegexp(regexp.MustCompile(`data-charset="([^"]+)"`))); got != "big5" {
		t.Fatalf("GetCharsetFromHTMLWithOptions = %q", got)
	}
	if got := GetMimeType("payload.JSON"); got != "application/json" {
		t.Fatalf("GetMimeType = %q", got)
	}
	if got := BuildBasicAuth("user", "pass"); !strings.HasPrefix(got, "Basic ") {
		t.Fatalf("BuildBasicAuth = %q", got)
	}
}

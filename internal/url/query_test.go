package url

import (
	"strings"
	"testing"
)

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
	if got := DecodeQuery("a=1&a=2&b=x+y"); len(got["a"]) != 2 || got["b"][0] != "x y" {
		t.Fatalf("DecodeQuery: %#v", got)
	}
	if got := DecodeQuery("bad=%zz"); len(got) != 0 {
		t.Fatalf("DecodeQuery invalid = %#v", got)
	}
	if got := AppendQuery("https://example.com/path?x=1", map[string]any{"y": 2}); !strings.Contains(got, "x=1") || !strings.Contains(got, "y=2") {
		t.Fatalf("AppendQuery: %q", got)
	}
}

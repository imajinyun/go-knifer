package url

import "testing"

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

func TestEncodeAndURLBuilder(t *testing.T) {
	if got := EncodePath("/a b/c+d"); got != "/a%20b/c+d" {
		t.Fatalf("EncodePath = %q", got)
	}
	if got := EncodePathSegment("a/b"); got != "a%2Fb" {
		t.Fatalf("EncodePathSegment = %q", got)
	}
	if got := EncodeQuery("a b+c"); got != "a+b%2Bc" {
		t.Fatalf("EncodeQuery = %q", got)
	}
	if got := EncodeQueryWithOptions("a b", WithQueryEscapeFunc(func(s string) string { return "query:" + s })); got != "query:a b" {
		t.Fatalf("EncodeQueryWithOptions = %q", got)
	}
	if got := FormURLEncodeWithOptions("a b", WithQueryEscapeFunc(func(s string) string { return "form:" + s })); got != "form:a b" {
		t.Fatalf("FormURLEncodeWithOptions = %q", got)
	}
	if got := EncodePathSegmentWithOptions("a/b", WithPathEscapeFunc(func(s string) string { return "path:" + s })); got != "path:a/b" {
		t.Fatalf("EncodePathSegmentWithOptions = %q", got)
	}
	if got := EncodeWithOptions("a b", WithQueryEscapeFunc(func(s string) string { return "encode:" + s })); got != "encode:a b" {
		t.Fatalf("EncodeWithOptions = %q", got)
	}
	if got, _ := DecodePlus("a+b%2Bc", false); got != "a+b+c" {
		t.Fatalf("DecodePlus = %q", got)
	}
	if got, _ := DecodeWithOptions("a+b%2Bc", WithPlusAsSpace(false)); got != "a+b+c" {
		t.Fatalf("DecodeWithOptions = %q", got)
	}
	built := NewHTTPURLBuilder("example.com").AddPathSegment("a b").AddQuery("q", "go net").SetFragment("top 1").Build()
	if built != "http://example.com/a%20b?q=go+net#top%201" {
		t.Fatalf("URLBuilder = %q", built)
	}
}

func TestEncodeHelpersCoverPathFragmentAndFallbacks(t *testing.T) {
	decoded, err := DecodeForPath("a+b%2Fc")
	if err != nil {
		t.Fatalf("DecodeForPath: %v", err)
	}
	if decoded != "a+b/c" {
		t.Fatalf("DecodeForPath = %q", decoded)
	}

	if got := EncodeAll("AZ az-._~:/?"); got != "AZ%20az-._~%3A%2F%3F" {
		t.Fatalf("EncodeAll = %q", got)
	}
	if got := EncodeFragment("a b/?#中"); got != "a%20b/?%23%E4%B8%AD" {
		t.Fatalf("EncodeFragment = %q", got)
	}
	if got := FormURLEncode("a b+c"); got != "a+b%2Bc" {
		t.Fatalf("FormURLEncode = %q", got)
	}
	if got := URLEncodeWithOptions("x y", WithQueryEscapeFunc(func(string) string { return "custom" })); got != "custom" {
		t.Fatalf("URLEncodeWithOptions = %q", got)
	}
	if got := EncodeQueryWithOptions("a b", WithQueryEscapeFunc(nil)); got != "a+b" {
		t.Fatalf("EncodeQueryWithOptions nil fallback = %q", got)
	}
	if got := EncodePathSegmentWithOptions("a/b", WithPathEscapeFunc(nil)); got != "a%2Fb" {
		t.Fatalf("EncodePathSegmentWithOptions nil fallback = %q", got)
	}
}

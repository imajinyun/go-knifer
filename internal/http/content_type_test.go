package http

import "testing"

// Mirrors hutool-http ContentTypeTest.
func TestContentTypeBuild(t *testing.T) {
	got := ContentTypeJSON.WithCharset("UTF-8")
	want := "application/json;charset=UTF-8"
	if got != want {
		t.Fatalf("WithCharset = %q, want %q", got, want)
	}
}

func TestContentTypeBuildFunc(t *testing.T) {
	got := BuildContentType(string(ContentTypeJSON), "UTF-8")
	if got != "application/json;charset=UTF-8" {
		t.Fatalf("BuildContentType = %q", got)
	}
}

func TestContentTypeGetWithLeadingSpace(t *testing.T) {
	json := " {\n     \"name\": \"hutool\"\n }"
	if got := GuessContentType(json); got != ContentTypeJSON {
		t.Fatalf("GuessContentType = %v", got)
	}
}

func TestContentTypeGuess(t *testing.T) {
	cases := []struct {
		in   string
		want ContentType
	}{
		{`{"a":1}`, ContentTypeJSON},
		{`[1,2,3]`, ContentTypeJSON},
		{`<root/>`, ContentTypeXML},
		{``, ""},
		{`hello`, ""},
	}
	for _, c := range cases {
		if got := GuessContentType(c.in); got != c.want {
			t.Fatalf("GuessContentType(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestIsDefaultContentType(t *testing.T) {
	if !IsDefaultContentType("") {
		t.Fatal("empty should be default")
	}
	if !IsDefaultContentType("application/x-www-form-urlencoded") {
		t.Fatal("form-urlencoded should be default")
	}
	if IsDefaultContentType("application/json") {
		t.Fatal("json should not be default")
	}
}

func TestIsFormURLEncoded(t *testing.T) {
	if !IsFormURLEncoded("application/x-www-form-urlencoded") {
		t.Fatal("expect true")
	}
	if !IsFormURLEncoded("APPLICATION/X-WWW-FORM-URLENCODED;charset=UTF-8") {
		t.Fatal("case-insensitive prefix")
	}
	if IsFormURLEncoded("application/json") {
		t.Fatal("expect false")
	}
}

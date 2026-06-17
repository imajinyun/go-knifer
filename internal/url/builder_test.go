package url

import "testing"

func TestParseURLBuilderAndAccessors(t *testing.T) {
	b, err := ParseURLBuilder("https://example.com:8443/a%20b/c?q=go+net&q=again#frag%20ment")
	if err != nil {
		t.Fatalf("ParseURLBuilder: %v", err)
	}
	if b.Scheme() != "https" || b.SchemeWithDefault("http") != "https" {
		t.Fatalf("scheme accessors = %q", b.Scheme())
	}
	if b.Host() != "example.com" || b.Authority() != "example.com:8443" {
		t.Fatalf("host/authority = %q/%q", b.Host(), b.Authority())
	}
	if b.Port() != 8443 || b.PortWithDefault(80) != 8443 {
		t.Fatalf("port accessors = %d", b.Port())
	}
	if got := b.PathString(); got != "/a%20b/c" {
		t.Fatalf("PathString = %q", got)
	}
	if got := b.Query()["q"]; len(got) != 2 || got[0] != "go net" || got[1] != "again" {
		t.Fatalf("Query = %#v", b.Query())
	}
	if b.QueryString() != "q=go+net&q=again" {
		t.Fatalf("QueryString = %q", b.QueryString())
	}
	if b.Fragment() != "frag ment" || b.FragmentEncoded() != "frag%20ment" {
		t.Fatalf("fragment = %q encoded=%q", b.Fragment(), b.FragmentEncoded())
	}
	if b.String() != "https://example.com:8443/a%20b/c?q=go+net&q=again#frag%20ment" {
		t.Fatalf("String = %q", b.String())
	}
}

func TestURLBuilderPathQueryAndDefaults(t *testing.T) {
	b := NewURLBuilder()
	if b.SchemeWithDefault("https") != "https" || b.PortWithDefault(443) != 443 || b.PathString() != "" {
		t.Fatalf("default accessors failed")
	}
	if got := b.SetWithEndTag(true).PathString(); got != "/" {
		t.Fatalf("empty PathString with end tag = %q", got)
	}

	b.SetScheme("https").SetHost("example.com").SetPort(443).
		SetPath("/api/v1").AddPath("users/list").AddPathSegment("a b").
		SetQuery("?q=go+net").AddQuery("page", 2).SetFragment("top/1")
	if got := b.Build(); got != "https://example.com:443/api/v1/users/list/a%20b/?page=2&q=go+net#top/1" {
		t.Fatalf("Build = %q", got)
	}
	oldQuery := b.QueryString()
	b.SetQuery("bad=%zz")
	if b.QueryString() != oldQuery {
		t.Fatalf("invalid SetQuery changed query from %q to %q", oldQuery, b.QueryString())
	}
}

func TestParseURLBuilderInvalidRaw(t *testing.T) {
	if _, err := ParseURLBuilder("http://[::1"); err == nil {
		t.Fatal("ParseURLBuilder invalid URL error = nil")
	}
}

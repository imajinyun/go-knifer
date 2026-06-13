package url

import "testing"

func TestNormalizeAndComplete(t *testing.T) {
	if got := Normalize("\\example.com\\a b", true, true); got != "http://example.com/a%20b" {
		t.Fatalf("Normalize: %q", got)
	}
	if got := NormalizeWithOptions("example.com/a", false, false, WithDefaultScheme("https")); got != "https://example.com/a" {
		t.Fatalf("NormalizeWithOptions: %q", got)
	}
	got, err := Complete("example.com/dir/", "a.html")
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if got != "http://example.com/dir/a.html" {
		t.Fatalf("Complete got %q", got)
	}
	if got := NormalizeUsingOptions("example.com//a b", WithDefaultScheme("https"), WithEncodePath(true), WithReplaceSlash(true)); got != "https://example.com/a%20b" {
		t.Fatalf("NormalizeUsingOptions: %q", got)
	}
}

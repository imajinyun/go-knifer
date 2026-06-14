package http

import "testing"

func TestParseLinuxDesktop(t *testing.T) {
	uaStr := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36"
	ua := ParseUserAgent(uaStr)
	if ua.OS != "Linux" {
		t.Fatalf("os: %q", ua.OS)
	}
	if ua.Platform != "Linux" {
		t.Fatalf("platform: %q", ua.Platform)
	}
	if ua.IsMobile {
		t.Fatal("desktop should not be mobile")
	}
}

func TestParseUnknown(t *testing.T) {
	ua := ParseUserAgent("Random/1.0")
	if ua.Browser != "Unknown" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.OS != "Unknown" {
		t.Fatalf("os: %q", ua.OS)
	}
	if ua.Engine != "Unknown" {
		t.Fatalf("engine: %q", ua.Engine)
	}
}

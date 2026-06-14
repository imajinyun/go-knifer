package http

import "testing"

func TestParseMobileSafariIPhone(t *testing.T) {
	uaStr := "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	ua := ParseUserAgent(uaStr)
	if ua.Browser != "Safari" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.Version != "5.0.2" {
		t.Fatalf("version: %q", ua.Version)
	}
	if ua.OS != "iOS" {
		t.Fatalf("os: %q", ua.OS)
	}
	if ua.Platform != "iPhone" {
		t.Fatalf("platform: %q", ua.Platform)
	}
	if !ua.IsMobile {
		t.Fatal("should be mobile")
	}
}

func TestParseChromeAndroid(t *testing.T) {
	uaStr := "Mozilla/5.0 (Linux; Android 9; MIX 3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.80 Mobile Safari/537.36"
	ua := ParseUserAgent(uaStr)
	if ua.Browser != "Chrome" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.Version != "70.0.3538.80" {
		t.Fatalf("version: %q", ua.Version)
	}
	if ua.OS != "Android" {
		t.Fatalf("os: %q", ua.OS)
	}
	if ua.Platform != "Android" {
		t.Fatalf("platform: %q", ua.Platform)
	}
	if !ua.IsMobile {
		t.Fatal("should be mobile")
	}
}

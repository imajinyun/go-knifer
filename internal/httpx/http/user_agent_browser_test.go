package http

import "testing"

// Covers representative browser cases from the utility toolkit-http useragent/UserAgentUtilTest.

func TestParseDesktopChromeWindows7(t *testing.T) {
	uaStr := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1"
	ua := ParseUserAgent(uaStr)
	if ua.Browser != "Chrome" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.Version != "14.0.835.163" {
		t.Fatalf("version: %q", ua.Version)
	}
	if ua.Engine != "WebKit" {
		t.Fatalf("engine: %q", ua.Engine)
	}
	if ua.OS != "Windows 7" {
		t.Fatalf("os: %q", ua.OS)
	}
	if ua.Platform != "Windows" {
		t.Fatalf("platform: %q", ua.Platform)
	}
	if ua.IsMobile {
		t.Fatal("desktop should not be mobile")
	}
}

func TestParseChromeWindows10(t *testing.T) {
	uaStr := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36"
	ua := ParseUserAgent(uaStr)
	if ua.Browser != "Chrome" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.Version != "70.0.3538.102" {
		t.Fatalf("version: %q", ua.Version)
	}
	if ua.OS != "Windows 10/11" {
		t.Fatalf("os: %q", ua.OS)
	}
	if ua.IsMobile {
		t.Fatal("should not be mobile")
	}
}

func TestParseMSEdge(t *testing.T) {
	uaStr := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.69 Safari/537.36 Edg/81.0.416.34"
	ua := ParseUserAgent(uaStr)
	if ua.Browser != "Edge" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.Version != "81.0.416.34" {
		t.Fatalf("version: %q", ua.Version)
	}
}

func TestParseFirefox(t *testing.T) {
	uaStr := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0"
	ua := ParseUserAgent(uaStr)
	if ua.Browser != "Firefox" {
		t.Fatalf("browser: %q", ua.Browser)
	}
	if ua.Version != "120.0" {
		t.Fatalf("version: %q", ua.Version)
	}
	if ua.Engine != "Gecko" {
		t.Fatalf("engine: %q", ua.Engine)
	}
	if ua.OS != "macOS" {
		t.Fatalf("os: %q", ua.OS)
	}
}

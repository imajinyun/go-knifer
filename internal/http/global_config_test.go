package http

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

// Covers the utility toolkit-http HttpGlobalConfigTest.

func TestGlobalTimeout(t *testing.T) {
	old := GetGlobalTimeout()
	defer SetGlobalTimeout(old)

	SetGlobalTimeout(7 * time.Second)
	if got := GetGlobalTimeout(); got != 7*time.Second {
		t.Fatalf("timeout: %v", got)
	}
}

func TestGlobalUserAgent(t *testing.T) {
	old := GetGlobalUserAgent()
	defer SetGlobalUserAgent(old)

	SetGlobalUserAgent("gokit-test/1.0")
	if got := GetGlobalUserAgent(); got != "gokit-test/1.0" {
		t.Fatalf("ua: %q", got)
	}
}

func TestGlobalFollowRedirects(t *testing.T) {
	old := GetGlobalFollowRedirects()
	defer SetGlobalFollowRedirects(old)

	SetGlobalFollowRedirects(false)
	if GetGlobalFollowRedirects() {
		t.Fatal("expected false")
	}
}

func TestGlobalMaxRedirects(t *testing.T) {
	old := GetGlobalMaxRedirects()
	defer SetGlobalMaxRedirects(old)

	SetGlobalMaxRedirects(3)
	if got := GetGlobalMaxRedirects(); got != 3 {
		t.Fatalf("max: %d", got)
	}
}

func TestGlobalIgnoreEOFError(t *testing.T) {
	old := IsIgnoreEOFError()
	defer SetIgnoreEOFError(old)

	SetIgnoreEOFError(false)
	if IsIgnoreEOFError() {
		t.Fatal("expected false")
	}
}

func TestGlobalTrustAnyHost(t *testing.T) {
	old := IsTrustAnyHost()
	defer SetTrustAnyHost(old)

	SetTrustAnyHost(true)
	if !IsTrustAnyHost() {
		t.Fatal("expected true")
	}
}

func TestGlobalHeadersDefault(t *testing.T) {
	headers := CloneGlobalHeaders()
	if headers.Get("User-Agent") == "" {
		t.Fatal("default UA missing")
	}
	if headers.Get("Accept") == "" {
		t.Fatal("default Accept missing")
	}
	if got := headers.Get("Accept-Encoding"); strings.Contains(got, "br") {
		t.Fatalf("default Accept-Encoding = %q should not advertise br without brotli decoding support", got)
	}
}

func TestGlobalHeadersSetAndRemove(t *testing.T) {
	SetGlobalHeader("X-Test", "v1")
	defer RemoveGlobalHeader("X-Test")

	headers := CloneGlobalHeaders()
	if headers.Get("X-Test") != "v1" {
		t.Fatalf("X-Test: %q", headers.Get("X-Test"))
	}

	RemoveGlobalHeader("X-Test")
	if got := CloneGlobalHeaders().Get("X-Test"); got != "" {
		t.Fatalf("after remove: %q", got)
	}
}

func TestGlobalCookieJar(t *testing.T) {
	jar := GetCookieJar()
	if jar == nil {
		t.Fatal("default jar should not be nil")
	}
	CloseCookie()
	if GetCookieJar() != nil {
		t.Fatal("after close should be nil")
	}
	// Restore the default jar.
	SetCookieJar(jar)
	if GetCookieJar() == nil {
		t.Fatal("restored jar nil")
	}

	// Customize the jar.
	var custom http.CookieJar
	SetCookieJar(custom)
	if GetCookieJar() != nil {
		t.Fatal("custom nil jar")
	}
	SetCookieJar(jar)
}

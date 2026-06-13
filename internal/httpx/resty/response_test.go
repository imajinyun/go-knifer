package resty

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseReadLimitOptions(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("abcdef"))
	}))
	defer srv.Close()

	limited := Get(srv.URL, WithMaxResponseBytes(3)).Execute()
	if got := limited.Bytes(); len(got) != 0 || limited.Err() == nil {
		t.Fatalf("limited Bytes() = %q err=%v, want max bytes error", string(got), limited.Err())
	}

	SetGlobalMaxResponseBytes(3)
	globalLimited := Get(srv.URL).Execute()
	SetGlobalMaxResponseBytes(0)
	if got := globalLimited.Bytes(); len(got) != 0 || globalLimited.Err() == nil {
		t.Fatalf("global limited Bytes() = %q err=%v, want max bytes error", string(got), globalLimited.Err())
	}

	unlimited := Get(srv.URL, WithMaxResponseBytes(0)).Execute()
	if got := unlimited.Body(); got != "abcdef" || unlimited.Err() != nil {
		t.Fatalf("unlimited override Body() = %q err=%v", got, unlimited.Err())
	}
}

func TestResponseHeadersCookiesAndLength(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Cookie"); !strings.Contains(got, "k=v") {
			t.Fatalf("Cookie = %q, want k=v", got)
		}
		w.Header().Set("X-Test", "yes")
		w.Header().Add("Set-Cookie", "sid=abc; Path=/")
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	resp := Get(srv.URL).Cookie("k", "v").Execute()
	if resp.Err() != nil {
		t.Fatalf("Execute() error = %v", resp.Err())
	}
	if got := resp.Headers()["X-Test"]; len(got) != 1 || got[0] != "yes" {
		t.Fatalf("Headers()[X-Test] = %v, want [yes]", got)
	}
	cookies := resp.Cookies()
	if len(cookies) != 1 || cookies[0].Name != "sid" || cookies[0].Value != "abc" {
		t.Fatalf("Cookies() = %+v, want sid=abc", cookies)
	}
	if got := resp.ContentLength(); got != int64(len("hello")) {
		t.Fatalf("ContentLength() = %d, want %d", got, len("hello"))
	}
}

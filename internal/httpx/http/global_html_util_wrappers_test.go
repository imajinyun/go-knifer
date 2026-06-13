package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func TestAdditionalGlobalHTMLAndUtilWrappers(t *testing.T) {
	previous := SnapshotGlobalConfig()
	defer ConfigureGlobalConfig(previous)
	SetGlobalBoundary("boundary-extra")
	SetGlobalDecodeURL(true)
	AddGlobalHeader("X-Extra", "one")
	defer RemoveGlobalHeader("X-Extra")
	if GetGlobalBoundary() != "boundary-extra" || !IsGlobalDecodeURL() {
		t.Fatalf("global boundary/decode = %q/%v", GetGlobalBoundary(), IsGlobalDecodeURL())
	}
	if got := CloneGlobalHeaders()["X-Extra"]; len(got) != 1 || got[0] != "one" {
		t.Fatalf("CloneGlobalHeaders X-Extra = %v", got)
	}
	if got := CleanHTMLWithOptions("a[drop]b", WithHTMLTagRegexp(regexp.MustCompile(`\[.*?\]`)), WithHTMLCommentRegexp(regexp.MustCompile(`$^`))); got != "ab" {
		t.Fatalf("CleanHTMLWithOptions = %q", got)
	}
	if got := FilterHTMLTagWithOptions("<custom>drop</custom><p>keep</p>", []string{"custom"}, WithHTMLFilterCompileFunc(regexp.Compile)); got != "<p>keep</p>" {
		t.Fatalf("FilterHTMLTagWithOptions = %q", got)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write([]byte(r.Method + ":" + string(body)))
	}))
	defer srv.Close()
	if got, err := GetStringSafeE(srv.URL, WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})); err != nil || got != "GET:" {
		t.Fatalf("GetStringSafeE = %q, %v", got, err)
	}
	if got, err := GetWithTimeoutE(srv.URL, time.Second); err != nil || got != "GET:" {
		t.Fatalf("GetWithTimeoutE = %q, %v", got, err)
	}
	if got, err := PostStringE(srv.URL, "body"); err != nil || got != "POST:body" {
		t.Fatalf("PostStringE = %q, %v", got, err)
	}
	if got, err := PostStringSafeE(srv.URL, "safe", WithURLPolicy(URLPolicy{AllowedSchemes: []string{"http", "https"}, RejectPrivate: false})); err != nil || got != "POST:safe" {
		t.Fatalf("PostStringSafeE = %q, %v", got, err)
	}
}

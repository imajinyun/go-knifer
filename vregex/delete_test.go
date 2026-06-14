package vregex

import (
	"regexp"
	"strings"
	"testing"
)

func TestRegexDeleteFacade(t *testing.T) {
	if got := DelFirst(`\d+`, "a1b22"); got != "ab22" {
		t.Fatalf("DelFirst = %q", got)
	}
	if got := DelFirstRe(regexp.MustCompile(`\d+`), "a1b22"); got != "ab22" {
		t.Fatalf("DelFirstRe = %q", got)
	}
	if got := DelLast(`\d+`, "a1b22"); got != "a1b" {
		t.Fatalf("DelLast = %q", got)
	}
	if got := DelLastRe(regexp.MustCompile(`\d+`), "a1b22"); got != "a1b" {
		t.Fatalf("DelLastRe = %q", got)
	}
	if got := DelAll(`\d+`, "a1b22"); got != "ab" {
		t.Fatalf("DelAll = %q", got)
	}
	if got := DelAllRe(regexp.MustCompile(`\d+`), "a1b22"); got != "ab" {
		t.Fatalf("DelAllRe = %q", got)
	}
	if got := DelPre(`\d+`, "a1b22"); got != "b22" {
		t.Fatalf("DelPre = %q", got)
	}
	if got := DelPreRe(regexp.MustCompile(`\d+`), "a1b22"); got != "b22" {
		t.Fatalf("DelPreRe = %q", got)
	}
}

func TestRegexDeleteFacadeWithOptions(t *testing.T) {
	opt := WithCompileFunc(func(pattern string) (*regexp.Regexp, error) {
		return regexp.Compile(strings.ReplaceAll(pattern, "TOKEN", `\d+`))
	})

	if got := DelFirstWithOptions(`TOKEN`, "a123b456", opt); got != "ab456" {
		t.Fatalf("DelFirstWithOptions = %q", got)
	}
	if got := DelLastWithOptions(`TOKEN`, "a123b456", opt); got != "a123b" {
		t.Fatalf("DelLastWithOptions = %q", got)
	}
	if got := DelAllWithOptions(`TOKEN`, "a123b456", opt); got != "ab" {
		t.Fatalf("DelAllWithOptions = %q", got)
	}
	if got := DelPreWithOptions(`TOKEN`, "a123b", opt); got != "b" {
		t.Fatalf("DelPreWithOptions = %q", got)
	}
}

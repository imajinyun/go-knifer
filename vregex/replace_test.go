package vregex

import (
	"regexp"
	"strings"
	"testing"
)

func TestRegexReplaceFacade(t *testing.T) {
	if Replace(`\d`, "a1b2", "*") != "a*b*" || Replace(`(`, "x", "*") != "x" {
		t.Fatal("Replace failed")
	}
	if got := ReplaceFirst(`\d+`, "a123b456", "X"); got != "aXb456" {
		t.Fatalf("ReplaceFirst = %q", got)
	}
	if got := ReplaceFirstRe(regexp.MustCompile(`\d+`), "a123b456", "X"); got != "aXb456" {
		t.Fatalf("ReplaceFirstRe = %q", got)
	}
	if got := ReplaceAll("中文1234", `(\d+)`, `($1)`); got != "中文(1234)" {
		t.Fatalf("ReplaceAll = %q", got)
	}
	if got := ReplaceAllRe("中文1234", regexp.MustCompile(`(\d+)`), `($1)`); got != "中文(1234)" {
		t.Fatalf("ReplaceAllRe = %q", got)
	}
	if got := ReplaceAllFunc("a1b22", `\d+`, func(m MatchResult) string { return "[" + m.Text + "]" }); got != "a[1]b[22]" {
		t.Fatalf("ReplaceAllFunc = %q", got)
	}
	if got := ReplaceAllFuncRe("a1b22", regexp.MustCompile(`\d+`), func(m MatchResult) string { return "[" + m.Text + "]" }); got != "a[1]b[22]" {
		t.Fatalf("ReplaceAllFuncRe = %q", got)
	}
}

func TestRegexReplaceFacadeWithOptions(t *testing.T) {
	opt := WithCompileFunc(func(pattern string) (*regexp.Regexp, error) {
		return regexp.Compile(strings.ReplaceAll(pattern, "TOKEN", `\d+`))
	})

	if got := ReplaceWithOptions(`TOKEN`, "a123b456", "X", opt); got != "aXbX" {
		t.Fatalf("ReplaceWithOptions = %q", got)
	}
	if got := ReplaceFirstWithOptions(`TOKEN`, "a123b456", "X", opt); got != "aXb456" {
		t.Fatalf("ReplaceFirstWithOptions = %q", got)
	}
	if got := ReplaceAllWithOptions("x123", `x(TOKEN)`, "$1", opt); got != "123" {
		t.Fatalf("ReplaceAllWithOptions = %q", got)
	}
	if got := ReplaceAllFuncWithOptions("a1b22", `TOKEN`, func(m MatchResult) string { return "[" + m.Text + "]" }, opt); got != "a[1]b[22]" {
		t.Fatalf("ReplaceAllFuncWithOptions = %q", got)
	}
}

package vregex

import "testing"

func TestRegexEscapeFacade(t *testing.T) {
	if got := Escape("a+b(c)"); got != `a\+b\(c\)` {
		t.Fatalf("Escape = %q", got)
	}
	if got := EscapeChar('+'); got != `\+` {
		t.Fatalf("EscapeChar = %q", got)
	}
	if got := EscapeChar('a'); got != "a" {
		t.Fatalf("EscapeChar non-keyword = %q", got)
	}
}

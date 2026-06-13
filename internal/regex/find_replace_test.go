package regex

import (
	"reflect"
	"testing"
)

func TestFindCountIndexMatchReplaceAndEscape(t *testing.T) {
	if got := FindAllGroup1(`(\d+)`, "a1b22c333"); !reflect.DeepEqual(got, []string{"1", "22", "333"}) {
		t.Fatalf("FindAllGroup1 = %#v", got)
	}
	if got := Count(`\d+`, "a1b22c333"); got != 3 {
		t.Fatalf("Count = %d", got)
	}
	if !Contains(`\d+`, "abc1") || Contains(`\d+`, "abc") {
		t.Fatalf("Contains failed")
	}
	first := IndexOf(`\d+`, "ab12cd34")
	last := LastIndexOf(`\d+`, "ab12cd34")
	if first == nil || first.Start != 2 || first.End != 4 || last == nil || last.Start != 6 || last.End != 8 {
		t.Fatalf("IndexOf/LastIndexOf failed: %#v %#v", first, last)
	}
	if n, ok := GetFirstNumber("abc123def"); !ok || n != 123 {
		t.Fatalf("GetFirstNumber = %d %v", n, ok)
	}
	if !IsMatch(`\d+`, "123") || IsMatch(`\d+`, "a123") {
		t.Fatalf("IsMatch failed")
	}
	if got := ReplaceAll("中文1234", `(\d+)`, `($1)`); got != "中文(1234)" {
		t.Fatalf("ReplaceAll = %q", got)
	}
	if got := ReplaceAllFunc("a1b22", `\d+`, func(m MatchResult) string { return "[" + m.Text + "]" }); got != "a[1]b[22]" {
		t.Fatalf("ReplaceAllFunc = %q", got)
	}
	if got := Escape("a+b(c)"); got != `a\+b\(c\)` {
		t.Fatalf("Escape = %q", got)
	}
}

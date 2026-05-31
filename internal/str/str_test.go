package str

import "testing"

// Tests aligned with hutool-core CharSequenceUtilTest and StrUtilTest.

func TestIsEmptyAndBlank(t *testing.T) {
	if !IsEmpty("") || IsEmpty("a") {
		t.Fatalf("IsEmpty failed")
	}
	if !IsBlank("  \t\n") || IsBlank(" a ") {
		t.Fatalf("IsBlank failed")
	}
	if !IsAllEmpty("", "") || IsAllEmpty("", "a") {
		t.Fatalf("IsAllEmpty failed")
	}
	if !IsAllBlank("  ", "\t") || IsAllBlank(" ", "x") {
		t.Fatalf("IsAllBlank failed")
	}
	if !HasEmpty("a", "") || HasEmpty("a", "b") {
		t.Fatalf("HasEmpty failed")
	}
	if !HasBlank("a", " ") || HasBlank("a", "b") {
		t.Fatalf("HasBlank failed")
	}
}

func TestSubAndSlicing(t *testing.T) {
	if got := Sub("hutool", 0, 3); got != "hut" {
		t.Fatalf("Sub: %q", got)
	}
	if got := Sub("hutool", -3, -1); got != "oo" {
		t.Fatalf("Sub negative: %q", got)
	}
	if got := SubBefore("a.b.c", ".", false); got != "a" {
		t.Fatalf("SubBefore first: %q", got)
	}
	if got := SubBefore("a.b.c", ".", true); got != "a.b" {
		t.Fatalf("SubBefore last: %q", got)
	}
	if got := SubAfter("a.b.c", ".", false); got != "b.c" {
		t.Fatalf("SubAfter first: %q", got)
	}
	if got := SubAfter("a.b.c", ".", true); got != "c" {
		t.Fatalf("SubAfter last: %q", got)
	}
}

func TestSplitTrim(t *testing.T) {
	got := SplitTrim(" a , b ,, c ", ",")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: %v", got)
	}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("idx %d: %q != %q", i, got[i], v)
		}
	}
}

func TestPadAndRepeat(t *testing.T) {
	if got := PadLeft("12", 5, '0'); got != "00012" {
		t.Fatalf("PadLeft: %q", got)
	}
	if got := PadRight("12", 5, '0'); got != "12000" {
		t.Fatalf("PadRight: %q", got)
	}
	if got := Repeat("ab", 3); got != "ababab" {
		t.Fatalf("Repeat: %q", got)
	}
}

func TestReverseAndContains(t *testing.T) {
	if got := Reverse("中国"); got != "国中" {
		t.Fatalf("Reverse: %q", got)
	}
	if !ContainsIgnoreCase("HelloWorld", "world") {
		t.Fatalf("ContainsIgnoreCase failed")
	}
	if !ContainsAny("abc", "x", "b") {
		t.Fatalf("ContainsAny failed")
	}
	if !ContainsAll("abcde", "a", "c") || ContainsAll("abc", "a", "z") {
		t.Fatalf("ContainsAll failed")
	}
}

func TestStartsEndsWith(t *testing.T) {
	if !StartsWith("hutool", "hu") || StartsWith("hutool", "tool") {
		t.Fatalf("StartsWith failed")
	}
	if !EndsWith("hutool", "ool") || EndsWith("hutool", "hu") {
		t.Fatalf("EndsWith failed")
	}
	if !EqualsIgnoreCase("HUTOOL", "hutool") {
		t.Fatalf("EqualsIgnoreCase failed")
	}
}

func TestFormat(t *testing.T) {
	if got := Format("name={}, age={}", "tom", 12); got != "name=tom, age=12" {
		t.Fatalf("Format: %q", got)
	}
	// Escaping.
	if got := Format("\\{}={}", "x"); got != "{}=x" {
		t.Fatalf("Format escape: %q", got)
	}
	// More placeholders than arguments.
	if got := Format("a={},b={}", 1); got != "a=1,b={}" {
		t.Fatalf("Format extra: %q", got)
	}
}

func TestPrefixSuffix(t *testing.T) {
	if got := RemovePrefix("foo_bar", "foo_"); got != "bar" {
		t.Fatalf("RemovePrefix: %q", got)
	}
	if got := RemoveSuffix("foo.txt", ".txt"); got != "foo" {
		t.Fatalf("RemoveSuffix: %q", got)
	}
	if got := AddPrefixIfNot("bar", "foo_"); got != "foo_bar" {
		t.Fatalf("AddPrefixIfNot: %q", got)
	}
	if got := AddPrefixIfNot("foo_bar", "foo_"); got != "foo_bar" {
		t.Fatalf("AddPrefixIfNot already: %q", got)
	}
	if got := AddSuffixIfNot("foo", ".txt"); got != "foo.txt" {
		t.Fatalf("AddSuffixIfNot: %q", got)
	}
}

func TestLength(t *testing.T) {
	if Length("中国a") != 3 || RuneLen("你好") != 2 {
		t.Fatalf("Length rune count failed")
	}
}

func TestEmojiHelpers(t *testing.T) {
	if !ContainsEmoji("hi😀") || !ContainsEmoji("go❤️") || !ContainsEmoji("1️⃣") {
		t.Fatal("ContainsEmoji() = false")
	}
	if got := RemoveEmoji("hi😀 go❤️ 1️⃣"); got != "hi go " {
		t.Fatalf("RemoveEmoji() = %q", got)
	}
}

package str

import "testing"

// Tests cover the utility toolkit-core CharSequenceUtilTest and StrUtilTest.

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
	if got := Sub("foobar", 0, 3); got != "foo" {
		t.Fatalf("Sub: %q", got)
	}
	if got := Sub("foobar", -3, -1); got != "ba" {
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
	if !StartsWith("foobar", "foo") || StartsWith("foobar", "bar") {
		t.Fatalf("StartsWith failed")
	}
	if !EndsWith("foobar", "bar") || EndsWith("foobar", "foo") {
		t.Fatalf("EndsWith failed")
	}
	if !EqualsIgnoreCase("FOOBAR", "foobar") {
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

func TestUnicodeEscapeHelpers(t *testing.T) {
	escaped := EscapeUnicode("Go 中国 😀")
	if escaped != `Go \u4E2D\u56FD \uD83D\uDE00` {
		t.Fatalf("EscapeUnicode() = %q", escaped)
	}
	if got := UnescapeUnicode(escaped); got != "Go 中国 😀" {
		t.Fatalf("UnescapeUnicode() = %q", got)
	}
	if got := UnescapeUnicode(`bad \u12GZ keep`); got != `bad \u12GZ keep` {
		t.Fatalf("UnescapeUnicode malformed = %q", got)
	}
	if got := UnescapeUnicode(`\uD83Dbroken`); got != string(rune(0xD83D))+"broken" {
		t.Fatalf("UnescapeUnicode lone surrogate = %q", got)
	}
}

func TestAntPathMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{name: "double star crosses segments", pattern: "/api/**/users/?", path: "/api/v1/admin/users/a", want: true},
		{name: "single star stays in segment", pattern: "/api/*/users", path: "/api/v1/admin/users", want: false},
		{name: "single char wildcard", pattern: "/file-?.txt", path: "/file-a.txt", want: true},
		{name: "literal mismatch", pattern: "/file-?.txt", path: "/file-ab.txt", want: false},
		{name: "double star zero segment", pattern: "/api/**", path: "/api", want: true},
		{name: "empty pattern", pattern: "", path: "", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AntPathMatch(tt.pattern, tt.path); got != tt.want {
				t.Fatalf("AntPathMatch(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
	if !AntPathMatchWithSeparator("foo.**.bar", "foo.a.b.bar", ".") {
		t.Fatal("AntPathMatchWithSeparator custom separator failed")
	}
	if AntPathMatchWithSeparator("foo.*.bar", "foo.a.b.bar", ".") {
		t.Fatal("single star should not cross custom separator segments")
	}
}

func TestTextSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want float64
	}{
		{name: "both empty", a: "", b: " \t", want: 1},
		{name: "one empty", a: "abc", b: " ", want: 0},
		{name: "same runes", a: "go go", b: "og", want: 1},
		{name: "partial overlap", a: "abc", b: "bcd", want: 0.5},
		{name: "unicode", a: "你好世界", b: "你好", want: 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JaccardSimilarity(tt.a, tt.b); got != tt.want {
				t.Fatalf("JaccardSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNGramSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		n    int
		want float64
	}{
		{name: "invalid n", a: "abc", b: "abc", n: 0, want: 0},
		{name: "both empty", a: "", b: " ", n: 2, want: 1},
		{name: "same short text", a: "go", b: "go", n: 3, want: 1},
		{name: "partial bigram", a: "abcd", b: "abef", n: 2, want: 0.2},
		{name: "unicode bigram", a: "你好世界", b: "你好呀", n: 2, want: 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NGramSimilarity(tt.a, tt.b, tt.n); got != tt.want {
				t.Fatalf("NGramSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimHashAndHammingDistance64(t *testing.T) {
	if SimHash("") != 0 {
		t.Fatal("SimHash(empty) != 0")
	}
	a := SimHash("go knifer toolkit")
	b := SimHash("go knifer toolkit")
	c := SimHash("database security password")
	if a != b {
		t.Fatal("SimHash() should be deterministic")
	}
	if HammingDistance64(a, b) != 0 {
		t.Fatal("HammingDistance64(same, same) != 0")
	}
	if HammingDistance64(0, ^uint64(0)) != 64 {
		t.Fatal("HammingDistance64(0, all bits) != 64")
	}
	if HammingDistance64(a, c) == 0 {
		t.Fatal("different text should not produce identical SimHash in this fixture")
	}
}

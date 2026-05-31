package version

import "testing"

func TestCompareVersion(t *testing.T) {
	cases := []struct {
		left  string
		right string
		want  int
	}{
		{"", "v1", -1},
		{"v1", "v1", 0},
		{"v1", "", 1},
		{"1.0.0", "1.0.2", -1},
		{"1.0.2", "1.0.2a", -1},
		{"1.0.3", "1.0.2a", 1},
		{"1.0a", "1.0.1", 1},
		{"1.0.1", "1.0a", -1},
		{"1.13.0", "1.12.1c", 1},
		{"V0.0.20170102", "V0.0.20170101", 1},
		{"1.0.0", "1.0", 0},
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0+2", "1.0.0+1", 1},
	}
	for _, tc := range cases {
		got := sign(CompareVersion(tc.left, tc.right))
		if got != tc.want {
			t.Fatalf("CompareVersion(%q, %q)=%d, want %d", tc.left, tc.right, got, tc.want)
		}
	}
}

func TestComparisonHelpers(t *testing.T) {
	if !IsGreaterThan("1.0.3", "1.0.2") {
		t.Fatal("IsGreaterThan failed")
	}
	if !IsGreaterThanOrEqual("1.0.2", "1.0.2") {
		t.Fatal("IsGreaterThanOrEqual failed")
	}
	if !IsLessThan("1.0.1", "1.0.2") {
		t.Fatal("IsLessThan failed")
	}
	if !IsLessThanOrEqual("1.0.2", "1.0.2") {
		t.Fatal("IsLessThanOrEqual failed")
	}
	if IsGreaterThan("1.0.1", "1.0.2") || IsLessThan("1.0.3", "1.0.2") {
		t.Fatal("comparison helpers returned unexpected true")
	}
}

func TestMatchEl(t *testing.T) {
	cases := []struct {
		current string
		expr    string
		want    bool
	}{
		{"1.0.2", ">=1.0.2", true},
		{"1.0.2", "≥1.0.2", true},
		{"1.0.2", "<1.0.1;1.0.2", true},
		{"1.0.2", "<1.0.2", false},
		{"1.0.2", "<=1.0.2", true},
		{"1.0.2", "≤1.0.2", true},
		{"1.0.2", "1.0.0-1.1.1", true},
		{"1.0.2", "1.0.3-1.1.1", false},
		{"1.0.2", "-1.0.2", true},
		{"1.0.2", "1.0.2-", true},
		{"1.0.2", "", false},
	}
	for _, tc := range cases {
		if got := MatchEl(tc.current, tc.expr); got != tc.want {
			t.Fatalf("MatchEl(%q, %q)=%v, want %v", tc.current, tc.expr, got, tc.want)
		}
	}
}

func TestMatchElWithDelimiterAndAnyMatch(t *testing.T) {
	if !MatchElByDelimiter("1.0.2", "<1.0.1,1.0.2", ",") {
		t.Fatal("comma separated exact match failed")
	}
	if !MatchElByDelimiter("1.0.2", "1.0.1,1.0.2-1.1.1", ",") {
		t.Fatal("comma separated range match failed")
	}
	if MatchElWithDelimiterErr("1.0.2", ">1.0.0", "-") == nil {
		t.Fatal("expected invalid delimiter error")
	}
	if !AnyMatch("1.0.2", "<1.0.1", "1.0.2") {
		t.Fatal("AnyMatch failed")
	}
	if !AnyMatchSlice("1.0.2", []string{"<1.0.1", "1.0.2"}) {
		t.Fatal("AnyMatchSlice failed")
	}
}

func TestNullComparisonExpression(t *testing.T) {
	if !MatchEl("1.0.0", ">null") {
		t.Fatal("non-empty version should be greater than null expression")
	}
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	if n > 0 {
		return 1
	}
	return 0
}

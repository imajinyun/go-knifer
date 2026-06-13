package captcha

import (
	"strconv"
	"strings"
	"testing"
)

func TestRandomGenerator_GenerateLength(t *testing.T) {
	g := NewRandomGenerator(6)
	for i := 0; i < 50; i++ {
		s := g.Generate()
		if len(s) != 6 {
			t.Fatalf("len=%d, want 6, got %q", len(s), s)
		}
	}
}

func TestRandomGenerator_VerifyIgnoreCase(t *testing.T) {
	g := NewRandomGenerator(4)
	if !g.Verify("AbCd", "abcd") {
		t.Fatalf("expect ignore case match")
	}
	if g.Verify("ab", "  ") {
		t.Fatalf("blank input should fail")
	}
}

func TestRandomGenerator_CustomBase(t *testing.T) {
	g := NewRandomGeneratorWithBase("01", 8)
	s := g.Generate()
	for _, c := range s {
		if c != '0' && c != '1' {
			t.Fatalf("unexpected char %q in %q", c, s)
		}
	}
}

func TestRandomGenerator_GenWithOptions(t *testing.T) {
	g := NewRandomGeneratorWithBase("abcd", 4)
	idx := 0
	code := g.GenWithOptions(WithGeneratorRandomInt(func(max int) int {
		v := idx
		idx++
		return v % max
	}))
	if code != "abcd" {
		t.Fatalf("GenWithOptions code = %q, want abcd", code)
	}
}

func TestMathGenerator_GenerateAndVerify(t *testing.T) {
	g := NewMathGenerator()
	for i := 0; i < 50; i++ {
		code := g.Generate()
		if !strings.HasSuffix(code, "=") {
			t.Fatalf("code should end with '=': %q", code)
		}
		v, ok := evalMathExpr(code)
		if !ok {
			t.Fatalf("expr eval failed: %q", code)
		}
		if !g.Verify(code, strconv.Itoa(v)) {
			t.Fatalf("verify failed: code=%q want=%d", code, v)
		}
	}
}

func TestMathGenerator_GenWithOptions(t *testing.T) {
	g := NewMathGeneratorWith(1, false)
	values := []int{1, 7, 3}
	idx := 0
	code := g.GenWithOptions(WithGeneratorRandomInt(func(max int) int {
		v := values[idx]
		idx++
		return v % max
	}))
	if code != "7-3=" {
		t.Fatalf("GenWithOptions code = %q, want 7-3=", code)
	}
	if !g.Verify(code, "4") {
		t.Fatalf("GenWithOptions code should verify: %q", code)
	}
}

func TestMathGenerator_VerifyWithOptionsUsesParser(t *testing.T) {
	g := NewMathGenerator()
	parseCalls := 0
	if !g.VerifyWithOptions("left+right=", "sum", WithGeneratorIntParser(func(text string) (int, error) {
		parseCalls++
		switch text {
		case "left":
			return 2, nil
		case "right":
			return 3, nil
		case "sum":
			return 5, nil
		default:
			return strconv.Atoi(text)
		}
	})) {
		t.Fatal("VerifyWithOptions should use custom parser")
	}
	if parseCalls != 3 {
		t.Fatalf("parser calls = %d, want 3", parseCalls)
	}
}

func TestMathGenerator_NoNegative(t *testing.T) {
	g := NewMathGeneratorWith(2, false)
	for i := 0; i < 100; i++ {
		code := g.Generate()
		v, ok := evalMathExpr(code)
		if !ok {
			t.Fatalf("eval fail: %q", code)
		}
		// Check only when the operator is '-'.
		if strings.Contains(code, "-") && v < 0 {
			t.Fatalf("expected non-negative, got %d for %q", v, code)
		}
	}
}

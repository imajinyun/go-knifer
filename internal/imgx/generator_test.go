package imgx

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestRandomGenerator_GenerateLength(t *testing.T) {
	g := NewRandomGenerator(6)
	for i := 0; i < 50; i++ {
		s := g.Gen()
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
	s := g.Gen()
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
		code := g.Gen()
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
		code := g.Gen()
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

func TestRandomGeneratorDefaultsAndNilOptions(t *testing.T) {
	g := NewRandomGeneratorWithBase("", 0)
	idx := 0
	code := g.GenWithOptions(nil, WithGeneratorRandomInt(nil), WithGeneratorRandomInt(func(max int) int {
		idx++
		return max - 1
	}))
	if len(code) != 4 {
		t.Fatalf("default length code = %q, want len 4", code)
	}
	if idx != 4 {
		t.Fatalf("random calls = %d, want 4", idx)
	}
}

func TestGeneratorHelpersEdgeCases(t *testing.T) {
	if got := normalizeRandomIndex(-7, 5); got != 2 {
		t.Fatalf("normalizeRandomIndex(-7, 5) = %d, want 2", got)
	}
	if got := normalizeRandomIndex(7, 0); got != 0 {
		t.Fatalf("normalizeRandomIndex(7, 0) = %d, want 0", got)
	}
	if got := padRight("abcd", 2, '.'); got != "abcd" {
		t.Fatalf("padRight long input = %q, want abcd", got)
	}
	if got := NewMathGeneratorWith(0, true); got.NumberLength != 2 {
		t.Fatalf("NewMathGeneratorWith default length = %d, want 2", got.NumberLength)
	}
	if got := NewMathGeneratorWith(3, true).Length(); got != 8 {
		t.Fatalf("MathGenerator.Length = %d, want 8", got)
	}
}

func TestMathGeneratorVerifyRejectsInvalidInputs(t *testing.T) {
	g := NewMathGenerator()
	parserErr := errors.New("parse failed")
	if g.VerifyWithOptions("1+2=", "3", WithGeneratorIntParser(func(string) (int, error) {
		return 0, parserErr
	})) {
		t.Fatal("VerifyWithOptions should reject parser errors")
	}
	if g.VerifyWithOptions("not an expression", "3") {
		t.Fatal("VerifyWithOptions should reject invalid expressions")
	}
	if _, ok := evalMathExprWithOptions("left+right=", WithGeneratorIntParser(func(string) (int, error) {
		return 0, parserErr
	})); ok {
		t.Fatal("evalMathExprWithOptions should reject operand parser errors")
	}
	if _, ok := evalMathExpr("1/2="); ok {
		t.Fatal("evalMathExpr should reject unsupported operators")
	}
}

func TestMathGeneratorNonNegativeSubtractionWithZeroLeftOperand(t *testing.T) {
	g := NewMathGeneratorWith(1, false)
	values := []int{1, 0}
	idx := 0
	code := g.GenWithOptions(WithGeneratorRandomInt(func(max int) int {
		v := values[idx]
		idx++
		return v % max
	}))
	if code != "0-0=" {
		t.Fatalf("GenWithOptions zero subtraction = %q, want 0-0=", code)
	}
}

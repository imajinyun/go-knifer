package vcaptcha_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vcaptcha"
)

func TestFacadeRandomGenerator(t *testing.T) {
	g := vcaptcha.NewRandomGenerator(4)
	code := g.Generate()
	if len(code) != 4 {
		t.Fatalf("expected code length 4, got %d", len(code))
	}
	if !g.Verify(code, code) {
		t.Fatal("expected generated code to verify")
	}
	if g.Verify(code, "wrong") {
		t.Fatal("expected wrong code to fail verification")
	}
}

func TestFacadeMathGenerator(t *testing.T) {
	g := vcaptcha.NewMathGenerator()
	code := g.Generate()
	if len(code) == 0 {
		t.Fatal("expected non-empty math code")
	}
	// MathGenerator produces expressions like "1+2="; Verify needs the computed answer.
	// We just smoke-test that generation and verification accept a correct answer.
	if !g.Verify("1+1=", "2") {
		t.Fatal("expected 1+1= to verify with answer 2")
	}
}

func TestFacadeVerifyIgnoreCase(t *testing.T) {
	if !vcaptcha.VerifyCaptchaIgnoreCase("ABC", "abc") {
		t.Fatal("expected case-insensitive verification to pass")
	}
	if vcaptcha.VerifyCaptchaIgnoreCase("ABC", "def") {
		t.Fatal("expected different code to fail verification")
	}
}

func TestFacadeLineCaptcha(t *testing.T) {
	c := vcaptcha.NewLineCaptcha(100, 40)
	if c == nil {
		t.Fatal("expected non-nil line captcha")
	}
	c.CreateCode()
	code := c.Code()
	if len(code) == 0 {
		t.Fatal("expected line captcha to have non-empty code")
	}
}

func TestFacadeCircleCaptcha(t *testing.T) {
	c := vcaptcha.NewCircleCaptcha(100, 40)
	if c == nil {
		t.Fatal("expected non-nil circle captcha")
	}
	c.CreateCode()
	code := c.Code()
	if len(code) == 0 {
		t.Fatal("expected circle captcha to have non-empty code")
	}
}

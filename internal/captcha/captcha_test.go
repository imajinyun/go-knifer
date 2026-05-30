package captcha

import (
	"bytes"
	"image/gif"
	"image/png"
	"strconv"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Generator tests.
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// Captcha implementation tests.
// ---------------------------------------------------------------------------

func TestLineCaptcha(t *testing.T) {
	c := NewLineCaptcha(200, 80)
	c.CreateCode()
	if c.Code() == "" {
		t.Fatal("empty code")
	}
	if !c.Verify(c.Code()) {
		t.Fatal("self verify failed")
	}
	if _, err := png.Decode(bytes.NewReader(c.ImageBytes())); err != nil {
		t.Fatalf("image not valid PNG: %v", err)
	}
}

func TestCircleCaptcha(t *testing.T) {
	c := NewCircleCaptchaWith(200, 80, 4, 8)
	c.CreateCode()
	if len(c.Code()) != 4 {
		t.Fatalf("code length = %d, want 4", len(c.Code()))
	}
	if _, err := png.Decode(bytes.NewReader(c.ImageBytes())); err != nil {
		t.Fatalf("invalid PNG: %v", err)
	}
}

func TestShearCaptcha(t *testing.T) {
	c := NewShearCaptchaWith(200, 80, 5, 4)
	c.CreateCode()
	if c.Code() == "" {
		t.Fatal("empty code")
	}
	if _, err := png.Decode(bytes.NewReader(c.ImageBytes())); err != nil {
		t.Fatalf("invalid PNG: %v", err)
	}
}

func TestGifCaptcha(t *testing.T) {
	c := NewGifCaptcha(200, 80)
	c.CreateCode()
	if c.Code() == "" {
		t.Fatal("empty code")
	}
	g, err := gif.DecodeAll(bytes.NewReader(c.ImageBytes()))
	if err != nil {
		t.Fatalf("invalid GIF: %v", err)
	}
	if len(g.Image) != len(c.Code()) {
		t.Fatalf("frames=%d, want %d", len(g.Image), len(c.Code()))
	}
}

func TestCaptchaUtilFactories(t *testing.T) {
	if CreateLineCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
	if CreateCircleCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
	if CreateShearCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
	if CreateGifCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
}

func TestICaptchaInterface(t *testing.T) {
	var _ ICaptcha = NewLineCaptcha(100, 40)
	var _ ICaptcha = NewCircleCaptcha(100, 40)
	var _ ICaptcha = NewShearCaptcha(100, 40)
	var _ ICaptcha = NewGifCaptcha(100, 40)
}

func TestImageBase64Data(t *testing.T) {
	c := NewLineCaptcha(100, 40)
	s := c.ImageBase64Data()
	if !strings.HasPrefix(s, "data:image/png;base64,") {
		t.Fatalf("unexpected data uri prefix: %q", s[:30])
	}
}

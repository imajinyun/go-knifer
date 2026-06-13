package captcha

import (
	"bytes"
	"image/png"
	"testing"
)

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

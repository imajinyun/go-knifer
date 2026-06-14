package imgx

import (
	"bytes"
	"image/png"
	"testing"
)

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

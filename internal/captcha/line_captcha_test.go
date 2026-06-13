package captcha

import (
	"bytes"
	"image/png"
	"testing"
)

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

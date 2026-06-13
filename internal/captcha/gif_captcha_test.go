package captcha

import (
	"bytes"
	"image/gif"
	"testing"
)

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

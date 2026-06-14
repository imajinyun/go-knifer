package vimg_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vimg"
)

func TestFacadeLineCaptcha(t *testing.T) {
	c := vimg.NewLineCaptcha(100, 40)
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
	c := vimg.NewCircleCaptcha(100, 40)
	if c == nil {
		t.Fatal("expected non-nil circle captcha")
	}
	c.CreateCode()
	code := c.Code()
	if len(code) == 0 {
		t.Fatal("expected circle captcha to have non-empty code")
	}
}

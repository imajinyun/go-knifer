package captcha

import "testing"

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

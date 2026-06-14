package vcaptcha_test

import (
	"image/color"
	"testing"

	"github.com/imajinyun/go-knifer/vcaptcha"
)

func TestFacadeCaptchaOptions(t *testing.T) {
	colorCalls := 0
	line := vcaptcha.NewLineCaptchaWithOptions(100, 40,
		vcaptcha.WithGenerator(fixedGenerator{code: "ABCD"}),
		vcaptcha.WithBackground(color.Black),
		vcaptcha.WithInterfereCount(0),
		vcaptcha.WithRandomInt(func(max int) int { return 0 }),
		vcaptcha.WithColorFunc(func() color.Color {
			colorCalls++
			return color.RGBA{R: 1, G: 2, B: 3, A: 255}
		}),
	)
	if got := line.Code(); got != "ABCD" {
		t.Fatalf("line captcha code = %q, want ABCD", got)
	}
	if !line.Verify("ABCD") {
		t.Fatal("line captcha should verify fixed code")
	}
	if colorCalls != len("ABCD") {
		t.Fatalf("custom color func calls=%d, want %d", colorCalls, len("ABCD"))
	}

	circle := vcaptcha.NewCircleCaptchaWithOptions(100, 40, vcaptcha.WithGenerator(fixedGenerator{code: "WXYZ"}))
	if got := circle.Code(); got != "WXYZ" {
		t.Fatalf("circle captcha code = %q, want WXYZ", got)
	}
	shear := vcaptcha.NewShearCaptchaWithOptions(100, 40, vcaptcha.WithGenerator(fixedGenerator{code: "EFGH"}))
	if got := shear.Code(); got != "EFGH" {
		t.Fatalf("shear captcha code = %q, want EFGH", got)
	}
	gif := vcaptcha.NewGifCaptchaWithOptions(100, 40, vcaptcha.WithGenerator(fixedGenerator{code: "IJKL"}), vcaptcha.WithGIFRepeat(1), vcaptcha.WithGIFDelay(5))
	if got := gif.Code(); got != "IJKL" {
		t.Fatalf("gif captcha code = %q, want IJKL", got)
	}
	if gif.Repeat != 1 || gif.Delay != 5 {
		t.Fatalf("gif options not applied: repeat=%d delay=%d", gif.Repeat, gif.Delay)
	}
}

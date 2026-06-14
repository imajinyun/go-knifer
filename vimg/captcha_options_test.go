package vimg_test

import (
	"image/color"
	"testing"

	"github.com/imajinyun/go-knifer/vimg"
)

func TestFacadeCaptchaOptions(t *testing.T) {
	colorCalls := 0
	line := vimg.NewLineCaptchaWithOptions(100, 40,
		vimg.WithGenerator(fixedGenerator{code: "ABCD"}),
		vimg.WithBackground(color.Black),
		vimg.WithInterfereCount(0),
		vimg.WithRandomInt(func(max int) int { return 0 }),
		vimg.WithColorFunc(func() color.Color {
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

	circle := vimg.NewCircleCaptchaWithOptions(100, 40, vimg.WithGenerator(fixedGenerator{code: "WXYZ"}))
	if got := circle.Code(); got != "WXYZ" {
		t.Fatalf("circle captcha code = %q, want WXYZ", got)
	}
	shear := vimg.NewShearCaptchaWithOptions(100, 40, vimg.WithGenerator(fixedGenerator{code: "EFGH"}))
	if got := shear.Code(); got != "EFGH" {
		t.Fatalf("shear captcha code = %q, want EFGH", got)
	}
	gif := vimg.NewGifCaptchaWithOptions(100, 40, vimg.WithGenerator(fixedGenerator{code: "IJKL"}), vimg.WithGIFRepeat(1), vimg.WithGIFDelay(5))
	if got := gif.Code(); got != "IJKL" {
		t.Fatalf("gif captcha code = %q, want IJKL", got)
	}
	if gif.Repeat != 1 || gif.Delay != 5 {
		t.Fatalf("gif options not applied: repeat=%d delay=%d", gif.Repeat, gif.Delay)
	}
}

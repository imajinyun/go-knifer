package captcha

import (
	"bytes"
	"image"
	"image/png"

	baseutil "github.com/imajinyun/go-knifer/internal/base"
)

// LineCaptcha 对应 hutool LineCaptcha：使用干扰线方式生成的图形验证码。
type LineCaptcha struct {
	AbstractCaptcha
}

// NewLineCaptcha 默认 5 位字符，150 条干扰线。
func NewLineCaptcha(width, height int) *LineCaptcha {
	return NewLineCaptchaWith(width, height, 5, 150)
}

// NewLineCaptchaWith 自定义字符数与干扰线条数。
func NewLineCaptchaWith(width, height, codeCount, lineCount int) *LineCaptcha {
	c := &LineCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = lineCount
	c.SetGenerator(NewRandomGenerator(codeCount))
	return c
}

// CreateCode 生成新的验证码字符串与图片。
func (c *LineCaptcha) CreateCode() {
	c.generateCode()
	c.setImageBytes(c.renderPNG(c.code))
}

// renderPNG 根据 code 渲染 PNG 字节。
func (c *LineCaptcha) renderPNG(code string) []byte {
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	fillBackground(img, c.bg())
	// 干扰线
	for i := 0; i < c.InterfereCount; i++ {
		xs := baseutil.RandomInt(c.Width)
		ys := baseutil.RandomInt(c.Height)
		xe := xs + baseutil.RandomInt(atLeastOne(c.Width/8))
		ye := ys + baseutil.RandomInt(atLeastOne(c.Height/8))
		drawLine(img, xs, ys, xe, ye, randomColor())
	}
	// 字符
	scale := computeScale(c.Height)
	drawString(img, code, c.Width, c.Height, scale)

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

// ImageBytes 重写以惰性生成。
func (c *LineCaptcha) ImageBytes() []byte {
	if c.imgBytes == nil {
		c.CreateCode()
	}
	return c.imgBytes
}

// Code 重写以惰性生成。
func (c *LineCaptcha) Code() string {
	if c.code == "" {
		c.CreateCode()
	}
	return c.code
}

// 让 ICaptcha 接口可用：通过组合自动获得 Verify/Write/WriteToFile 等方法。
var _ ICaptcha = (*LineCaptcha)(nil)

func atLeastOne(v int) int {
	if v > 1 {
		return v
	}
	return 1
}

// computeScale 根据图像高度计算字体放大倍数（5x7 字体）。
func computeScale(height int) int {
	scale := height / (fontHeight + 4)
	if scale < 1 {
		scale = 1
	}
	return scale
}

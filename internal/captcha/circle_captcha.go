package captcha

import (
	"bytes"
	"image"
	"image/png"

	baseutil "github.com/imajinyun/go-knifer/internal/base"
)

// CircleCaptcha 对应 hutool CircleCaptcha：使用干扰圆圈生成的图形验证码。
type CircleCaptcha struct {
	AbstractCaptcha
}

// NewCircleCaptcha 默认 5 位字符，15 个圆圈。
func NewCircleCaptcha(width, height int) *CircleCaptcha {
	return NewCircleCaptchaWith(width, height, 5, 15)
}

// NewCircleCaptchaWith 自定义字符数与干扰圆圈数。
func NewCircleCaptchaWith(width, height, codeCount, circleCount int) *CircleCaptcha {
	c := &CircleCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = circleCount
	c.SetGenerator(NewRandomGenerator(codeCount))
	return c
}

// CreateCode 生成新的验证码字符串与图片。
func (c *CircleCaptcha) CreateCode() {
	c.generateCode()
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	fillBackground(img, c.bg())
	half := c.Height >> 1
	for i := 0; i < c.InterfereCount; i++ {
		cx := baseutil.RandomInt(c.Width)
		cy := baseutil.RandomInt(c.Height)
		rx := baseutil.RandomInt(atLeastOne(half))
		ry := baseutil.RandomInt(atLeastOne(half))
		drawOval(img, cx, cy, rx, ry, randomColor())
	}
	drawString(img, c.code, c.Width, c.Height, computeScale(c.Height))

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	c.setImageBytes(buf.Bytes())
}

// ImageBytes 惰性渲染。
func (c *CircleCaptcha) ImageBytes() []byte {
	if c.imgBytes == nil {
		c.CreateCode()
	}
	return c.imgBytes
}

// Code 惰性生成。
func (c *CircleCaptcha) Code() string {
	if c.code == "" {
		c.CreateCode()
	}
	return c.code
}

var _ ICaptcha = (*CircleCaptcha)(nil)

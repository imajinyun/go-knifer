package captcha

import (
	"bytes"
	"image"
	"image/png"

	baseutil "github.com/imajinyun/go-knifer/internal/base"
)

// ShearCaptcha 对应 hutool ShearCaptcha：扭曲干扰验证码。
type ShearCaptcha struct {
	AbstractCaptcha
}

// NewShearCaptcha 默认 5 位字符，干扰线宽 4。
func NewShearCaptcha(width, height int) *ShearCaptcha {
	return NewShearCaptchaWith(width, height, 5, 4)
}

// NewShearCaptchaWith 自定义字符数与干扰线宽度。
func NewShearCaptchaWith(width, height, codeCount, thickness int) *ShearCaptcha {
	c := &ShearCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = thickness
	c.SetGenerator(NewRandomGenerator(codeCount))
	return c
}

// CreateCode 生成验证码字符串与图片。
func (c *ShearCaptcha) CreateCode() {
	c.generateCode()
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	fillBackground(img, c.bg())

	// 1) 字符
	drawString(img, c.code, c.Width, c.Height, computeScale(c.Height))

	// 2) 扭曲
	shearX(img, c.bg())
	shearY(img, c.bg())

	// 3) 干扰线（粗线）
	thickness := c.InterfereCount
	if thickness <= 0 {
		thickness = 4
	}
	x1 := 0
	y1 := baseutil.RandomInt(c.Height) + 1
	x2 := c.Width
	y2 := baseutil.RandomInt(c.Height) + 1
	drawThickLine(img, x1, y1, x2, y2, thickness, randomColor())

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	c.setImageBytes(buf.Bytes())
}

// ImageBytes 惰性渲染。
func (c *ShearCaptcha) ImageBytes() []byte {
	if c.imgBytes == nil {
		c.CreateCode()
	}
	return c.imgBytes
}

// Code 惰性生成。
func (c *ShearCaptcha) Code() string {
	if c.code == "" {
		c.CreateCode()
	}
	return c.code
}

var _ ICaptcha = (*ShearCaptcha)(nil)

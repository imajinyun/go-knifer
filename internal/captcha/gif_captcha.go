package captcha

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"

	baseutil "github.com/imajinyun/go-knifer/internal/base"
)

// GifCaptcha 对应 hutool GifCaptcha：GIF 动图验证码。
//
// 每帧仅高亮一位字符，其它字符以浅色绘制，制造闪烁效果。
type GifCaptcha struct {
	AbstractCaptcha

	// Repeat 帧循环次数；0 表示无限循环（与 hutool 一致）。
	Repeat int
	// Delay 单帧延迟，单位：1/100 秒（GIF 标准），默认 10。
	Delay int
}

// NewGifCaptcha 默认 5 位验证码，10 个干扰元素。
func NewGifCaptcha(width, height int) *GifCaptcha {
	return NewGifCaptchaWith(width, height, 5)
}

// NewGifCaptchaWith 自定义字符数。
func NewGifCaptchaWith(width, height, codeCount int) *GifCaptcha {
	c := &GifCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = 10
	c.Repeat = 0
	c.Delay = 10
	c.SetGenerator(NewRandomGenerator(codeCount))
	return c
}

// CreateCode 生成新的验证码字符串与 GIF 图片。
func (c *GifCaptcha) CreateCode() {
	c.generateCode()

	frames := make([]*image.Paletted, 0, len(c.code))
	delays := make([]int, 0, len(c.code))
	disposals := make([]byte, 0, len(c.code))

	pal := append(color.Palette{}, palette.Plan9...)
	for hi := 0; hi < len(c.code); hi++ {
		rgba := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
		fillBackground(rgba, c.bg())
		// 干扰圆圈
		half := c.Height >> 1
		for i := 0; i < c.InterfereCount; i++ {
			cx := baseutil.RandomInt(c.Width)
			cy := baseutil.RandomInt(c.Height)
			rx := baseutil.RandomInt(maxInt(half, 1))
			ry := baseutil.RandomInt(maxInt(half, 1))
			drawOval(rgba, cx, cy, rx, ry, randomColor())
		}
		// 字符：高亮当前 hi 位置
		drawCodeFrame(rgba, c.code, hi, c.Width, c.Height)

		// 转换为 paletted
		p := image.NewPaletted(rgba.Bounds(), pal)
		for y := 0; y < c.Height; y++ {
			for x := 0; x < c.Width; x++ {
				p.Set(x, y, rgba.RGBAAt(x, y))
			}
		}
		frames = append(frames, p)
		delays = append(delays, c.Delay)
		disposals = append(disposals, gif.DisposalBackground)
	}

	var buf bytes.Buffer
	_ = gif.EncodeAll(&buf, &gif.GIF{
		Image:     frames,
		Delay:     delays,
		LoopCount: c.Repeat,
		Disposal:  disposals,
	})
	c.setImageBytes(buf.Bytes())
}

// ImageBytes 惰性渲染。
func (c *GifCaptcha) ImageBytes() []byte {
	if c.imgBytes == nil {
		c.CreateCode()
	}
	return c.imgBytes
}

// Code 惰性生成。
func (c *GifCaptcha) Code() string {
	if c.code == "" {
		c.CreateCode()
	}
	return c.code
}

// ImageBase64Data 重写 GIF 类型。
func (c *GifCaptcha) ImageBase64Data() string {
	if c.imgBytes == nil {
		c.CreateCode()
	}
	return "data:image/gif;base64," + c.ImageBase64()
}

var _ ICaptcha = (*GifCaptcha)(nil)

// drawCodeFrame 绘制 GIF 单帧文本：highlight 索引位置使用鲜亮颜色，
// 其它位置使用偏浅的颜色。
func drawCodeFrame(img *image.RGBA, code string, highlight, w, h int) {
	scale := computeScale(h)
	charW := fontWidth*scale + scale
	totalW := charW * len(code)
	startX := (w - totalW) / 2
	charH := fontHeight * scale
	startY := (h - charH) / 2
	for i := 0; i < len(code); i++ {
		var c color.Color
		if i == highlight {
			c = randomColor()
		} else {
			r := uint8(160 + baseutil.RandomInt(80))
			g := uint8(160 + baseutil.RandomInt(80))
			b := uint8(160 + baseutil.RandomInt(80))
			c = color.RGBA{R: r, G: g, B: b, A: 255}
		}
		drawChar(img, code[i], startX+i*charW, startY, scale, c)
	}
}

package captcha

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"
)

// ICaptcha 对应 hutool-captcha ICaptcha 接口。
type ICaptcha interface {
	// CreateCode 生成验证码字符串并渲染图片。
	CreateCode()
	// Code 返回验证码文字内容。
	Code() string
	// Verify 校验用户输入是否正确（通常忽略大小写）。
	Verify(userInputCode string) bool
	// ImageBytes 返回图片字节。
	ImageBytes() []byte
	// ImageBase64 返回图片 Base64。
	ImageBase64() string
	// ImageBase64Data 返回带 data-uri 前缀的 Base64。
	ImageBase64Data() string
	// Write 将图片写出到 io.Writer。
	Write(w io.Writer) error
	// WriteToFile 将图片写出到文件路径。
	WriteToFile(path string) error
}

// AbstractCaptcha 对应 hutool-captcha AbstractCaptcha，所有验证码的公共骨架。
type AbstractCaptcha struct {
	Width          int         // 图片宽度
	Height         int         // 图片高度
	InterfereCount int         // 干扰元素数量
	FontSize       float64     // 字体大小比例（占 Height 的比例，默认 0.75）
	Background     color.Color // 背景色，nil 时表示白色

	generator CodeGenerator
	code      string
	imgBytes  []byte
}

// Code 返回当前验证码字符串。
func (a *AbstractCaptcha) Code() string {
	if a.code == "" {
		a.code = a.ensureGenerator().Generate()
	}
	return a.code
}

// Verify 使用 generator 校验用户输入。
func (a *AbstractCaptcha) Verify(userInputCode string) bool {
	return a.ensureGenerator().Verify(a.Code(), userInputCode)
}

// ImageBytes 返回图片字节；如未生成则为 nil。
func (a *AbstractCaptcha) ImageBytes() []byte { return a.imgBytes }

// ImageBase64 返回图片 Base64 编码。
func (a *AbstractCaptcha) ImageBase64() string {
	return base64.StdEncoding.EncodeToString(a.getImageBytes())
}

// ImageBase64Data 返回带 data-uri 前缀的 Base64（PNG格式）。
func (a *AbstractCaptcha) ImageBase64Data() string {
	return "data:image/png;base64," + a.ImageBase64()
}

// Write 写出图片到 io.Writer。
func (a *AbstractCaptcha) Write(w io.Writer) error {
	b := a.getImageBytes()
	if len(b) == 0 {
		return fmt.Errorf("gkcaptcha: empty image, call CreateCode first")
	}
	_, err := w.Write(b)
	return err
}

// WriteToFile 写出到文件。
func (a *AbstractCaptcha) WriteToFile(path string) error {
	b := a.getImageBytes()
	if len(b) == 0 {
		return fmt.Errorf("gkcaptcha: empty image, call CreateCode first")
	}
	return os.WriteFile(path, b, 0o644)
}

// Generator 返回底层 CodeGenerator。
func (a *AbstractCaptcha) Generator() CodeGenerator { return a.generator }

// SetGenerator 替换 CodeGenerator 并重置状态。
func (a *AbstractCaptcha) SetGenerator(g CodeGenerator) {
	a.generator = g
	a.code = ""
	a.imgBytes = nil
}

// SetBackground 设置背景色。
func (a *AbstractCaptcha) SetBackground(bg color.Color) { a.Background = bg }

// getImageBytes 惰性生成图片。
func (a *AbstractCaptcha) getImageBytes() []byte {
	return a.imgBytes
}

func (a *AbstractCaptcha) ensureGenerator() CodeGenerator {
	if a.generator == nil {
		a.generator = NewRandomGenerator(5)
	}
	return a.generator
}

func (a *AbstractCaptcha) generateCode() {
	a.code = a.ensureGenerator().Generate()
}

func (a *AbstractCaptcha) setImageBytes(b []byte) {
	a.imgBytes = b
}

func (a *AbstractCaptcha) bg() color.Color {
	if a.Background == nil {
		return color.White
	}
	return a.Background
}

// VerifyIgnoreCase 忽略大小写比较（兼容工具函数）。
func VerifyIgnoreCase(code, input string) bool {
	return strings.EqualFold(strings.TrimSpace(code), strings.TrimSpace(input))
}

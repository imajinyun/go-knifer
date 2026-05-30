package captcha

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"
)

// ICaptcha mirrors the hutool-captcha ICaptcha interface.
type ICaptcha interface {
	// CreateCode generates the captcha text and renders the image.
	CreateCode()
	// Code returns the captcha text.
	Code() string
	// Verify reports whether the user input is valid, usually case-insensitively.
	Verify(userInputCode string) bool
	// ImageBytes returns the encoded image bytes.
	ImageBytes() []byte
	// ImageBase64 returns the Base64-encoded image.
	ImageBase64() string
	// ImageBase64Data returns the Base64 image with a data URI prefix.
	ImageBase64Data() string
	// Write writes the image to an io.Writer.
	Write(w io.Writer) error
	// WriteToFile writes the image to a file path.
	WriteToFile(path string) error
}

// AbstractCaptcha mirrors hutool-captcha AbstractCaptcha and holds shared captcha state.
type AbstractCaptcha struct {
	Width          int         // Image width.
	Height         int         // Image height.
	InterfereCount int         // Number of interference elements.
	FontSize       float64     // Font size ratio against Height; default is 0.75.
	Background     color.Color // Background color; nil means white.

	generator CodeGenerator
	code      string
	imgBytes  []byte
}

// Code returns the current captcha text.
func (a *AbstractCaptcha) Code() string {
	if a.code == "" {
		a.code = a.ensureGenerator().Generate()
	}
	return a.code
}

// Verify uses the generator to validate user input.
func (a *AbstractCaptcha) Verify(userInputCode string) bool {
	return a.ensureGenerator().Verify(a.Code(), userInputCode)
}

// ImageBytes returns image bytes, or nil if not generated yet.
func (a *AbstractCaptcha) ImageBytes() []byte { return a.imgBytes }

// ImageBase64 returns the Base64-encoded image.
func (a *AbstractCaptcha) ImageBase64() string {
	return base64.StdEncoding.EncodeToString(a.getImageBytes())
}

// ImageBase64Data returns a PNG data URI containing the Base64 image.
func (a *AbstractCaptcha) ImageBase64Data() string {
	return "data:image/png;base64," + a.ImageBase64()
}

// Write writes the image to an io.Writer.
func (a *AbstractCaptcha) Write(w io.Writer) error {
	b := a.getImageBytes()
	if len(b) == 0 {
		return fmt.Errorf("gkcaptcha: empty image, call CreateCode first")
	}
	_, err := w.Write(b)
	return err
}

// WriteToFile writes the image to a file.
func (a *AbstractCaptcha) WriteToFile(path string) error {
	b := a.getImageBytes()
	if len(b) == 0 {
		return fmt.Errorf("gkcaptcha: empty image, call CreateCode first")
	}
	return os.WriteFile(path, b, 0o644)
}

// Generator returns the underlying CodeGenerator.
func (a *AbstractCaptcha) Generator() CodeGenerator { return a.generator }

// SetGenerator replaces the CodeGenerator and resets generated state.
func (a *AbstractCaptcha) SetGenerator(g CodeGenerator) {
	a.generator = g
	a.code = ""
	a.imgBytes = nil
}

// SetBackground sets the background color.
func (a *AbstractCaptcha) SetBackground(bg color.Color) { a.Background = bg }

// getImageBytes returns lazily generated image bytes.
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

// VerifyIgnoreCase compares code and input case-insensitively for compatibility helpers.
func VerifyIgnoreCase(code, input string) bool {
	return strings.EqualFold(strings.TrimSpace(code), strings.TrimSpace(input))
}

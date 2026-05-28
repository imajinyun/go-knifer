package vcaptcha

import "github.com/imajinyun/go-knifer/internal/captcha"

// Captcha is the common captcha interface.
type Captcha = captcha.ICaptcha

// ICaptcha is the common captcha interface.
type ICaptcha = captcha.ICaptcha

// CodeGenerator generates captcha verification code.
type CodeGenerator = captcha.CodeGenerator

// AbstractCaptcha contains common captcha fields and behavior.
type AbstractCaptcha = captcha.AbstractCaptcha

// LineCaptcha draws line-interference captcha images.
type LineCaptcha = captcha.LineCaptcha

// CircleCaptcha draws circle-interference captcha images.
type CircleCaptcha = captcha.CircleCaptcha

// ShearCaptcha draws sheared captcha images.
type ShearCaptcha = captcha.ShearCaptcha

// GifCaptcha draws animated GIF captcha images.
type GifCaptcha = captcha.GifCaptcha

// RandomGenerator generates random captcha strings.
type RandomGenerator = captcha.RandomGenerator

// MathGenerator generates math-expression captcha strings.
type MathGenerator = captcha.MathGenerator

// VerifyCaptchaIgnoreCase verifies code ignoring case.
func VerifyCaptchaIgnoreCase(code, input string) bool { return captcha.VerifyIgnoreCase(code, input) }

// NewRandomGenerator creates a random captcha generator.
func NewRandomGenerator(length int) *RandomGenerator { return captcha.NewRandomGenerator(length) }

// NewRandomGeneratorWithBase creates a random captcha generator using base.
func NewRandomGeneratorWithBase(base string, length int) *RandomGenerator {
	return captcha.NewRandomGeneratorWithBase(base, length)
}

// NewMathGenerator creates a math captcha generator.
func NewMathGenerator() *MathGenerator { return captcha.NewMathGenerator() }

// NewLineCaptcha creates a line-interference captcha.
func NewLineCaptcha(width, height int) *LineCaptcha { return captcha.NewLineCaptcha(width, height) }

// NewLineCaptchaWith creates a line-interference captcha with options.
func NewLineCaptchaWith(width, height, codeCount, lineCount int) *LineCaptcha {
	return captcha.NewLineCaptchaWith(width, height, codeCount, lineCount)
}

// NewCircleCaptcha creates a circle-interference captcha.
func NewCircleCaptcha(width, height int) *CircleCaptcha {
	return captcha.NewCircleCaptcha(width, height)
}

// NewCircleCaptchaWith creates a circle-interference captcha with options.
func NewCircleCaptchaWith(width, height, codeCount, circleCount int) *CircleCaptcha {
	return captcha.NewCircleCaptchaWith(width, height, codeCount, circleCount)
}

// NewShearCaptcha creates a sheared captcha.
func NewShearCaptcha(width, height int) *ShearCaptcha { return captcha.NewShearCaptcha(width, height) }

// NewShearCaptchaWith creates a sheared captcha with options.
func NewShearCaptchaWith(width, height, codeCount, thickness int) *ShearCaptcha {
	return captcha.NewShearCaptchaWith(width, height, codeCount, thickness)
}

// NewGifCaptcha creates a GIF captcha.
func NewGifCaptcha(width, height int) *GifCaptcha { return captcha.NewGifCaptcha(width, height) }

// NewGifCaptchaWith creates a GIF captcha with options.
func NewGifCaptchaWith(width, height, codeCount int) *GifCaptcha {
	return captcha.NewGifCaptchaWith(width, height, codeCount)
}

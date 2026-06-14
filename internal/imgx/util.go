package imgx

// CaptchaUtil-style package functions create graphical captchas.
//
// Following Go style, factory methods are exposed as package-level functions.

// CreateLineCaptcha creates a line captcha with 5 characters and 150 lines by default.
func CreateLineCaptcha(width, height int) *LineCaptcha {
	return CreateLineCaptchaWithOptions(width, height)
}

// CreateLineCaptchaWithOptions creates a line captcha customized by options.
func CreateLineCaptchaWithOptions(width, height int, opts ...CaptchaOption) *LineCaptcha {
	return NewLineCaptchaWithOptions(width, height, opts...)
}

// CreateLineCaptchaWith creates a line captcha with custom options.
func CreateLineCaptchaWith(width, height, codeCount, lineCount int) *LineCaptcha {
	return CreateLineCaptchaWithOptions(width, height, WithGenerator(NewRandomGenerator(codeCount)), WithInterfereCount(lineCount))
}

// CreateLineCaptchaByGenerator creates a line captcha with a custom generator.
func CreateLineCaptchaByGenerator(width, height int, generator CodeGenerator, lineCount int) *LineCaptcha {
	return CreateLineCaptchaWithOptions(width, height, WithGenerator(generator), WithInterfereCount(lineCount))
}

// CreateCircleCaptcha creates a circle captcha with 5 characters and 15 circles by default.
func CreateCircleCaptcha(width, height int) *CircleCaptcha {
	return CreateCircleCaptchaWithOptions(width, height)
}

// CreateCircleCaptchaWithOptions creates a circle captcha customized by options.
func CreateCircleCaptchaWithOptions(width, height int, opts ...CaptchaOption) *CircleCaptcha {
	return NewCircleCaptchaWithOptions(width, height, opts...)
}

// CreateCircleCaptchaWith creates a circle captcha with custom options.
func CreateCircleCaptchaWith(width, height, codeCount, circleCount int) *CircleCaptcha {
	return CreateCircleCaptchaWithOptions(width, height, WithGenerator(NewRandomGenerator(codeCount)), WithInterfereCount(circleCount))
}

// CreateCircleCaptchaByGenerator creates a circle captcha with a custom generator.
func CreateCircleCaptchaByGenerator(width, height int, generator CodeGenerator, circleCount int) *CircleCaptcha {
	return CreateCircleCaptchaWithOptions(width, height, WithGenerator(generator), WithInterfereCount(circleCount))
}

// CreateShearCaptcha creates a shear captcha with 5 characters and line width 4 by default.
func CreateShearCaptcha(width, height int) *ShearCaptcha {
	return CreateShearCaptchaWithOptions(width, height)
}

// CreateShearCaptchaWithOptions creates a shear captcha customized by options.
func CreateShearCaptchaWithOptions(width, height int, opts ...CaptchaOption) *ShearCaptcha {
	return NewShearCaptchaWithOptions(width, height, opts...)
}

// CreateShearCaptchaWith creates a shear captcha with custom options.
func CreateShearCaptchaWith(width, height, codeCount, thickness int) *ShearCaptcha {
	return CreateShearCaptchaWithOptions(width, height, WithGenerator(NewRandomGenerator(codeCount)), WithInterfereCount(thickness))
}

// CreateShearCaptchaByGenerator creates a shear captcha with a custom generator.
func CreateShearCaptchaByGenerator(width, height int, generator CodeGenerator, thickness int) *ShearCaptcha {
	return CreateShearCaptchaWithOptions(width, height, WithGenerator(generator), WithInterfereCount(thickness))
}

// CreateGifCaptcha creates an animated GIF captcha with 5 characters by default.
func CreateGifCaptcha(width, height int) *GifCaptcha {
	return CreateGifCaptchaWithOptions(width, height)
}

// CreateGifCaptchaWithOptions creates an animated GIF captcha customized by options.
func CreateGifCaptchaWithOptions(width, height int, opts ...CaptchaOption) *GifCaptcha {
	return NewGifCaptchaWithOptions(width, height, opts...)
}

// CreateGifCaptchaWith creates an animated GIF captcha with a custom character count.
func CreateGifCaptchaWith(width, height, codeCount int) *GifCaptcha {
	return CreateGifCaptchaWithOptions(width, height, WithGenerator(NewRandomGenerator(codeCount)))
}

// CreateGifCaptchaByGenerator creates an animated GIF captcha with a custom generator.
func CreateGifCaptchaByGenerator(width, height int, generator CodeGenerator) *GifCaptcha {
	return CreateGifCaptchaWithOptions(width, height, WithGenerator(generator))
}

package captcha

// CaptchaUtil-style package functions create graphical captchas.
//
// Following Go style, factory methods are exposed as package-level functions.

// CreateLineCaptcha creates a line captcha with 5 characters and 150 lines by default.
func CreateLineCaptcha(width, height int) *LineCaptcha {
	return NewLineCaptcha(width, height)
}

// CreateLineCaptchaWith creates a line captcha with custom options.
func CreateLineCaptchaWith(width, height, codeCount, lineCount int) *LineCaptcha {
	return NewLineCaptchaWith(width, height, codeCount, lineCount)
}

// CreateLineCaptchaByGenerator creates a line captcha with a custom generator.
func CreateLineCaptchaByGenerator(width, height int, generator CodeGenerator, lineCount int) *LineCaptcha {
	c := &LineCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = lineCount
	c.SetGenerator(generator)
	return c
}

// CreateCircleCaptcha creates a circle captcha with 5 characters and 15 circles by default.
func CreateCircleCaptcha(width, height int) *CircleCaptcha {
	return NewCircleCaptcha(width, height)
}

// CreateCircleCaptchaWith creates a circle captcha with custom options.
func CreateCircleCaptchaWith(width, height, codeCount, circleCount int) *CircleCaptcha {
	return NewCircleCaptchaWith(width, height, codeCount, circleCount)
}

// CreateCircleCaptchaByGenerator creates a circle captcha with a custom generator.
func CreateCircleCaptchaByGenerator(width, height int, generator CodeGenerator, circleCount int) *CircleCaptcha {
	c := &CircleCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = circleCount
	c.SetGenerator(generator)
	return c
}

// CreateShearCaptcha creates a shear captcha with 5 characters and line width 4 by default.
func CreateShearCaptcha(width, height int) *ShearCaptcha {
	return NewShearCaptcha(width, height)
}

// CreateShearCaptchaWith creates a shear captcha with custom options.
func CreateShearCaptchaWith(width, height, codeCount, thickness int) *ShearCaptcha {
	return NewShearCaptchaWith(width, height, codeCount, thickness)
}

// CreateShearCaptchaByGenerator creates a shear captcha with a custom generator.
func CreateShearCaptchaByGenerator(width, height int, generator CodeGenerator, thickness int) *ShearCaptcha {
	c := &ShearCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = thickness
	c.SetGenerator(generator)
	return c
}

// CreateGifCaptcha creates an animated GIF captcha with 5 characters by default.
func CreateGifCaptcha(width, height int) *GifCaptcha {
	return NewGifCaptcha(width, height)
}

// CreateGifCaptchaWith creates an animated GIF captcha with a custom character count.
func CreateGifCaptchaWith(width, height, codeCount int) *GifCaptcha {
	return NewGifCaptchaWith(width, height, codeCount)
}

// CreateGifCaptchaByGenerator creates an animated GIF captcha with a custom generator.
func CreateGifCaptchaByGenerator(width, height int, generator CodeGenerator) *GifCaptcha {
	c := &GifCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = 10
	c.Repeat = 0
	c.Delay = 10
	c.SetGenerator(generator)
	return c
}

package captcha

// CaptchaUtil 对应 hutool CaptchaUtil：图形验证码工厂入口。
//
// 因 Go 语言风格，所有方法以包级函数提供。

// CreateLineCaptcha 创建线干扰验证码（默认 5 位字符，150 条干扰线）。
func CreateLineCaptcha(width, height int) *LineCaptcha {
	return NewLineCaptcha(width, height)
}

// CreateLineCaptchaWith 自定义参数创建线干扰验证码。
func CreateLineCaptchaWith(width, height, codeCount, lineCount int) *LineCaptcha {
	return NewLineCaptchaWith(width, height, codeCount, lineCount)
}

// CreateLineCaptchaByGenerator 使用自定义 generator 创建线干扰验证码。
func CreateLineCaptchaByGenerator(width, height int, generator CodeGenerator, lineCount int) *LineCaptcha {
	c := &LineCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = lineCount
	c.SetGenerator(generator)
	return c
}

// CreateCircleCaptcha 创建圆圈干扰验证码（默认 5 位字符，15 个干扰圆）。
func CreateCircleCaptcha(width, height int) *CircleCaptcha {
	return NewCircleCaptcha(width, height)
}

// CreateCircleCaptchaWith 自定义参数。
func CreateCircleCaptchaWith(width, height, codeCount, circleCount int) *CircleCaptcha {
	return NewCircleCaptchaWith(width, height, codeCount, circleCount)
}

// CreateCircleCaptchaByGenerator 使用自定义 generator 创建。
func CreateCircleCaptchaByGenerator(width, height int, generator CodeGenerator, circleCount int) *CircleCaptcha {
	c := &CircleCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = circleCount
	c.SetGenerator(generator)
	return c
}

// CreateShearCaptcha 创建扭曲验证码（默认 5 位字符，干扰线宽 4）。
func CreateShearCaptcha(width, height int) *ShearCaptcha {
	return NewShearCaptcha(width, height)
}

// CreateShearCaptchaWith 自定义参数。
func CreateShearCaptchaWith(width, height, codeCount, thickness int) *ShearCaptcha {
	return NewShearCaptchaWith(width, height, codeCount, thickness)
}

// CreateShearCaptchaByGenerator 使用自定义 generator 创建。
func CreateShearCaptchaByGenerator(width, height int, generator CodeGenerator, thickness int) *ShearCaptcha {
	c := &ShearCaptcha{}
	c.Width = width
	c.Height = height
	c.InterfereCount = thickness
	c.SetGenerator(generator)
	return c
}

// CreateGifCaptcha 创建 GIF 动图验证码（默认 5 位字符）。
func CreateGifCaptcha(width, height int) *GifCaptcha {
	return NewGifCaptcha(width, height)
}

// CreateGifCaptchaWith 自定义字符数。
func CreateGifCaptchaWith(width, height, codeCount int) *GifCaptcha {
	return NewGifCaptchaWith(width, height, codeCount)
}

// CreateGifCaptchaByGenerator 使用自定义 generator 创建。
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

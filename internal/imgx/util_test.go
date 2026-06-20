package imgx

import "testing"

func TestCaptchaUtilFactories(t *testing.T) {
	if CreateLineCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
	if CreateCircleCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
	if CreateShearCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
	if CreateGifCaptcha(120, 40) == nil {
		t.Fatal("nil")
	}
}

func TestCaptchaUtilFactoryVariants(t *testing.T) {
	line := CreateLineCaptchaWith(120, 40, 3, 7)
	if line.InterfereCount != 7 || len(line.Code()) != 3 {
		t.Fatalf("CreateLineCaptchaWith interfere=%d code=%q", line.InterfereCount, line.Code())
	}
	lineByGenerator := CreateLineCaptchaByGenerator(120, 40, fixedGenerator{code: "LINE"}, 8)
	if lineByGenerator.InterfereCount != 8 || lineByGenerator.Code() != "LINE" {
		t.Fatalf("CreateLineCaptchaByGenerator interfere=%d code=%q", lineByGenerator.InterfereCount, lineByGenerator.Code())
	}

	circle := CreateCircleCaptchaWith(120, 40, 4, 9)
	if circle.InterfereCount != 9 || len(circle.Code()) != 4 {
		t.Fatalf("CreateCircleCaptchaWith interfere=%d code=%q", circle.InterfereCount, circle.Code())
	}
	circleByGenerator := CreateCircleCaptchaByGenerator(120, 40, fixedGenerator{code: "CIRC"}, 10)
	if circleByGenerator.InterfereCount != 10 || circleByGenerator.Code() != "CIRC" {
		t.Fatalf("CreateCircleCaptchaByGenerator interfere=%d code=%q", circleByGenerator.InterfereCount, circleByGenerator.Code())
	}

	shear := CreateShearCaptchaWith(120, 40, 2, 6)
	if shear.InterfereCount != 6 || len(shear.Code()) != 2 {
		t.Fatalf("CreateShearCaptchaWith interfere=%d code=%q", shear.InterfereCount, shear.Code())
	}
	shearByGenerator := CreateShearCaptchaByGenerator(120, 40, fixedGenerator{code: "SHR"}, 5)
	if shearByGenerator.InterfereCount != 5 || shearByGenerator.Code() != "SHR" {
		t.Fatalf("CreateShearCaptchaByGenerator interfere=%d code=%q", shearByGenerator.InterfereCount, shearByGenerator.Code())
	}

	gifCaptcha := CreateGifCaptchaWith(120, 40, 2)
	if len(gifCaptcha.Code()) != 2 {
		t.Fatalf("CreateGifCaptchaWith code=%q", gifCaptcha.Code())
	}
	gifByGenerator := CreateGifCaptchaByGenerator(120, 40, fixedGenerator{code: "GIF"})
	if gifByGenerator.Code() != "GIF" {
		t.Fatalf("CreateGifCaptchaByGenerator code=%q", gifByGenerator.Code())
	}
}

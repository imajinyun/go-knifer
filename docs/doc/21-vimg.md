# vimg Quickstart

`vimg` provides basic image processing and captcha helpers, covering image metadata reads, format conversion, thumbnail generation, captcha generation/verification, and file writes.

## Read image information

```go
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/imajinyun/go-knifer/vimg"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 80, 40))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}

	width, height, format, err := vimg.Info(bytes.NewReader(buf.Bytes()))
	if err != nil {
		panic(err)
	}
	fmt.Println(width, height, format)
}
```

## Convert formats and generate thumbnails

```go
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/imajinyun/go-knifer/vimg"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 160, 80))
	var src bytes.Buffer
	if err := png.Encode(&src, img); err != nil {
		panic(err)
	}

	var jpegOut bytes.Buffer
	if err := vimg.ConvertFormat(&jpegOut, bytes.NewReader(src.Bytes()), "jpeg"); err != nil {
		panic(err)
	}

	var thumb bytes.Buffer
	if err := vimg.Thumbnail(&thumb, bytes.NewReader(src.Bytes()), 64, "png"); err != nil {
		panic(err)
	}
	fmt.Println(jpegOut.Len() > 0, thumb.Len() > 0)
}
```

## Generate and verify captchas

```go
package main

import (
	"fmt"
	"strings"

	"github.com/imajinyun/go-knifer/vimg"
)

func main() {
	captcha := vimg.NewLineCaptchaWithOptions(120, 48,
		vimg.WithGenerator(vimg.NewRandomGeneratorWithBase("ABC123", 4)),
		vimg.WithInterfereCount(8),
	)
	captcha.CreateCode()

	code := captcha.Code()
	fmt.Println(captcha.Verify(strings.ToLower(code)))
	fmt.Println(strings.HasPrefix(captcha.ImageBase64Data(), "data:image/png;base64,"))
}
```

## Use math captchas and write options

```go
package main

import (
	"fmt"
	"path/filepath"

	"github.com/imajinyun/go-knifer/vimg"
)

func main() {
	generator := vimg.NewMathGeneratorWith(1, false)
	fmt.Println(generator.Verify("1+1=", "2"))

	captcha := vimg.NewGifCaptchaWithOptions(120, 48, vimg.WithGenerator(generator))
	captcha.CreateCode()
	path := filepath.Join("tmp", "captcha.gif")
	_ = captcha.WriteToFileWithOptions(path,
		vimg.WithCreateParents(true),
		vimg.WithFilePerm(0o600),
		vimg.WithOverwrite(true),
	)
}
```

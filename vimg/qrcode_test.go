package vimg

import (
	"bytes"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func TestFacadeQRCode(t *testing.T) {
	const content = "facade qr payload"
	pngBytes, err := QRCodePNG(content,
		WithQRCodeSize(150),
		WithQRCodeMargin(2),
		WithQRCodeErrorCorrection(QRErrorCorrectionQuartile),
	)
	if err != nil {
		t.Fatalf("QRCodePNG: %v", err)
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("png config: %v", err)
	}
	if cfg.Width != 150 || cfg.Height != 150 {
		t.Fatalf("QRCodePNG size = %dx%d, want 150x150", cfg.Width, cfg.Height)
	}
	result, err := DecodeQRCode(bytes.NewReader(pngBytes), WithDecodeTryHarder(true))
	if err != nil {
		t.Fatalf("DecodeQRCode: %v", err)
	}
	if result.Text != content || result.Format != BarcodeFormatQRCode {
		t.Fatalf("DecodeQRCode = (%q, %v), want (%q, %v)", result.Text, result.Format, content, BarcodeFormatQRCode)
	}
}

func TestFacadeBarcodeRenderers(t *testing.T) {
	pngBytes, err := BarcodePNG("123456789012", BarcodeFormatEAN13, WithBarcodeSize(220, 90))
	if err != nil {
		t.Fatalf("BarcodePNG: %v", err)
	}
	if _, err := png.DecodeConfig(bytes.NewReader(pngBytes)); err != nil {
		t.Fatalf("png config: %v", err)
	}

	svg, err := QRCodeSVG("facade svg", WithQRCodeForeground(color.Black))
	if err != nil {
		t.Fatalf("QRCodeSVG: %v", err)
	}
	if !strings.Contains(svg, "<svg") || !strings.Contains(svg, "<path") {
		t.Fatalf("QRCodeSVG missing expected tags: %q", svg[:min(len(svg), 80)])
	}

	ascii, err := QRCodeASCII("facade ascii")
	if err != nil {
		t.Fatalf("QRCodeASCII: %v", err)
	}
	if !strings.Contains(ascii, "██") {
		t.Fatalf("QRCodeASCII missing block chars: %q", ascii[:min(len(ascii), 80)])
	}

	data, err := BarcodeBase64Data("facade data", BarcodeFormatQRCode)
	if err != nil {
		t.Fatalf("BarcodeBase64Data: %v", err)
	}
	if !strings.HasPrefix(data, "data:image/png;base64,") {
		t.Fatalf("BarcodeBase64Data prefix = %q", data[:min(len(data), 30)])
	}
}

package str

import (
	"bytes"
	"errors"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
)

func TestBOM(t *testing.T) {
	tests := []struct {
		name   string
		prefix []byte
		want   BOMType
	}{
		{name: "utf8", prefix: []byte{0xEF, 0xBB, 0xBF}, want: BOMUTF8},
		{name: "utf16 le", prefix: []byte{0xFF, 0xFE}, want: BOMUTF16LE},
		{name: "utf16 be", prefix: []byte{0xFE, 0xFF}, want: BOMUTF16BE},
		{name: "utf32 le", prefix: []byte{0xFF, 0xFE, 0x00, 0x00}, want: BOMUTF32LE},
		{name: "utf32 be", prefix: []byte{0x00, 0x00, 0xFE, 0xFF}, want: BOMUTF32BE},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := append(append([]byte(nil), tt.prefix...), 'g', 'o')
			if got := HasBOM(data); got != tt.want {
				t.Fatalf("HasBOM = %q, want %q", got, tt.want)
			}
			if got := StripBOM(data); !bytes.Equal(got, []byte("go")) {
				t.Fatalf("StripBOM = %v", got)
			}
		})
	}

	plain := []byte("go")
	stripped := StripBOM(plain)
	if !bytes.Equal(stripped, plain) {
		t.Fatalf("StripBOM without bom = %v", stripped)
	}
	stripped[0] = 'n'
	if string(plain) != "go" {
		t.Fatalf("StripBOM returned aliased data; original = %q", plain)
	}
}

func TestCharsetRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		charset string
		text    string
	}{
		{name: "utf8", charset: "utf-8", text: "plain"},
		{name: "gbk", charset: "gbk", text: "中文"},
		{name: "gb18030", charset: "gb18030", text: "中文"},
		{name: "big5", charset: "big5", text: "中文"},
		{name: "shift jis", charset: "shift_jis", text: "日本語"},
		{name: "euc kr", charset: "euc-kr", text: "한국어"},
		{name: "iso 8859 1", charset: "iso-8859-1", text: "é"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := FromUTF8([]byte(tt.text), tt.charset)
			if err != nil {
				t.Fatalf("FromUTF8 error = %v", err)
			}
			utf8, err := ToUTF8(append([]byte{0xEF, 0xBB, 0xBF}, encoded...), tt.charset)
			if err != nil {
				t.Fatalf("ToUTF8 error = %v", err)
			}
			if string(utf8) != tt.text {
				t.Fatalf("ToUTF8 = %q, want %q", utf8, tt.text)
			}
		})
	}
}

func TestCharsetUnsupported(t *testing.T) {
	_, err := ToUTF8([]byte("x"), "unknown-charset")
	if !errors.Is(err, knifer.ErrCodeUnsupported) {
		t.Fatalf("ToUTF8 unsupported err = %v", err)
	}
	_, err = FromUTF8([]byte("x"), "unknown-charset")
	if !errors.Is(err, knifer.ErrCodeUnsupported) {
		t.Fatalf("FromUTF8 unsupported err = %v", err)
	}
}

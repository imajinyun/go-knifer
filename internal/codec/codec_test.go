package codec

import "testing"

func TestBase64(t *testing.T) {
	src := "Hello, 世界"
	enc := Base64EncodeStr(src)
	dec, err := Base64DecodeStr(enc)
	if err != nil {
		t.Fatalf("Base64 decode err: %v", err)
	}
	if dec != src {
		t.Fatalf("Base64 mismatch: %q", dec)
	}
}

func TestBase64URL(t *testing.T) {
	data := []byte{0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}
	enc := Base64URLEncode(data)
	dec, err := Base64URLDecode(enc)
	if err != nil {
		t.Fatalf("Base64URL decode err: %v", err)
	}
	if string(dec) != string(data) {
		t.Fatalf("Base64URL mismatch")
	}
}

func TestHex(t *testing.T) {
	if HexEncodeStr("AB") != "4142" {
		t.Fatalf("HexEncode: %s", HexEncodeStr("AB"))
	}
	got, err := HexDecodeStr("4142")
	if err != nil || got != "AB" {
		t.Fatalf("HexDecode: %v %q", err, got)
	}
}

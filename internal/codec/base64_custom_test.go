package codec

import (
	"encoding/base64"
	"testing"
)

func TestBase64EncodeWithEncoding(t *testing.T) {
	data := []byte("custom?")
	tests := []struct {
		name string
		enc  *base64.Encoding
		want string
	}{
		{name: "custom raw url encoding", enc: base64.RawURLEncoding, want: base64.RawURLEncoding.EncodeToString(data)},
		{name: "nil uses standard encoding", enc: nil, want: base64.StdEncoding.EncodeToString(data)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64EncodeWithEncoding(data, tt.enc); got != tt.want {
				t.Fatalf("Base64EncodeWithEncoding(%q) = %q, want %q", data, got, tt.want)
			}
		})
	}
}

func TestBase64DecodeWithEncoding(t *testing.T) {
	data := []byte("custom?")
	tests := []struct {
		name string
		in   string
		enc  *base64.Encoding
		want []byte
	}{
		{name: "custom raw url encoding", in: base64.RawURLEncoding.EncodeToString(data), enc: base64.RawURLEncoding, want: data},
		{name: "nil uses standard encoding", in: base64.StdEncoding.EncodeToString(data), enc: nil, want: data},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64DecodeWithEncoding(tt.in, tt.enc)
			if err != nil {
				t.Fatalf("Base64DecodeWithEncoding(%q) error = %v", tt.in, err)
			}
			if string(got) != string(tt.want) {
				t.Fatalf("Base64DecodeWithEncoding(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

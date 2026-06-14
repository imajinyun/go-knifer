package codec

import "testing"

func TestBase64EncodeStr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "ascii", in: "hello", want: "aGVsbG8="},
		{name: "unicode", in: "Hello, 世界", want: "SGVsbG8sIOS4lueVjA=="},
		{name: "empty", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64EncodeStr(tt.in); got != tt.want {
				t.Fatalf("Base64EncodeStr(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBase64DecodeStr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "ascii", in: "aGVsbG8=", want: "hello"},
		{name: "unicode", in: "SGVsbG8sIOS4lueVjA==", want: "Hello, 世界"},
		{name: "empty", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64DecodeStr(tt.in)
			if err != nil {
				t.Fatalf("Base64DecodeStr(%q) error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("Base64DecodeStr(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBase64RoundTrip(t *testing.T) {
	src := "Hello, 世界"
	enc := Base64EncodeStr(src)
	dec, err := Base64DecodeStr(enc)
	if err != nil {
		t.Fatalf("Base64DecodeStr(%q) error = %v", enc, err)
	}
	if dec != src {
		t.Fatalf("Base64DecodeStr(Base64EncodeStr(%q)) = %q, want %q", src, dec, src)
	}
}

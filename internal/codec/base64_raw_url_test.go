package codec

import "testing"

func TestBase64RawURLEncode(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want string
	}{
		{name: "unpadded", in: []byte("custom?"), want: "Y3VzdG9tPw"},
		{name: "empty", in: nil, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64RawURLEncode(tt.in); got != tt.want {
				t.Fatalf("Base64RawURLEncode(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBase64RawURLDecode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []byte
	}{
		{name: "unpadded", in: "Y3VzdG9tPw", want: []byte("custom?")},
		{name: "empty", in: "", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64RawURLDecode(tt.in)
			if err != nil {
				t.Fatalf("Base64RawURLDecode(%q) error = %v", tt.in, err)
			}
			if string(got) != string(tt.want) {
				t.Fatalf("Base64RawURLDecode(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestBase64RawURLRoundTrip(t *testing.T) {
	data := []byte("custom?")
	enc := Base64RawURLEncode(data)
	dec, err := Base64RawURLDecode(enc)
	if err != nil {
		t.Fatalf("Base64RawURLDecode(%q) error = %v", enc, err)
	}
	if string(dec) != string(data) {
		t.Fatalf("Base64RawURLDecode(Base64RawURLEncode(%q)) = %q, want %q", data, dec, data)
	}
}

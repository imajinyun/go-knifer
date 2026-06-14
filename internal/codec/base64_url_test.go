package codec

import "testing"

func TestBase64URLEncode(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want string
	}{
		{name: "url alphabet", in: []byte{0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}, want: "-vv8_f7_"},
		{name: "empty", in: nil, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64URLEncode(tt.in); got != tt.want {
				t.Fatalf("Base64URLEncode(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBase64URLDecode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []byte
	}{
		{name: "url alphabet", in: "-vv8_f7_", want: []byte{0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
		{name: "empty", in: "", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64URLDecode(tt.in)
			if err != nil {
				t.Fatalf("Base64URLDecode(%q) error = %v", tt.in, err)
			}
			if string(got) != string(tt.want) {
				t.Fatalf("Base64URLDecode(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestBase64URLRoundTrip(t *testing.T) {
	data := []byte{0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}
	enc := Base64URLEncode(data)
	dec, err := Base64URLDecode(enc)
	if err != nil {
		t.Fatalf("Base64URLDecode(%q) error = %v", enc, err)
	}
	if string(dec) != string(data) {
		t.Fatalf("Base64URLDecode(Base64URLEncode(%v)) = %v, want %v", data, dec, data)
	}
}

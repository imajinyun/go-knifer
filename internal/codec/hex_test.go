package codec

import "testing"

func TestHexEncodeStr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "ascii", in: "AB", want: "4142"},
		{name: "unicode", in: "世界", want: "e4b896e7958c"},
		{name: "empty", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexEncodeStr(tt.in); got != tt.want {
				t.Fatalf("HexEncodeStr(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestHexDecodeStr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "ascii", in: "4142", want: "AB"},
		{name: "unicode", in: "e4b896e7958c", want: "世界"},
		{name: "empty", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexDecodeStr(tt.in)
			if err != nil {
				t.Fatalf("HexDecodeStr(%q) error = %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("HexDecodeStr(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

package crypto

import "testing"

func TestSignParams(t *testing.T) {
	params := map[string]any{"b": 2, "a": 1, "skip": nil}
	if got := SignParams(params, SHA256Hex, "&", "=", true, "secret"); got != SHA256Hex([]byte("a=1&b=2&secret")) {
		t.Fatalf("SignParams() = %s", got)
	}
	if got := SignParamsSHA256(map[string]any{"b": 2, "a": 1}, "z"); got != SHA256Hex([]byte("a1b2z")) {
		t.Fatalf("SignParamsSHA256() = %s", got)
	}
}

package jwt

import "testing"

func TestSetSignerNilIsSafe(t *testing.T) {
	j := New()
	j.SetSigner(nil)
	if j.Signer() != nil {
		t.Fatal("SetSigner(nil) should clear signer")
	}
	if got := j.Algorithm(); got != "" {
		t.Fatalf("SetSigner(nil) should not write alg, got %q", got)
	}
	if err := j.SetSignerE(nil); err == nil {
		t.Fatal("SetSignerE(nil) should return an error")
	}
}

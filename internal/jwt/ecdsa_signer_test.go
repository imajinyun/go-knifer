package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestECDSASigner_RoundTrip(t *testing.T) {
	cases := []struct {
		alg   string
		curve elliptic.Curve
	}{
		{AlgES256, elliptic.P256()},
		{AlgES384, elliptic.P384()},
		{AlgES512, elliptic.P521()},
	}
	for _, c := range cases {
		priv, err := ecdsa.GenerateKey(c.curve, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		signer, err := NewECDSASigner(c.alg, priv, nil)
		if err != nil {
			t.Fatalf("%s: %v", c.alg, err)
		}
		token, err := New().AddPayloads(map[string]any{"u": 1}).SetSigner(signer).Sign()
		if err != nil {
			t.Fatalf("%s sign: %v", c.alg, err)
		}
		j, err := Of(token)
		if err != nil {
			t.Fatal(err)
		}
		if !j.VerifyWith(signer) {
			t.Fatalf("%s verify failed", c.alg)
		}
	}
}

func TestECDSASignerWithOptionsRandomReader(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := NewECDSASignerWithOptions(AlgES256, priv, nil, WithSignerRandomReader(rand.Reader))
	if err != nil {
		t.Fatalf("NewECDSASignerWithOptions: %v", err)
	}
	token, err := New().SetPayload("u", 1).SetSigner(signer).Sign()
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	j, err := Of(token)
	if err != nil {
		t.Fatal(err)
	}
	if !j.VerifyWith(signer) {
		t.Fatal("verify with custom ECDSA random reader failed")
	}
	verifier, err := NewECDSASignerWithOptions(AlgES256, nil, &priv.PublicKey)
	if err != nil {
		t.Fatalf("public-only signer: %v", err)
	}
	if token, err := New().SetPayload("u", 1).SetSigner(verifier).Sign(); err == nil || token != "" {
		t.Fatalf("Sign should reject public-only ECDSA empty signature, token=%q err=%v", token, err)
	}
}

func TestECDSASigner_CurveMismatch(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewECDSASigner(AlgES256, priv, nil); err == nil {
		t.Fatal("expected curve mismatch error")
	}
}

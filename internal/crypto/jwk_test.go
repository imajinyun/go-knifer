package crypto

import (
	"crypto/rsa"
	"errors"
	"strings"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
)

func testRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	priv, err := GenRSAKey(1024)
	if err != nil {
		t.Fatalf("GenRSAKey: %v", err)
	}
	return priv
}

func TestRSAJWKPublicRoundTrip(t *testing.T) {
	priv := testRSAKey(t)
	jwk, err := RSAPublicKeyToJWK(&priv.PublicKey, "kid-1")
	if err != nil {
		t.Fatalf("RSAPublicKeyToJWK: %v", err)
	}
	if jwk.KeyType != "RSA" || jwk.KeyID != "kid-1" || jwk.Modulus == "" || jwk.Exponent == "" || jwk.Private != "" {
		t.Fatalf("public jwk = %#v", jwk)
	}
	data, err := MarshalJWK(jwk)
	if err != nil {
		t.Fatalf("MarshalJWK: %v", err)
	}
	parsed, err := ParseJWK(data)
	if err != nil {
		t.Fatalf("ParseJWK: %v", err)
	}
	pub, err := JWKToRSAPublicKey(parsed)
	if err != nil {
		t.Fatalf("JWKToRSAPublicKey: %v", err)
	}
	if pub.N.Cmp(priv.N) != 0 || pub.E != priv.E {
		t.Fatalf("public key mismatch")
	}
}

func TestRSAJWKPrivateRoundTrip(t *testing.T) {
	priv := testRSAKey(t)
	jwk, err := RSAPrivateKeyToJWK(priv, "kid-private")
	if err != nil {
		t.Fatalf("RSAPrivateKeyToJWK: %v", err)
	}
	if jwk.Private == "" || jwk.PrimeP == "" || jwk.PrimeQ == "" {
		t.Fatalf("private jwk missing private fields: %#v", jwk)
	}
	parsed, err := JWKToRSAPrivateKey(jwk)
	if err != nil {
		t.Fatalf("JWKToRSAPrivateKey: %v", err)
	}
	if parsed.N.Cmp(priv.N) != 0 || parsed.D.Cmp(priv.D) != 0 {
		t.Fatalf("private key mismatch")
	}
}

func TestJWKSSelectByKeyID(t *testing.T) {
	priv := testRSAKey(t)
	key1, err := RSAPublicKeyToJWK(&priv.PublicKey, "kid-1")
	if err != nil {
		t.Fatalf("RSAPublicKeyToJWK key1: %v", err)
	}
	key2, err := RSAPublicKeyToJWK(&priv.PublicKey, "kid-2")
	if err != nil {
		t.Fatalf("RSAPublicKeyToJWK key2: %v", err)
	}
	data, err := MarshalJWKS([]JWK{key1, key2})
	if err != nil {
		t.Fatalf("MarshalJWKS: %v", err)
	}
	set, err := ParseJWKS(data)
	if err != nil {
		t.Fatalf("ParseJWKS: %v", err)
	}
	selected, err := SelectJWKByKeyID(set, "kid-2")
	if err != nil {
		t.Fatalf("SelectJWKByKeyID: %v", err)
	}
	if selected.KeyID != "kid-2" {
		t.Fatalf("selected kid = %q", selected.KeyID)
	}
	_, err = SelectJWKByKeyID(set, "missing")
	if !errors.Is(err, knifer.ErrCodeNotFound) || !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("SelectJWKByKeyID missing error = %v", err)
	}
}

func TestJWKMalformedErrors(t *testing.T) {
	if _, err := ParseJWK([]byte("{")); !errors.Is(err, ErrInvalidJWK) || !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("ParseJWK malformed error = %v", err)
	}
	if _, err := ParseJWKS([]byte(`{"keys":[]}`)); !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("ParseJWKS empty error = %v", err)
	}
	if _, err := JWKToRSAPublicKey(JWK{KeyType: "EC"}); !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("JWKToRSAPublicKey EC error = %v", err)
	}
	if _, err := JWKToRSAPublicKey(JWK{KeyType: "RSA", Modulus: "!", Exponent: "AQAB"}); !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("JWKToRSAPublicKey malformed modulus error = %v", err)
	}
	if _, err := MarshalJWK(JWK{}); !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("MarshalJWK empty error = %v", err)
	}
	if _, err := MarshalJWKS([]JWK{}); !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("MarshalJWKS empty error = %v", err)
	}
	if _, err := SelectJWKByKeyID(JWKS{Keys: []JWK{{KeyType: "RSA"}}}, " "); !errors.Is(err, ErrInvalidJWK) {
		t.Fatalf("SelectJWKByKeyID empty kid error = %v", err)
	}
	if data, err := MarshalJWKS([]JWK{{KeyType: "RSA", KeyID: "kid"}}); err != nil || strings.Contains(string(data), "http") {
		t.Fatalf("JWKS should be local JSON only: %s, %v", data, err)
	}
}

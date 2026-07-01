package vcrypto_test

import (
	"errors"
	"testing"

	knifer "github.com/imajinyun/knifer-go"
	"github.com/imajinyun/knifer-go/vcrypto"
)

func TestFacadeRSAJWKAndJWKS(t *testing.T) {
	priv, err := vcrypto.GenRSAKey(1024)
	if err != nil {
		t.Fatalf("GenRSAKey: %v", err)
	}
	jwk, err := vcrypto.RSAPublicKeyToJWK(&priv.PublicKey, "kid-1")
	if err != nil {
		t.Fatalf("RSAPublicKeyToJWK: %v", err)
	}
	data, err := vcrypto.MarshalJWK(jwk)
	if err != nil {
		t.Fatalf("MarshalJWK: %v", err)
	}
	parsed, err := vcrypto.ParseJWK(data)
	if err != nil {
		t.Fatalf("ParseJWK: %v", err)
	}
	pub, err := vcrypto.JWKToRSAPublicKey(parsed)
	if err != nil {
		t.Fatalf("JWKToRSAPublicKey: %v", err)
	}
	if pub.N.Cmp(priv.N) != 0 || pub.E != priv.E {
		t.Fatalf("public key mismatch")
	}
	jwksData, err := vcrypto.MarshalJWKS([]vcrypto.JWK{jwk})
	if err != nil {
		t.Fatalf("MarshalJWKS: %v", err)
	}
	set, err := vcrypto.ParseJWKS(jwksData)
	if err != nil {
		t.Fatalf("ParseJWKS: %v", err)
	}
	selected, err := vcrypto.SelectJWKByKeyID(set, "kid-1")
	if err != nil {
		t.Fatalf("SelectJWKByKeyID: %v", err)
	}
	if selected.KeyID != "kid-1" {
		t.Fatalf("selected kid = %q", selected.KeyID)
	}
}

func TestFacadeJWKErrors(t *testing.T) {
	if _, err := vcrypto.ParseJWK([]byte("{")); !errors.Is(err, vcrypto.ErrInvalidJWK) || !errors.Is(err, knifer.ErrCodeInvalidInput) {
		t.Fatalf("ParseJWK malformed error = %v", err)
	}
	if _, err := vcrypto.SelectJWKByKeyID(vcrypto.JWKS{Keys: []vcrypto.JWK{{KeyID: "kid-1"}}}, "missing"); !errors.Is(err, vcrypto.ErrInvalidJWK) || !errors.Is(err, knifer.ErrCodeNotFound) {
		t.Fatalf("SelectJWKByKeyID missing error = %v", err)
	}
}

package vjwt_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"testing"

	"github.com/imajinyun/go-knifer/vjwt"
)

func TestStrictHMACSignerRejectsWeakKey(t *testing.T) {
	if _, err := vjwt.NewHMACSignerStrict(vjwt.JWTAlgHS256, []byte("weak")); err == nil {
		t.Fatal("NewHMACSignerStrict should reject weak key")
	}
	if _, err := vjwt.CreateSignerStrict(vjwt.JWTAlgHS256, []byte("weak")); err == nil {
		t.Fatal("CreateSignerStrict should reject weak key")
	}
	if minBytes, err := vjwt.MinHMACKeyBytes(vjwt.JWTAlgHS256); err != nil || minBytes != vjwt.MinHMACKeyBytesHS256 {
		t.Fatalf("MinHMACKeyBytes = %d, %v", minBytes, err)
	}
}

func TestFacadeSignerFactoriesAndAlgorithms(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	for _, tt := range []struct {
		name string
		fn   func([]byte) vjwt.JWTSigner
		alg  string
	}{
		{name: "JWTSignerHS256", fn: vjwt.JWTSignerHS256, alg: vjwt.JWTAlgHS256},
		{name: "HS256", fn: vjwt.HS256, alg: vjwt.JWTAlgHS256},
		{name: "HS384", fn: vjwt.HS384, alg: vjwt.JWTAlgHS384},
		{name: "HS512", fn: vjwt.HS512, alg: vjwt.JWTAlgHS512},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(key).Algorithm(); got != tt.alg {
				t.Fatalf("Algorithm = %q, want %q", got, tt.alg)
			}
		})
	}
	signer, err := vjwt.JWTSignerHMAC(vjwt.JWTAlgHS384, key)
	if err != nil || signer.Algorithm() != vjwt.JWTAlgHS384 {
		t.Fatalf("JWTSignerHMAC alg=%q err=%v", signer.Algorithm(), err)
	}
	if got := vjwt.MustHMACSigner(vjwt.JWTAlgHS512, key).Algorithm(); got != vjwt.JWTAlgHS512 {
		t.Fatalf("MustHMACSigner alg = %q", got)
	}
	if got := vjwt.AlgorithmName(vjwt.JWTAlgPS256); got != "SHA256withRSA_PSS" {
		t.Fatalf("AlgorithmName(PS256) = %q", got)
	}
	if _, err := vjwt.NewRSAPSSSigner(vjwt.JWTAlgPS256, nil, nil); err == nil {
		t.Fatal("NewRSAPSSSigner(nil keys) error = nil")
	}
	if _, err := vjwt.NewECDSASigner(vjwt.JWTAlgES256, nil, nil); err == nil {
		t.Fatal("NewECDSASigner(nil keys) error = nil")
	}

	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("rsa.GenerateKey: %v", err)
	}
	reader := zeroReader{}
	psSigner, err := vjwt.NewRSAPSSSignerWithOptions(vjwt.JWTAlgPS256, rsaKey, &rsaKey.PublicKey, vjwt.WithSignerRandomReader(reader), vjwt.WithRSAPSSOptions(nil))
	if err != nil || psSigner.Algorithm() != vjwt.JWTAlgPS256 {
		t.Fatalf("NewRSAPSSSignerWithOptions alg=%q err=%v", psSigner.Algorithm(), err)
	}
	if got := vjwt.PS256WithOptions(rsaKey, &rsaKey.PublicKey, vjwt.WithSignerRandomReader(reader)).Algorithm(); got != vjwt.JWTAlgPS256 {
		t.Fatalf("PS256WithOptions alg = %q", got)
	}
	if got := vjwt.PS384(rsaKey, &rsaKey.PublicKey).Algorithm(); got != vjwt.JWTAlgPS384 {
		t.Fatalf("PS384 alg = %q", got)
	}
	if got := vjwt.PS512WithOptions(rsaKey, &rsaKey.PublicKey).Algorithm(); got != vjwt.JWTAlgPS512 {
		t.Fatalf("PS512WithOptions alg = %q", got)
	}

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey: %v", err)
	}
	ecSigner, err := vjwt.JWTSignerECDSA(vjwt.JWTAlgES256, ecdsaKey, &ecdsaKey.PublicKey)
	if err != nil || ecSigner.Algorithm() != vjwt.JWTAlgES256 {
		t.Fatalf("JWTSignerECDSA alg=%q err=%v", ecSigner.Algorithm(), err)
	}
	if got := vjwt.JWTSignerES256(ecdsaKey, &ecdsaKey.PublicKey).Algorithm(); got != vjwt.JWTAlgES256 {
		t.Fatalf("JWTSignerES256 alg = %q", got)
	}
	if got := vjwt.ES256WithOptions(ecdsaKey, &ecdsaKey.PublicKey, vjwt.WithSignerRandomReader(reader)).Algorithm(); got != vjwt.JWTAlgES256 {
		t.Fatalf("ES256WithOptions alg = %q", got)
	}
	ecdsa384Key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey P384: %v", err)
	}
	if got := vjwt.ES384(ecdsa384Key, &ecdsa384Key.PublicKey).Algorithm(); got != vjwt.JWTAlgES384 {
		t.Fatalf("ES384 alg = %q", got)
	}
	ecdsa521Key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey P521: %v", err)
	}
	if got := vjwt.ES512WithOptions(ecdsa521Key, &ecdsa521Key.PublicKey).Algorithm(); got != vjwt.JWTAlgES512 {
		t.Fatalf("ES512WithOptions alg = %q", got)
	}
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 1
	}
	return len(p), nil
}

var _ io.Reader = zeroReader{}

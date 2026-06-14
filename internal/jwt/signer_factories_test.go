package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestSignerUtilFactories(t *testing.T) {
	// HS*
	if HS256([]byte("k")).Algorithm() != AlgHS256 {
		t.Fatal()
	}
	if HS384([]byte("k")).Algorithm() != AlgHS384 {
		t.Fatal()
	}
	if HS512([]byte("k")).Algorithm() != AlgHS512 {
		t.Fatal()
	}
	// PS*
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	if PS256(priv, nil).Algorithm() != AlgPS256 {
		t.Fatal()
	}

	// ES*
	ec, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if ES256(ec, nil).Algorithm() != AlgES256 {
		t.Fatal()
	}
}

func TestAlgorithmName(t *testing.T) {
	pairs := map[string]string{
		AlgHS256: "HmacSHA256",
		AlgHS384: "HmacSHA384",
		AlgHS512: "HmacSHA512",
		AlgPS256: "SHA256withRSA_PSS",
		AlgES256: "SHA256withECDSA",
	}
	for id, name := range pairs {
		if got := AlgorithmName(id); got != name {
			t.Fatalf("%s -> %s, want %s", id, got, name)
		}
	}
	if AlgorithmName("UNKNOWN") != "UNKNOWN" {
		t.Fatal("unknown should be returned as-is")
	}
}

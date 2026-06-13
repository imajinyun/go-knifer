package crypto

import (
	"bytes"
	stdcrypto "crypto"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestRSAOAEP(t *testing.T) {
	priv, err := GenerateRSAKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("hello")
	cipherText, err := RSAEncryptOAEP(plain, &priv.PublicKey, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := RSADecryptOAEP(cipherText, priv, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("RSADecryptOAEP() = %q", out)
	}
}

func TestRSAPKCS1PSS(t *testing.T) {
	priv, err := GenerateRSAKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("signed payload")
	digest := sha256.Sum256(plain)
	pssSig, err := RSASignPSS(priv, stdcrypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerifyPSS(&priv.PublicKey, stdcrypto.SHA256, digest[:], pssSig); err != nil {
		t.Fatal(err)
	}
	pssOptions := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: stdcrypto.SHA256}
	pssSig, err = RSASignPSSWithOptions(priv, stdcrypto.SHA256, digest[:], WithRSAPSSOptions(pssOptions))
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerifyPSSWithOptions(&priv.PublicKey, stdcrypto.SHA256, digest[:], pssSig, WithRSAPSSOptions(pssOptions)); err != nil {
		t.Fatal(err)
	}
	oaepCipherText, err := RSAEncryptOAEPWithOptions(plain, &priv.PublicKey, []byte("label"), WithRSAOAEPHash(sha256.New))
	if err != nil {
		t.Fatal(err)
	}
	oaepOut, err := RSADecryptOAEPWithOptions(oaepCipherText, priv, []byte("label"), WithRSAOAEPHash(sha256.New))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(oaepOut, plain) {
		t.Fatalf("RSADecryptOAEPWithOptions() = %q", oaepOut)
	}
	quickSig, err := SignSHA256WithRSA(plain, priv)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifySHA256WithRSA(plain, quickSig, &priv.PublicKey); err != nil {
		t.Fatal(err)
	}
}

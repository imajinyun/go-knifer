package vcrypto_test

import (
	"bytes"
	stdcrypto "crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/imajinyun/go-knifer/vcrypto"
)

func TestRSAAndPEM(t *testing.T) {
	priv, err := vcrypto.GenerateRSAKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM, err := vcrypto.PublicKeyToPEM(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pub, err := vcrypto.ParseRSAPublicKeyPEM(pubPEM)
	if err != nil {
		t.Fatal(err)
	}
	parsedPriv, err := vcrypto.ParseRSAPrivateKeyPEM(vcrypto.PrivateKeyToPEM(priv))
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte("rsa message")
	cipherText, err := vcrypto.RSAEncryptOAEP(plain, pub, nil)
	if err != nil {
		t.Fatal(err)
	}
	out, err := vcrypto.RSADecryptOAEP(cipherText, parsedPriv, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("RSADecryptOAEP() = %q", out)
	}

	digest := sha256.Sum256(plain)
	pssSig, err := vcrypto.RSASignPSS(parsedPriv, stdcrypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	if err := vcrypto.RSAVerifyPSS(pub, stdcrypto.SHA256, digest[:], pssSig); err != nil {
		t.Fatal(err)
	}
	pssOptions := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: stdcrypto.SHA256}
	pssSig, err = vcrypto.RSASignPSSWithOptions(parsedPriv, stdcrypto.SHA256, digest[:], vcrypto.WithRSAPSSOptions(pssOptions))
	if err != nil {
		t.Fatal(err)
	}
	if err := vcrypto.RSAVerifyPSSWithOptions(pub, stdcrypto.SHA256, digest[:], pssSig, vcrypto.WithRSAPSSOptions(pssOptions)); err != nil {
		t.Fatal(err)
	}
	oaepCipherText, err := vcrypto.RSAEncryptOAEPWithOptions(plain, pub, []byte("label"), vcrypto.WithRSAOAEPHash(sha256.New))
	if err != nil {
		t.Fatal(err)
	}
	oaepOut, err := vcrypto.RSADecryptOAEPWithOptions(oaepCipherText, parsedPriv, []byte("label"), vcrypto.WithRSAOAEPHash(sha256.New))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(oaepOut, plain) {
		t.Fatalf("RSADecryptOAEPWithOptions() = %q", oaepOut)
	}
	quickSig, err := vcrypto.SignSHA256WithRSA(plain, parsedPriv)
	if err != nil {
		t.Fatal(err)
	}
	if err := vcrypto.VerifySHA256WithRSA(plain, quickSig, pub); err != nil {
		t.Fatal(err)
	}
	pkcs8, err := vcrypto.PrivateKeyToPKCS8PEM(parsedPriv)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := vcrypto.ParseRSAPrivateKeyPEM(pkcs8); err != nil {
		t.Fatal(err)
	}
	if _, err := vcrypto.ParseRSAPublicKeyPEM(vcrypto.PublicKeyToPKCS1PEM(pub)); err != nil {
		t.Fatal(err)
	}
}

func TestAdditionalPEMCertificateAndRSAErrors(t *testing.T) {
	priv, err := vcrypto.GenerateRSAKey(1024)
	if err != nil {
		t.Fatal(err)
	}
	pub := &priv.PublicKey
	if _, err := vcrypto.ParseRSAPrivateKeyPEM([]byte("not pem")); !errors.Is(err, vcrypto.ErrInvalidKey) {
		t.Fatalf("ParseRSAPrivateKeyPEM invalid = %v", err)
	}
	if _, err := vcrypto.ParseRSAPublicKeyPEM([]byte("not pem")); !errors.Is(err, vcrypto.ErrInvalidKey) {
		t.Fatalf("ParseRSAPublicKeyPEM invalid = %v", err)
	}
	if _, err := vcrypto.ParseX509CertificatePEM([]byte("not pem")); !errors.Is(err, vcrypto.ErrInvalidKey) {
		t.Fatalf("ParseX509CertificatePEM invalid = %v", err)
	}

	certDER, err := x509.CreateCertificate(bytes.NewReader(bytes.Repeat([]byte{0x42}, 1024)), &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "go-knifer.test"},
		NotBefore:    time.Unix(0, 0),
		NotAfter:     time.Unix(3600, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}, &x509.Certificate{}, pub, priv)
	if err != nil {
		t.Fatalf("CreateCertificate: %v", err)
	}
	certPEM := pemEncode("CERTIFICATE", certDER)
	cert, err := vcrypto.ParseX509CertificatePEM(certPEM)
	if err != nil {
		t.Fatalf("ParseX509CertificatePEM: %v", err)
	}
	if cert.Subject.CommonName != "go-knifer.test" {
		t.Fatalf("certificate CN = %q", cert.Subject.CommonName)
	}
	certPub, err := vcrypto.PublicKeyFromCertificatePEM(certPEM)
	if err != nil {
		t.Fatalf("PublicKeyFromCertificatePEM: %v", err)
	}
	if certPub.N.Cmp(pub.N) != 0 {
		t.Fatal("PublicKeyFromCertificatePEM returned different key")
	}

	digest := sha256.Sum256([]byte("payload"))
	sig, err := vcrypto.RSASignPSSWithOptions(priv, stdcrypto.SHA256, digest[:], vcrypto.WithRSARandomReader(bytes.NewReader(bytes.Repeat([]byte{0x55}, 512))))
	if err != nil {
		t.Fatalf("RSASignPSSWithOptions: %v", err)
	}
	if err := vcrypto.RSAVerifyPSS(pub, stdcrypto.SHA256, digest[:], append([]byte(nil), sig[:len(sig)-1]...)); err == nil {
		t.Fatal("RSAVerifyPSS tampered signature error = nil")
	}
	data := []byte("rsa digest payload")
	digestSig, err := vcrypto.SignWithRSAOptions(data, priv, vcrypto.WithRSADigestHash(stdcrypto.SHA256, sha256.New), vcrypto.WithRSADigestRandomReader(bytes.NewReader(bytes.Repeat([]byte{0x66}, 512))))
	if err != nil {
		t.Fatalf("SignWithRSAOptions: %v", err)
	}
	if err := vcrypto.VerifyWithRSAOptions(data, digestSig, pub, vcrypto.WithRSADigestHash(stdcrypto.SHA256, sha256.New)); err != nil {
		t.Fatalf("VerifyWithRSAOptions: %v", err)
	}
	if err := vcrypto.VerifyWithRSAOptions([]byte("different"), digestSig, pub); err == nil {
		t.Fatal("VerifyWithRSAOptions tampered data error = nil")
	}
}

func pemEncode(typ string, der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der})
}

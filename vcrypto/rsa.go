package vcrypto

import (
	stdcrypto "crypto"
	"crypto/rsa"
	"hash"
	"io"

	cryptoimpl "github.com/imajinyun/go-knifer/internal/crypto"
)

// RSAOption customizes RSA helper behavior.
type RSAOption = cryptoimpl.RSAOption

// WithRSARandomReader sets the entropy source used by RSA helpers.
func WithRSARandomReader(reader io.Reader) RSAOption { return cryptoimpl.WithRSARandomReader(reader) }

// WithRSAOAEPHash sets the hash function used by RSA-OAEP helpers.
func WithRSAOAEPHash(newHash func() hash.Hash) RSAOption { return cryptoimpl.WithRSAOAEPHash(newHash) }

// WithRSAPSSOptions sets the PSS options used by RSA-PSS helpers.
func WithRSAPSSOptions(opts *rsa.PSSOptions) RSAOption { return cryptoimpl.WithRSAPSSOptions(opts) }

// GenerateRSAKey generates an RSA private key.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) { return cryptoimpl.GenerateRSAKey(bits) }

// GenerateRSAKeyWithOptions generates an RSA private key with options.
func GenerateRSAKeyWithOptions(bits int, opts ...RSAOption) (*rsa.PrivateKey, error) {
	return cryptoimpl.GenerateRSAKeyWithOptions(bits, opts...)
}

// RSAEncryptOAEP encrypts data using RSA-OAEP with SHA-256.
func RSAEncryptOAEP(plain []byte, pub *rsa.PublicKey, label []byte) ([]byte, error) {
	return cryptoimpl.RSAEncryptOAEP(plain, pub, label)
}

// RSAEncryptOAEPWithOptions encrypts data using RSA-OAEP with options.
func RSAEncryptOAEPWithOptions(plain []byte, pub *rsa.PublicKey, label []byte, opts ...RSAOption) ([]byte, error) {
	return cryptoimpl.RSAEncryptOAEPWithOptions(plain, pub, label, opts...)
}

// RSADecryptOAEP decrypts data using RSA-OAEP with SHA-256.
func RSADecryptOAEP(cipherText []byte, priv *rsa.PrivateKey, label []byte) ([]byte, error) {
	return cryptoimpl.RSADecryptOAEP(cipherText, priv, label)
}

// RSADecryptOAEPWithOptions decrypts data using RSA-OAEP with options.
func RSADecryptOAEPWithOptions(cipherText []byte, priv *rsa.PrivateKey, label []byte, opts ...RSAOption) ([]byte, error) {
	return cryptoimpl.RSADecryptOAEPWithOptions(cipherText, priv, label, opts...)
}

// RSAEncryptPKCS1v15 encrypts data using RSA PKCS#1 v1.5 padding.
func RSAEncryptPKCS1v15(plain []byte, pub *rsa.PublicKey) ([]byte, error) {
	return cryptoimpl.RSAEncryptPKCS1v15(plain, pub)
}

// RSAEncryptPKCS1v15WithOptions encrypts data using RSA PKCS#1 v1.5 padding with options.
func RSAEncryptPKCS1v15WithOptions(plain []byte, pub *rsa.PublicKey, opts ...RSAOption) ([]byte, error) {
	return cryptoimpl.RSAEncryptPKCS1v15WithOptions(plain, pub, opts...)
}

// RSADecryptPKCS1v15 decrypts data using RSA PKCS#1 v1.5 padding.
func RSADecryptPKCS1v15(cipherText []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return cryptoimpl.RSADecryptPKCS1v15(cipherText, priv)
}

// RSADecryptPKCS1v15WithOptions decrypts data using RSA PKCS#1 v1.5 padding with options.
func RSADecryptPKCS1v15WithOptions(cipherText []byte, priv *rsa.PrivateKey, opts ...RSAOption) ([]byte, error) {
	return cryptoimpl.RSADecryptPKCS1v15WithOptions(cipherText, priv, opts...)
}

// RSASignPKCS1v15 signs digest using RSA PKCS#1 v1.5.
func RSASignPKCS1v15(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte) ([]byte, error) {
	return cryptoimpl.RSASignPKCS1v15(priv, hash, digest)
}

// RSASignPKCS1v15WithOptions signs digest using RSA PKCS#1 v1.5 with options.
func RSASignPKCS1v15WithOptions(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte, opts ...RSAOption) ([]byte, error) {
	return cryptoimpl.RSASignPKCS1v15WithOptions(priv, hash, digest, opts...)
}

// RSAVerifyPKCS1v15 verifies an RSA PKCS#1 v1.5 signature.
func RSAVerifyPKCS1v15(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte) error {
	return cryptoimpl.RSAVerifyPKCS1v15(pub, hash, digest, sig)
}

// RSASignPSS signs digest using RSA-PSS.
func RSASignPSS(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte) ([]byte, error) {
	return cryptoimpl.RSASignPSS(priv, hash, digest)
}

// RSASignPSSWithOptions signs digest using RSA-PSS with options.
func RSASignPSSWithOptions(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte, opts ...RSAOption) ([]byte, error) {
	return cryptoimpl.RSASignPSSWithOptions(priv, hash, digest, opts...)
}

// RSAVerifyPSS verifies an RSA-PSS signature.
func RSAVerifyPSS(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte) error {
	return cryptoimpl.RSAVerifyPSS(pub, hash, digest, sig)
}

// RSAVerifyPSSWithOptions verifies an RSA-PSS signature with options.
func RSAVerifyPSSWithOptions(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte, opts ...RSAOption) error {
	return cryptoimpl.RSAVerifyPSSWithOptions(pub, hash, digest, sig, opts...)
}

// SignSHA256WithRSA signs data using SHA256withRSA.
func SignSHA256WithRSA(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return cryptoimpl.SignSHA256WithRSA(data, priv)
}

// VerifySHA256WithRSA verifies SHA256withRSA signature.
func VerifySHA256WithRSA(data, sig []byte, pub *rsa.PublicKey) error {
	return cryptoimpl.VerifySHA256WithRSA(data, sig, pub)
}

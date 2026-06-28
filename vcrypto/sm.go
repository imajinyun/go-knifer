package vcrypto

import (
	"hash"
	"io"

	cryptoimpl "github.com/imajinyun/knifer-go/internal/crypto"
)

// SM4Option customizes SM4 helper behavior.
type SM4Option = cryptoimpl.SM4Option

// SM2Option customizes SM2 helper behavior.
type SM2Option = cryptoimpl.SM2Option

// SM2PrivateKey is an SM2 private key.
type SM2PrivateKey = cryptoimpl.SM2PrivateKey

// SM2PublicKey is an SM2 public key.
type SM2PublicKey = cryptoimpl.SM2PublicKey

// WithSM4RandomOptions sets the entropy source options used when SM4SealGCM generates a nonce.
func WithSM4RandomOptions(opts ...RandomOption) SM4Option {
	return cryptoimpl.WithSM4RandomOptions(opts...)
}

// WithSM2RandomReader sets the entropy source used by SM2 encryption, signing, and key generation.
func WithSM2RandomReader(reader io.Reader) SM2Option {
	return cryptoimpl.WithSM2RandomReader(reader)
}

// WithSM2UID sets the SM2 user ID used by SM2Sign and SM2Verify.
func WithSM2UID(uid []byte) SM2Option { return cryptoimpl.WithSM2UID(uid) }

// SM3 returns the SM3 digest bytes of data.
func SM3(data []byte) []byte { return cryptoimpl.SM3(data) }

// SM3Hex returns the SM3 digest of data in lower-case hex form.
func SM3Hex(data []byte) string { return cryptoimpl.SM3Hex(data) }

// SM3Equal compares two SM3 digest values in constant time.
func SM3Equal(a, b []byte) bool { return cryptoimpl.SM3Equal(a, b) }

// HMACSM3Bytes returns HMAC-SM3 digest bytes.
func HMACSM3Bytes(key, data []byte) []byte { return cryptoimpl.HMACSM3Bytes(key, data) }

// HMACSM3Hex returns HMAC-SM3 in lower-case hex form.
func HMACSM3Hex(key, data []byte) string { return cryptoimpl.HMACSM3Hex(key, data) }

// SM3New returns a new SM3 hash.
func SM3New() hash.Hash { return cryptoimpl.SM3New() }

// GenSM4Key returns a random 16-byte SM4 key.
func GenSM4Key() ([]byte, error) { return cryptoimpl.GenSM4Key() }

// GenSM4KeyWithOptions returns a random SM4 key using custom random options.
func GenSM4KeyWithOptions(opts ...RandomOption) ([]byte, error) {
	return cryptoimpl.GenSM4KeyWithOptions(opts...)
}

// SM4EncryptECB encrypts plain data using SM4-ECB with PKCS#7 padding.
func SM4EncryptECB(plain, key []byte) ([]byte, error) {
	return cryptoimpl.SM4EncryptECB(plain, key)
}

// SM4DecryptECB decrypts SM4-ECB data with PKCS#7 padding.
func SM4DecryptECB(cipherText, key []byte) ([]byte, error) {
	return cryptoimpl.SM4DecryptECB(cipherText, key)
}

// SM4EncryptCBC encrypts plain data using SM4-CBC with PKCS#7 padding.
func SM4EncryptCBC(plain, key, iv []byte) ([]byte, error) {
	return cryptoimpl.SM4EncryptCBC(plain, key, iv)
}

// SM4DecryptCBC decrypts SM4-CBC data with PKCS#7 padding.
func SM4DecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	return cryptoimpl.SM4DecryptCBC(cipherText, key, iv)
}

// SM4SealGCM encrypts plain data using SM4-GCM and a freshly generated nonce.
func SM4SealGCM(plain, key, additionalData []byte) (nonce, cipherText []byte, err error) {
	return cryptoimpl.SM4SealGCM(plain, key, additionalData)
}

// SM4SealGCMWithOptions encrypts plain data using SM4-GCM and a freshly generated nonce.
func SM4SealGCMWithOptions(plain, key, additionalData []byte, opts ...SM4Option) (nonce, cipherText []byte, err error) {
	return cryptoimpl.SM4SealGCMWithOptions(plain, key, additionalData, opts...)
}

// SM4EncryptGCM encrypts plain data using SM4-GCM.
func SM4EncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.SM4EncryptGCM(plain, key, nonce, additionalData)
}

// SM4DecryptGCM decrypts SM4-GCM data.
func SM4DecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.SM4DecryptGCM(cipherText, key, nonce, additionalData)
}

// GenSM2Key generates an SM2 private key.
func GenSM2Key() (*SM2PrivateKey, error) { return cryptoimpl.GenSM2Key() }

// GenSM2KeyWithOptions generates an SM2 private key with options.
func GenSM2KeyWithOptions(opts ...SM2Option) (*SM2PrivateKey, error) {
	return cryptoimpl.GenSM2KeyWithOptions(opts...)
}

// SM2Encrypt encrypts plain data using SM2 public key encryption.
func SM2Encrypt(plain []byte, pub *SM2PublicKey) ([]byte, error) {
	return cryptoimpl.SM2Encrypt(plain, pub)
}

// SM2EncryptWithOptions encrypts plain data using SM2 public key encryption with options.
func SM2EncryptWithOptions(plain []byte, pub *SM2PublicKey, opts ...SM2Option) ([]byte, error) {
	return cryptoimpl.SM2EncryptWithOptions(plain, pub, opts...)
}

// SM2Decrypt decrypts SM2 ciphertext.
func SM2Decrypt(cipherText []byte, priv *SM2PrivateKey) ([]byte, error) {
	return cryptoimpl.SM2Decrypt(cipherText, priv)
}

// SM2Sign signs data using SM2 with the default user ID.
func SM2Sign(data []byte, priv *SM2PrivateKey) ([]byte, error) {
	return cryptoimpl.SM2Sign(data, priv)
}

// SM2SignWithOptions signs data using SM2 with options.
func SM2SignWithOptions(data []byte, priv *SM2PrivateKey, opts ...SM2Option) ([]byte, error) {
	return cryptoimpl.SM2SignWithOptions(data, priv, opts...)
}

// SM2Verify verifies an SM2 signature using the default user ID.
func SM2Verify(data, sig []byte, pub *SM2PublicKey) error {
	return cryptoimpl.SM2Verify(data, sig, pub)
}

// SM2VerifyWithOptions verifies an SM2 signature with options.
func SM2VerifyWithOptions(data, sig []byte, pub *SM2PublicKey, opts ...SM2Option) error {
	return cryptoimpl.SM2VerifyWithOptions(data, sig, pub, opts...)
}

// SM2PrivateKeyToPEM encodes an SM2 private key as PKCS#8 PEM.
func SM2PrivateKeyToPEM(priv *SM2PrivateKey) ([]byte, error) {
	return cryptoimpl.SM2PrivateKeyToPEM(priv)
}

// SM2PublicKeyToPEM encodes an SM2 public key as PKIX PEM.
func SM2PublicKeyToPEM(pub *SM2PublicKey) ([]byte, error) {
	return cryptoimpl.SM2PublicKeyToPEM(pub)
}

// ParseSM2PrivateKeyPEM parses a PKCS#8 or SEC1 SM2 private key PEM.
func ParseSM2PrivateKeyPEM(data []byte) (*SM2PrivateKey, error) {
	return cryptoimpl.ParseSM2PrivateKeyPEM(data)
}

// ParseSM2PublicKeyPEM parses a PKIX SM2 public key PEM.
func ParseSM2PublicKeyPEM(data []byte) (*SM2PublicKey, error) {
	return cryptoimpl.ParseSM2PublicKeyPEM(data)
}

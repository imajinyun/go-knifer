package vcrypto

import (
	"crypto/cipher"

	cryptoimpl "github.com/imajinyun/go-knifer/internal/crypto"
)

// AESBlockOption customizes AES block-mode helpers per call.
type AESBlockOption = cryptoimpl.AESBlockOption

// WithAESBlockFactory sets the cipher block factory used by AES helpers.
func WithAESBlockFactory(factory func([]byte) (cipher.Block, error)) AESBlockOption {
	return cryptoimpl.WithAESBlockFactory(factory)
}

// AESGCMOption customizes AES-GCM helper behavior.
type AESGCMOption = cryptoimpl.AESGCMOption

// WithGCMNonceSize sets a custom nonce size for AES-GCM helpers.
func WithGCMNonceSize(size int) AESGCMOption { return cryptoimpl.WithGCMNonceSize(size) }

// WithGCMTagSize sets a custom tag size for AES-GCM helpers.
func WithGCMTagSize(size int) AESGCMOption { return cryptoimpl.WithGCMTagSize(size) }

// WithGCMBlockFactory sets the cipher block factory used by AES-GCM helpers.
func WithGCMBlockFactory(factory func([]byte) (cipher.Block, error)) AESGCMOption {
	return cryptoimpl.WithGCMBlockFactory(factory)
}

// AESEncryptCTR encrypts or decrypts data using AES-CTR.
func AESEncryptCTR(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCTR(data, key, iv)
}

// AESEncryptCTRWithOptions encrypts or decrypts data using AES-CTR with options.
func AESEncryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESEncryptCTRWithOptions(data, key, iv, opts...)
}

// AESDecryptCTR decrypts or encrypts data using AES-CTR.
func AESDecryptCTR(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCTR(data, key, iv)
}

// AESDecryptCTRWithOptions decrypts or encrypts data using AES-CTR with options.
func AESDecryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESDecryptCTRWithOptions(data, key, iv, opts...)
}

// AESEncryptCFB encrypts data using AES-CFB.
func AESEncryptCFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCFB(data, key, iv)
}

// AESEncryptCFBWithOptions encrypts data using AES-CFB with options.
func AESEncryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESEncryptCFBWithOptions(data, key, iv, opts...)
}

// AESDecryptCFB decrypts data using AES-CFB.
func AESDecryptCFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCFB(data, key, iv)
}

// AESDecryptCFBWithOptions decrypts data using AES-CFB with options.
func AESDecryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESDecryptCFBWithOptions(data, key, iv, opts...)
}

// AESEncryptOFB encrypts or decrypts data using AES-OFB.
func AESEncryptOFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptOFB(data, key, iv)
}

// AESEncryptOFBWithOptions encrypts or decrypts data using AES-OFB with options.
func AESEncryptOFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESEncryptOFBWithOptions(data, key, iv, opts...)
}

// AESDecryptOFB decrypts or encrypts data using AES-OFB.
func AESDecryptOFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptOFB(data, key, iv)
}

// AESDecryptOFBWithOptions decrypts or encrypts data using AES-OFB with options.
func AESDecryptOFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESDecryptOFBWithOptions(data, key, iv, opts...)
}

// AESEncryptGCM encrypts plain data using AES-GCM.
func AESEncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptGCM(plain, key, nonce, additionalData)
}

// AESEncryptGCMWithOptions encrypts plain data using AES-GCM with options.
func AESEncryptGCMWithOptions(plain, key, nonce, additionalData []byte, opts ...AESGCMOption) ([]byte, error) {
	return cryptoimpl.AESEncryptGCMWithOptions(plain, key, nonce, additionalData, opts...)
}

// AESDecryptGCM decrypts AES-GCM data.
func AESDecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptGCM(cipherText, key, nonce, additionalData)
}

// AESDecryptGCMWithOptions decrypts AES-GCM data with options.
func AESDecryptGCMWithOptions(cipherText, key, nonce, additionalData []byte, opts ...AESGCMOption) ([]byte, error) {
	return cryptoimpl.AESDecryptGCMWithOptions(cipherText, key, nonce, additionalData, opts...)
}

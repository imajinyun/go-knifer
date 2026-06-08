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

// WithGCMRandomOptions sets the entropy source options used when AESSealGCM generates a nonce.
func WithGCMRandomOptions(opts ...RandomOption) AESGCMOption {
	return cryptoimpl.WithGCMRandomOptions(opts...)
}

// AESSealGCM encrypts plain data using AES-GCM and a freshly generated nonce.
// Prefer this helper for new encryption because AES-GCM authenticates both the
// ciphertext and additionalData. The returned nonce is not secret, but it must be
// stored or transmitted with the ciphertext and must never be reused with the
// same key.
func AESSealGCM(plain, key, additionalData []byte) (nonce, cipherText []byte, err error) {
	return cryptoimpl.AESSealGCM(plain, key, additionalData)
}

// AESSealGCMWithOptions encrypts plain data using AES-GCM and a freshly generated nonce.
func AESSealGCMWithOptions(plain, key, additionalData []byte, opts ...AESGCMOption) (nonce, cipherText []byte, err error) {
	return cryptoimpl.AESSealGCMWithOptions(plain, key, additionalData, opts...)
}

// AESOpenGCM decrypts data produced by AESSealGCM or AESEncryptGCM.
func AESOpenGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESOpenGCM(cipherText, key, nonce, additionalData)
}

// AESOpenGCMWithOptions decrypts AES-GCM data with options.
func AESOpenGCMWithOptions(cipherText, key, nonce, additionalData []byte, opts ...AESGCMOption) ([]byte, error) {
	return cryptoimpl.AESOpenGCMWithOptions(cipherText, key, nonce, additionalData, opts...)
}

// AESEncryptCTR encrypts or decrypts data using AES-CTR.
//
// Security: AES-CTR is an unauthenticated legacy/compatibility mode. It does not
// detect ciphertext tampering and must only be used with a separate MAC. Prefer
// AESSealGCM for new encrypted data.
func AESEncryptCTR(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCTR(data, key, iv)
}

// AESEncryptCTRWithOptions encrypts or decrypts data using AES-CTR with options.
//
// Security: AES-CTR is unauthenticated; prefer AESSealGCM for new encrypted data.
func AESEncryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESEncryptCTRWithOptions(data, key, iv, opts...)
}

// AESDecryptCTR decrypts or encrypts data using AES-CTR.
//
// Security: AES-CTR is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCTR(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCTR(data, key, iv)
}

// AESDecryptCTRWithOptions decrypts or encrypts data using AES-CTR with options.
//
// Security: AES-CTR is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESDecryptCTRWithOptions(data, key, iv, opts...)
}

// AESEncryptCFB encrypts data using AES-CFB.
//
// Security: AES-CFB is an unauthenticated legacy/compatibility mode. It does not
// detect ciphertext tampering and must only be used with a separate MAC. Prefer
// AESSealGCM for new encrypted data.
func AESEncryptCFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCFB(data, key, iv)
}

// AESEncryptCFBWithOptions encrypts data using AES-CFB with options.
//
// Security: AES-CFB is unauthenticated; prefer AESSealGCM for new encrypted data.
func AESEncryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESEncryptCFBWithOptions(data, key, iv, opts...)
}

// AESDecryptCFB decrypts data using AES-CFB.
//
// Security: AES-CFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCFB(data, key, iv)
}

// AESDecryptCFBWithOptions decrypts data using AES-CFB with options.
//
// Security: AES-CFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESDecryptCFBWithOptions(data, key, iv, opts...)
}

// AESEncryptOFB encrypts or decrypts data using AES-OFB.
//
// Security: AES-OFB is an unauthenticated legacy/compatibility mode. It does not
// detect ciphertext tampering and must only be used with a separate MAC. Prefer
// AESSealGCM for new encrypted data.
func AESEncryptOFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptOFB(data, key, iv)
}

// AESEncryptOFBWithOptions encrypts or decrypts data using AES-OFB with options.
//
// Security: AES-OFB is unauthenticated; prefer AESSealGCM for new encrypted data.
func AESEncryptOFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return cryptoimpl.AESEncryptOFBWithOptions(data, key, iv, opts...)
}

// AESDecryptOFB decrypts or encrypts data using AES-OFB.
//
// Security: AES-OFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptOFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptOFB(data, key, iv)
}

// AESDecryptOFBWithOptions decrypts or encrypts data using AES-OFB with options.
//
// Security: AES-OFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
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

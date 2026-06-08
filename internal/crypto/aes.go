package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type aesGCMConfig struct {
	nonceSize int
	tagSize   int
	blockFunc func([]byte) (cipher.Block, error)
	random    []RandomOption
}

// AESBlockOption customizes AES block-mode helpers per call.
type AESBlockOption func(*aesBlockConfig)

type aesBlockConfig struct {
	blockFunc func([]byte) (cipher.Block, error)
}

func applyAESBlockOptions(opts []AESBlockOption) aesBlockConfig {
	cfg := aesBlockConfig{blockFunc: aes.NewCipher}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.blockFunc == nil {
		cfg.blockFunc = aes.NewCipher
	}
	return cfg
}

// WithAESBlockFactory sets the cipher block factory used by AES helpers.
func WithAESBlockFactory(factory func([]byte) (cipher.Block, error)) AESBlockOption {
	return func(c *aesBlockConfig) { c.blockFunc = factory }
}

// AESGCMOption customizes AES-GCM helper behavior.
type AESGCMOption func(*aesGCMConfig)

// WithGCMNonceSize sets a custom nonce size for AES-GCM helpers.
func WithGCMNonceSize(size int) AESGCMOption {
	return func(c *aesGCMConfig) { c.nonceSize = size }
}

// WithGCMTagSize sets a custom tag size for AES-GCM helpers.
func WithGCMTagSize(size int) AESGCMOption {
	return func(c *aesGCMConfig) { c.tagSize = size }
}

// WithGCMBlockFactory sets the cipher block factory used by AES-GCM helpers.
func WithGCMBlockFactory(factory func([]byte) (cipher.Block, error)) AESGCMOption {
	return func(c *aesGCMConfig) { c.blockFunc = factory }
}

// WithGCMRandomOptions sets the entropy source options used when AESSealGCM generates a nonce.
func WithGCMRandomOptions(opts ...RandomOption) AESGCMOption {
	return func(c *aesGCMConfig) { c.random = append([]RandomOption(nil), opts...) }
}

func applyAESGCMOptions(opts []AESGCMOption) aesGCMConfig {
	cfg := aesGCMConfig{blockFunc: aes.NewCipher}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.blockFunc == nil {
		cfg.blockFunc = aes.NewCipher
	}
	return cfg
}

func newGCM(block cipher.Block, cfg aesGCMConfig) (cipher.AEAD, error) {
	if cfg.nonceSize > 0 && cfg.tagSize > 0 {
		return nil, errors.New("crypto: cannot set both GCM nonce size and tag size")
	}
	if cfg.nonceSize > 0 {
		return cipher.NewGCMWithNonceSize(block, cfg.nonceSize)
	}
	if cfg.tagSize > 0 {
		return cipher.NewGCMWithTagSize(block, cfg.tagSize)
	}
	return cipher.NewGCM(block)
}

// AESSealGCM encrypts plain data using AES-GCM and a freshly generated nonce.
// Prefer this helper for new encryption because AES-GCM authenticates both the
// ciphertext and additionalData. The returned nonce is not secret, but it must be
// stored or transmitted with the ciphertext and must never be reused with the
// same key.
func AESSealGCM(plain, key, additionalData []byte) (nonce, cipherText []byte, err error) {
	return AESSealGCMWithOptions(plain, key, additionalData)
}

// AESSealGCMWithOptions encrypts plain data using AES-GCM and a freshly generated nonce.
func AESSealGCMWithOptions(plain, key, additionalData []byte, opts ...AESGCMOption) (nonce, cipherText []byte, err error) {
	cfg := applyAESGCMOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := newGCM(block, cfg)
	if err != nil {
		return nil, nil, err
	}
	nonce, err = RandomBytesWithOptions(gcm.NonceSize(), cfg.random...)
	if err != nil {
		return nil, nil, err
	}
	return nonce, gcm.Seal(nil, nonce, plain, additionalData), nil
}

// AESOpenGCM decrypts data produced by AESSealGCM or AESEncryptGCM.
func AESOpenGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return AESDecryptGCM(cipherText, key, nonce, additionalData)
}

// AESOpenGCMWithOptions decrypts AES-GCM data with options.
func AESOpenGCMWithOptions(cipherText, key, nonce, additionalData []byte, opts ...AESGCMOption) ([]byte, error) {
	return AESDecryptGCMWithOptions(cipherText, key, nonce, additionalData, opts...)
}

// AESEncryptCTR encrypts or decrypts data using AES-CTR.
//
// Security: AES-CTR is an unauthenticated legacy/compatibility mode. It does not
// detect ciphertext tampering and must only be used with a separate MAC. Prefer
// AESSealGCM for new encrypted data.
func AESEncryptCTR(data, key, iv []byte) ([]byte, error) {
	return AESEncryptCTRWithOptions(data, key, iv)
}

// AESEncryptCTRWithOptions encrypts or decrypts data using AES-CTR with options.
//
// Security: AES-CTR is unauthenticated; prefer AESSealGCM for new encrypted data.
func AESEncryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesStream(data, key, iv, cipher.NewCTR, opts...)
}

// AESDecryptCTR decrypts or encrypts data using AES-CTR.
//
// Security: AES-CTR is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCTR(data, key, iv []byte) ([]byte, error) {
	return AESEncryptCTRWithOptions(data, key, iv)
}

// AESDecryptCTRWithOptions decrypts or encrypts data using AES-CTR with options.
//
// Security: AES-CTR is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return AESEncryptCTRWithOptions(data, key, iv, opts...)
}

// AESEncryptCFB encrypts data using AES-CFB.
//
// Security: AES-CFB is an unauthenticated legacy/compatibility mode. It does not
// detect ciphertext tampering and must only be used with a separate MAC. Prefer
// AESSealGCM for new encrypted data.
func AESEncryptCFB(data, key, iv []byte) ([]byte, error) {
	return AESEncryptCFBWithOptions(data, key, iv)
}

// AESEncryptCFBWithOptions encrypts data using AES-CFB with options.
//
// Security: AES-CFB is unauthenticated; prefer AESSealGCM for new encrypted data.
func AESEncryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesCFB(data, key, iv, false, opts...)
}

// AESDecryptCFB decrypts data using AES-CFB.
//
// Security: AES-CFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCFB(data, key, iv []byte) ([]byte, error) {
	return AESDecryptCFBWithOptions(data, key, iv)
}

// AESDecryptCFBWithOptions decrypts data using AES-CFB with options.
//
// Security: AES-CFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesCFB(data, key, iv, true, opts...)
}

// AESEncryptOFB encrypts or decrypts data using AES-OFB.
//
// Security: AES-OFB is an unauthenticated legacy/compatibility mode. It does not
// detect ciphertext tampering and must only be used with a separate MAC. Prefer
// AESSealGCM for new encrypted data.
func AESEncryptOFB(data, key, iv []byte) ([]byte, error) {
	return AESEncryptOFBWithOptions(data, key, iv)
}

// AESEncryptOFBWithOptions encrypts or decrypts data using AES-OFB with options.
//
// Security: AES-OFB is unauthenticated; prefer AESSealGCM for new encrypted data.
func AESEncryptOFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesOFB(data, key, iv, opts...)
}

// AESDecryptOFB decrypts or encrypts data using AES-OFB.
//
// Security: AES-OFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptOFB(data, key, iv []byte) ([]byte, error) {
	return AESEncryptOFBWithOptions(data, key, iv)
}

// AESDecryptOFBWithOptions decrypts or encrypts data using AES-OFB with options.
//
// Security: AES-OFB is unauthenticated; prefer AESOpenGCM for new encrypted data.
func AESDecryptOFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return AESEncryptOFBWithOptions(data, key, iv, opts...)
}

// AESEncryptGCM encrypts plain data using AES-GCM.
func AESEncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	return AESEncryptGCMWithOptions(plain, key, nonce, additionalData)
}

// AESEncryptGCMWithOptions encrypts plain data using AES-GCM with options.
func AESEncryptGCMWithOptions(plain, key, nonce, additionalData []byte, opts ...AESGCMOption) ([]byte, error) {
	cfg := applyAESGCMOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	gcm, err := newGCM(block, cfg)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, ErrInvalidIV
	}
	return gcm.Seal(nil, nonce, plain, additionalData), nil
}

// AESDecryptGCM decrypts AES-GCM data.
func AESDecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return AESDecryptGCMWithOptions(cipherText, key, nonce, additionalData)
}

// AESDecryptGCMWithOptions decrypts AES-GCM data with options.
func AESDecryptGCMWithOptions(cipherText, key, nonce, additionalData []byte, opts ...AESGCMOption) ([]byte, error) {
	cfg := applyAESGCMOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	gcm, err := newGCM(block, cfg)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, ErrInvalidIV
	}
	return gcm.Open(nil, nonce, cipherText, additionalData)
}

func aesBlockWithIV(key, iv []byte, opts ...AESBlockOption) (cipher.Block, error) {
	cfg := applyAESBlockOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIV
	}
	return block, nil
}

func aesStream(data, key, iv []byte, newStream func(cipher.Block, []byte) cipher.Stream, opts ...AESBlockOption) ([]byte, error) {
	block, err := aesBlockWithIV(key, iv, opts...)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	newStream(block, iv).XORKeyStream(out, data)
	return out, nil
}

func aesCFB(data, key, iv []byte, decrypt bool, opts ...AESBlockOption) ([]byte, error) {
	block, err := aesBlockWithIV(key, iv, opts...)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	feedback := append([]byte(nil), iv...)
	stream := make([]byte, block.BlockSize())
	for pos := 0; pos < len(data); pos += block.BlockSize() {
		block.Encrypt(stream, feedback)
		n := min(block.BlockSize(), len(data)-pos)
		for i := 0; i < n; i++ {
			out[pos+i] = data[pos+i] ^ stream[i]
		}
		if n == block.BlockSize() {
			if decrypt {
				copy(feedback, data[pos:pos+n])
			} else {
				copy(feedback, out[pos:pos+n])
			}
		}
	}
	return out, nil
}

func aesOFB(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	block, err := aesBlockWithIV(key, iv, opts...)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	feedback := append([]byte(nil), iv...)
	for pos := 0; pos < len(data); pos += block.BlockSize() {
		block.Encrypt(feedback, feedback)
		n := min(block.BlockSize(), len(data)-pos)
		for i := 0; i < n; i++ {
			out[pos+i] = data[pos+i] ^ feedback[i]
		}
	}
	return out, nil
}

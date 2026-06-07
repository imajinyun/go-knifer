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

// AESEncryptCBC encrypts plain data using AES-CBC with PKCS#7 padding.
func AESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	return AESEncryptCBCWithOptions(plain, key, iv)
}

// AESEncryptCBCWithOptions encrypts plain data using AES-CBC with options.
func AESEncryptCBCWithOptions(plain, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	cfg := applyAESBlockOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	return encryptCBC(block, plain, iv)
}

// AESDecryptCBC decrypts AES-CBC data using PKCS#7 padding.
func AESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	return AESDecryptCBCWithOptions(cipherText, key, iv)
}

// AESDecryptCBCWithOptions decrypts AES-CBC data with options.
func AESDecryptCBCWithOptions(cipherText, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	cfg := applyAESBlockOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	return decryptCBC(block, cipherText, iv)
}

// AESEncryptECB encrypts plain data using AES-ECB with PKCS#7 padding.
func AESEncryptECB(plain, key []byte) ([]byte, error) {
	return AESEncryptECBWithOptions(plain, key)
}

// AESEncryptECBWithOptions encrypts plain data using AES-ECB with options.
func AESEncryptECBWithOptions(plain, key []byte, opts ...AESBlockOption) ([]byte, error) {
	cfg := applyAESBlockOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	return encryptECB(block, plain), nil
}

// AESDecryptECB decrypts AES-ECB data using PKCS#7 padding.
func AESDecryptECB(cipherText, key []byte) ([]byte, error) {
	return AESDecryptECBWithOptions(cipherText, key)
}

// AESDecryptECBWithOptions decrypts AES-ECB data with options.
func AESDecryptECBWithOptions(cipherText, key []byte, opts ...AESBlockOption) ([]byte, error) {
	cfg := applyAESBlockOptions(opts)
	block, err := cfg.blockFunc(key)
	if err != nil {
		return nil, err
	}
	return decryptECB(block, cipherText)
}

// AESEncryptCTR encrypts or decrypts data using AES-CTR.
func AESEncryptCTR(data, key, iv []byte) ([]byte, error) {
	return AESEncryptCTRWithOptions(data, key, iv)
}

// AESEncryptCTRWithOptions encrypts or decrypts data using AES-CTR with options.
func AESEncryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesStream(data, key, iv, cipher.NewCTR, opts...)
}

// AESDecryptCTR decrypts or encrypts data using AES-CTR.
func AESDecryptCTR(data, key, iv []byte) ([]byte, error) {
	return AESEncryptCTRWithOptions(data, key, iv)
}

// AESDecryptCTRWithOptions decrypts or encrypts data using AES-CTR with options.
func AESDecryptCTRWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return AESEncryptCTRWithOptions(data, key, iv, opts...)
}

// AESEncryptCFB encrypts data using AES-CFB.
func AESEncryptCFB(data, key, iv []byte) ([]byte, error) {
	return AESEncryptCFBWithOptions(data, key, iv)
}

// AESEncryptCFBWithOptions encrypts data using AES-CFB with options.
func AESEncryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesCFB(data, key, iv, false, opts...)
}

// AESDecryptCFB decrypts data using AES-CFB.
func AESDecryptCFB(data, key, iv []byte) ([]byte, error) {
	return AESDecryptCFBWithOptions(data, key, iv)
}

// AESDecryptCFBWithOptions decrypts data using AES-CFB with options.
func AESDecryptCFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesCFB(data, key, iv, true, opts...)
}

// AESEncryptOFB encrypts or decrypts data using AES-OFB.
func AESEncryptOFB(data, key, iv []byte) ([]byte, error) {
	return AESEncryptOFBWithOptions(data, key, iv)
}

// AESEncryptOFBWithOptions encrypts or decrypts data using AES-OFB with options.
func AESEncryptOFBWithOptions(data, key, iv []byte, opts ...AESBlockOption) ([]byte, error) {
	return aesOFB(data, key, iv, opts...)
}

// AESDecryptOFB decrypts or encrypts data using AES-OFB.
func AESDecryptOFB(data, key, iv []byte) ([]byte, error) {
	return AESEncryptOFBWithOptions(data, key, iv)
}

// AESDecryptOFBWithOptions decrypts or encrypts data using AES-OFB with options.
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

func encryptCBC(block cipher.Block, plain, iv []byte) ([]byte, error) {
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIV
	}
	plain = pkcs7Pad(plain, block.BlockSize())
	out := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, plain)
	return out, nil
}

func decryptCBC(block cipher.Block, cipherText, iv []byte) ([]byte, error) {
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIV
	}
	if len(cipherText) == 0 || len(cipherText)%block.BlockSize() != 0 {
		return nil, ErrInvalidCipherText
	}
	out := make([]byte, len(cipherText))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, cipherText)
	return pkcs7Unpad(out, block.BlockSize())
}

func encryptECB(block cipher.Block, plain []byte) []byte {
	plain = pkcs7Pad(plain, block.BlockSize())
	out := make([]byte, len(plain))
	for bs, be := 0, block.BlockSize(); bs < len(plain); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(out[bs:be], plain[bs:be])
	}
	return out
}

func decryptECB(block cipher.Block, cipherText []byte) ([]byte, error) {
	if len(cipherText) == 0 || len(cipherText)%block.BlockSize() != 0 {
		return nil, ErrInvalidCipherText
	}
	out := make([]byte, len(cipherText))
	for bs, be := 0, block.BlockSize(); bs < len(cipherText); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(out[bs:be], cipherText[bs:be])
	}
	return pkcs7Unpad(out, block.BlockSize())
}

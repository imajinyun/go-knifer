package crypto

import (
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	"slices"

	gmcipher "github.com/emmansun/gmsm/cipher"
	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/sm3"
	"github.com/emmansun/gmsm/sm4"
	"github.com/emmansun/gmsm/smx509"
	knifer "github.com/imajinyun/knifer-go"
)

// SM2PrivateKey is an SM2 private key.
type SM2PrivateKey = sm2.PrivateKey

// SM2PublicKey is an SM2 public key.
type SM2PublicKey = ecdsa.PublicKey

type sm4Config struct {
	random []RandomOption
}

// SM4Option customizes SM4 helper behavior.
type SM4Option func(*sm4Config)

// WithSM4RandomOptions sets the entropy source options used when SM4SealGCM generates a nonce.
func WithSM4RandomOptions(opts ...RandomOption) SM4Option {
	return func(c *sm4Config) { c.random = slices.Clone(opts) }
}

func applySM4Options(opts []SM4Option) sm4Config {
	var cfg sm4Config
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

type sm2Config struct {
	random io.Reader
	uid    []byte
}

// SM2Option customizes SM2 helper behavior.
type SM2Option func(*sm2Config)

// WithSM2RandomReader sets the entropy source used by SM2 encryption, signing, and key generation.
func WithSM2RandomReader(reader io.Reader) SM2Option {
	return func(c *sm2Config) {
		if reader != nil {
			c.random = reader
		}
	}
}

// WithSM2UID sets the SM2 user ID used by SM2Sign and SM2Verify.
func WithSM2UID(uid []byte) SM2Option {
	return func(c *sm2Config) { c.uid = slices.Clone(uid) }
}

func applySM2Options(opts []SM2Option) sm2Config {
	cfg := sm2Config{random: rand.Reader}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.random == nil {
		cfg.random = rand.Reader
	}
	return cfg
}

// SM3 returns the SM3 digest bytes of data.
func SM3(data []byte) []byte {
	sum := sm3.Sum(data)
	return sum[:]
}

// SM3Hex returns the SM3 digest of data in lower-case hex form.
func SM3Hex(data []byte) string {
	sum := sm3.Sum(data)
	return hex.EncodeToString(sum[:])
}

// SM3Equal compares two SM3 digest values in constant time.
func SM3Equal(a, b []byte) bool { return hmac.Equal(a, b) }

// HMACSM3Bytes returns HMAC-SM3 digest bytes.
func HMACSM3Bytes(key, data []byte) []byte { return HMACBytes(sm3.New, key, data) }

// HMACSM3Hex returns HMAC-SM3 in lower-case hex form.
func HMACSM3Hex(key, data []byte) string { return HMACHex(sm3.New, key, data) }

// SM3New returns a new SM3 hash.
func SM3New() hash.Hash { return sm3.New() }

// GenSM4Key returns a random 16-byte SM4 key.
func GenSM4Key() ([]byte, error) { return GenSM4KeyWithOptions() }

// GenSM4KeyWithOptions returns a random SM4 key using custom random options.
func GenSM4KeyWithOptions(opts ...RandomOption) ([]byte, error) {
	return RandomBytesWithOptions(sm4.BlockSize, opts...)
}

func newSM4Block(key []byte) (cipher.Block, error) {
	if err := ValidateSM4Key(key); err != nil {
		return nil, err
	}
	return sm4.NewCipher(key)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	out := make([]byte, len(data)+padding)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(padding)
	}
	return out
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, ErrInvalidCipherText
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize || padding > len(data) {
		return nil, ErrInvalidCipherText
	}
	for _, b := range data[len(data)-padding:] {
		if int(b) != padding {
			return nil, ErrInvalidCipherText
		}
	}
	return slices.Clone(data[:len(data)-padding]), nil
}

// SM4EncryptECB encrypts plain data using SM4-ECB with PKCS#7 padding.
func SM4EncryptECB(plain, key []byte) ([]byte, error) {
	block, err := newSM4Block(key)
	if err != nil {
		return nil, err
	}
	plain = pkcs7Pad(plain, block.BlockSize())
	out := make([]byte, len(plain))
	encrypter := gmcipher.NewECBEncrypter(block)
	encrypter.CryptBlocks(out, plain)
	return out, nil
}

// SM4DecryptECB decrypts SM4-ECB data with PKCS#7 padding.
func SM4DecryptECB(cipherText, key []byte) ([]byte, error) {
	block, err := newSM4Block(key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%block.BlockSize() != 0 {
		return nil, ErrInvalidCipherText
	}
	out := make([]byte, len(cipherText))
	decrypter := gmcipher.NewECBDecrypter(block)
	decrypter.CryptBlocks(out, cipherText)
	plain, err := pkcs7Unpad(out, block.BlockSize())
	if err != nil {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "sm4-ecb invalid padding", err)
	}
	return plain, nil
}

// SM4EncryptCBC encrypts plain data using SM4-CBC with PKCS#7 padding.
func SM4EncryptCBC(plain, key, iv []byte) ([]byte, error) {
	block, err := newSM4Block(key)
	if err != nil {
		return nil, err
	}
	if err := ValidateSM4IV(iv); err != nil {
		return nil, err
	}
	plain = pkcs7Pad(plain, block.BlockSize())
	out := make([]byte, len(plain))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(out, plain)
	return out, nil
}

// SM4DecryptCBC decrypts SM4-CBC data with PKCS#7 padding.
func SM4DecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	block, err := newSM4Block(key)
	if err != nil {
		return nil, err
	}
	if err := ValidateSM4IV(iv); err != nil {
		return nil, err
	}
	if len(cipherText) == 0 || len(cipherText)%block.BlockSize() != 0 {
		return nil, ErrInvalidCipherText
	}
	out := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(out, cipherText)
	plain, err := pkcs7Unpad(out, block.BlockSize())
	if err != nil {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "sm4-cbc invalid padding", err)
	}
	return plain, nil
}

// SM4SealGCM encrypts plain data using SM4-GCM and a freshly generated nonce.
func SM4SealGCM(plain, key, additionalData []byte) (nonce, cipherText []byte, err error) {
	return SM4SealGCMWithOptions(plain, key, additionalData)
}

// SM4SealGCMWithOptions encrypts plain data using SM4-GCM and a freshly generated nonce.
func SM4SealGCMWithOptions(plain, key, additionalData []byte, opts ...SM4Option) (nonce, cipherText []byte, err error) {
	cfg := applySM4Options(opts)
	block, err := newSM4Block(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce, err = RandomBytesWithOptions(gcm.NonceSize(), cfg.random...)
	if err != nil {
		return nil, nil, err
	}
	return nonce, gcm.Seal(nil, nonce, plain, additionalData), nil
}

// SM4EncryptGCM encrypts plain data using SM4-GCM.
func SM4EncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	block, err := newSM4Block(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, ErrInvalidIV
	}
	return gcm.Seal(nil, nonce, plain, additionalData), nil
}

// SM4DecryptGCM decrypts SM4-GCM data.
func SM4DecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	block, err := newSM4Block(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, ErrInvalidIV
	}
	plain, err := gcm.Open(nil, nonce, cipherText, additionalData)
	if err != nil {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "sm4-gcm authentication failed", errors.Join(ErrInvalidCipherText, err))
	}
	return plain, nil
}

// GenSM2Key generates an SM2 private key.
func GenSM2Key() (*SM2PrivateKey, error) { return GenSM2KeyWithOptions() }

// GenSM2KeyWithOptions generates an SM2 private key with options.
func GenSM2KeyWithOptions(opts ...SM2Option) (*SM2PrivateKey, error) {
	cfg := applySM2Options(opts)
	return sm2.GenerateKey(cfg.random)
}

// SM2Encrypt encrypts plain data using SM2 public key encryption.
func SM2Encrypt(plain []byte, pub *SM2PublicKey) ([]byte, error) {
	return SM2EncryptWithOptions(plain, pub)
}

// SM2EncryptWithOptions encrypts plain data using SM2 public key encryption with options.
func SM2EncryptWithOptions(plain []byte, pub *SM2PublicKey, opts ...SM2Option) ([]byte, error) {
	if pub == nil {
		return nil, ErrInvalidKey
	}
	cfg := applySM2Options(opts)
	return sm2.Encrypt(cfg.random, pub, plain, sm2.ASN1EncrypterOpts)
}

// SM2Decrypt decrypts SM2 ciphertext.
func SM2Decrypt(cipherText []byte, priv *SM2PrivateKey) ([]byte, error) {
	if priv == nil {
		return nil, ErrInvalidKey
	}
	plain, err := priv.Decrypt(nil, cipherText, sm2.ASN1DecrypterOpts)
	if err != nil {
		return nil, knifer.WrapError(knifer.ErrCodeInvalidInput, "sm2 decryption failed", errors.Join(ErrInvalidCipherText, err))
	}
	return plain, nil
}

// SM2Sign signs data using SM2 with the default user ID.
func SM2Sign(data []byte, priv *SM2PrivateKey) ([]byte, error) {
	return SM2SignWithOptions(data, priv)
}

// SM2SignWithOptions signs data using SM2 with options.
func SM2SignWithOptions(data []byte, priv *SM2PrivateKey, opts ...SM2Option) ([]byte, error) {
	if priv == nil {
		return nil, ErrInvalidKey
	}
	cfg := applySM2Options(opts)
	return priv.Sign(cfg.random, data, sm2.NewSM2SignerOption(true, cfg.uid))
}

// SM2Verify verifies an SM2 signature using the default user ID.
func SM2Verify(data, sig []byte, pub *SM2PublicKey) error {
	return SM2VerifyWithOptions(data, sig, pub)
}

// SM2VerifyWithOptions verifies an SM2 signature with options.
func SM2VerifyWithOptions(data, sig []byte, pub *SM2PublicKey, opts ...SM2Option) error {
	if pub == nil {
		return ErrInvalidKey
	}
	cfg := applySM2Options(opts)
	if sm2.VerifyASN1WithSM2(pub, cfg.uid, data, sig) {
		return nil
	}
	return ErrInvalidSM2Signature
}

// SM2PrivateKeyToPEM encodes an SM2 private key as PKCS#8 PEM.
func SM2PrivateKeyToPEM(priv *SM2PrivateKey) ([]byte, error) {
	if priv == nil {
		return nil, ErrInvalidKey
	}
	b, err := smx509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b}), nil
}

// SM2PublicKeyToPEM encodes an SM2 public key as PKIX PEM.
func SM2PublicKeyToPEM(pub *SM2PublicKey) ([]byte, error) {
	if pub == nil {
		return nil, ErrInvalidKey
	}
	b, err := smx509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b}), nil
}

// ParseSM2PrivateKeyPEM parses a PKCS#8 or SEC1 SM2 private key PEM.
func ParseSM2PrivateKeyPEM(data []byte) (*SM2PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidKey
	}
	if key, err := smx509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if priv, ok := key.(*sm2.PrivateKey); ok {
			return priv, nil
		}
		return nil, ErrInvalidKey
	}
	priv, err := smx509.ParseSM2PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// ParseSM2PublicKeyPEM parses a PKIX SM2 public key PEM.
func ParseSM2PublicKeyPEM(data []byte) (*SM2PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidKey
	}
	key, err := smx509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := key.(*ecdsa.PublicKey)
	if !ok || !sm2.IsSM2PublicKey(pub) {
		return nil, ErrInvalidKey
	}
	return pub, nil
}

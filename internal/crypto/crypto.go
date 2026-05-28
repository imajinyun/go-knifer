package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"hash"
	"io"
)

var (
	// ErrInvalidKey 表示无效的加密密钥。ErrInvalidKey indicates an invalid cryptographic key.
	ErrInvalidKey = errors.New("invalid key")
	// ErrInvalidIV 表示无效的初始化向量。ErrInvalidIV indicates an invalid initialization vector.
	ErrInvalidIV = errors.New("invalid iv")
	// ErrInvalidCipherText 表示无效的密文数据。ErrInvalidCipherText indicates invalid encrypted data.
	ErrInvalidCipherText = errors.New("invalid cipher text")
)

// MD5Hex 返回数据的 MD5 小写十六进制摘要，对应 Hutool DigestUtil.md5Hex。MD5Hex returns the MD5 digest of data in lower-case hex form.
func MD5Hex(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// SHA1Hex 返回数据的 SHA1 小写十六进制摘要。SHA1Hex returns the SHA1 digest of data in lower-case hex form.
func SHA1Hex(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}

// SHA256Hex 返回数据的 SHA256 小写十六进制摘要。SHA256Hex returns the SHA256 digest of data in lower-case hex form.
func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// SHA512Hex 返回数据的 SHA512 小写十六进制摘要。SHA512Hex returns the SHA512 digest of data in lower-case hex form.
func SHA512Hex(data []byte) string {
	sum := sha512.Sum512(data)
	return hex.EncodeToString(sum[:])
}

// HMACHex 使用指定哈希算法计算 HMAC 小写十六进制摘要。HMACHex returns HMAC digest in lower-case hex form.
func HMACHex(fn func() hash.Hash, key, data []byte) string {
	h := hmac.New(fn, key)
	_, _ = h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// HMACMD5Hex 返回 HMAC-MD5 小写十六进制摘要。HMACMD5Hex returns HMAC-MD5 in lower-case hex form.
func HMACMD5Hex(key, data []byte) string { return HMACHex(md5.New, key, data) }

// HMACSHA1Hex 返回 HMAC-SHA1 小写十六进制摘要。HMACSHA1Hex returns HMAC-SHA1 in lower-case hex form.
func HMACSHA1Hex(key, data []byte) string { return HMACHex(sha1.New, key, data) }

// HMACSHA256Hex 返回 HMAC-SHA256 小写十六进制摘要。HMACSHA256Hex returns HMAC-SHA256 in lower-case hex form.
func HMACSHA256Hex(key, data []byte) string { return HMACHex(sha256.New, key, data) }

// HMACSHA512Hex 返回 HMAC-SHA512 小写十六进制摘要。HMACSHA512Hex returns HMAC-SHA512 in lower-case hex form.
func HMACSHA512Hex(key, data []byte) string { return HMACHex(sha512.New, key, data) }

// RandomBytes 返回 n 个密码学安全随机字节，对应 Hutool SecureUtil.generateRandomBytes。RandomBytes returns n cryptographically secure random bytes.
func RandomBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrInvalidKey
	}
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateAESKey 生成 AES 随机密钥，长度必须为 16、24 或 32 字节。GenerateAESKey returns a random AES key. Valid sizes are 16, 24, or 32 bytes.
func GenerateAESKey(size int) ([]byte, error) {
	if size != 16 && size != 24 && size != 32 {
		return nil, ErrInvalidKey
	}
	return RandomBytes(size)
}

// AESEncryptCBC 使用 AES-CBC 和 PKCS#7 填充加密明文数据。AESEncryptCBC encrypts plain data using AES-CBC with PKCS#7 padding.
func AESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIV
	}
	plain = pkcs7Pad(plain, block.BlockSize())
	out := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, plain)
	return out, nil
}

// AESDecryptCBC 使用 AES-CBC 和 PKCS#7 填充解密密文数据。AESDecryptCBC decrypts AES-CBC data using PKCS#7 padding.
func AESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
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

// AESEncryptGCM 使用 AES-GCM 加密明文数据。AESEncryptGCM encrypts plain data using AES-GCM.
func AESEncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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

// AESDecryptGCM 解密 AES-GCM 密文数据。AESDecryptGCM decrypts AES-GCM data.
func AESDecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
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
	return gcm.Open(nil, nonce, cipherText, additionalData)
}

// GenerateRSAKey 生成 RSA 私钥，对应 Hutool KeyUtil.generateKeyPair 的常见用法。GenerateRSAKey generates an RSA private key.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// RSAEncryptOAEP 使用 RSA-OAEP 和 SHA-256 加密数据。RSAEncryptOAEP encrypts data using RSA-OAEP with SHA-256.
func RSAEncryptOAEP(plain []byte, pub *rsa.PublicKey, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plain, label)
}

// RSADecryptOAEP 使用 RSA-OAEP 和 SHA-256 解密数据。RSADecryptOAEP decrypts data using RSA-OAEP with SHA-256.
func RSADecryptOAEP(cipherText []byte, priv *rsa.PrivateKey, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, cipherText, label)
}

// PrivateKeyToPEM 将 RSA 私钥编码为 PKCS#1 PEM。PrivateKeyToPEM encodes an RSA private key as PKCS#1 PEM.
func PrivateKeyToPEM(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
}

// PublicKeyToPEM 将 RSA 公钥编码为 PKIX PEM。PublicKeyToPEM encodes an RSA public key as PKIX PEM.
func PublicKeyToPEM(pub *rsa.PublicKey) ([]byte, error) {
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: b}), nil
}

// ParseRSAPrivateKeyPEM 解析 PKCS#1 或 PKCS#8 RSA 私钥 PEM。ParseRSAPrivateKeyPEM parses a PKCS#1 or PKCS#8 RSA private key PEM.
func ParseRSAPrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidKey
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	return priv, nil
}

// ParseRSAPublicKeyPEM 解析 PKIX 或 PKCS#1 RSA 公钥 PEM。ParseRSAPublicKeyPEM parses a PKIX or PKCS#1 RSA public key PEM.
func ParseRSAPublicKeyPEM(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidKey
	}
	if pub, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return pub, nil
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	return pub, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	return append(append([]byte(nil), data...), bytes.Repeat([]byte{byte(padding)}, padding)...)
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
	return append([]byte(nil), data[:len(data)-padding]...), nil
}

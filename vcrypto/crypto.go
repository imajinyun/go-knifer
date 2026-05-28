package vcrypto

import (
	"crypto/rsa"

	cryptoimpl "github.com/imajinyun/go-knifer/internal/crypto"
)

// ErrInvalidKey 表示无效的加密密钥。ErrInvalidKey indicates an invalid cryptographic key.
var ErrInvalidKey = cryptoimpl.ErrInvalidKey

// ErrInvalidIV 表示无效的初始化向量。ErrInvalidIV indicates an invalid initialization vector.
var ErrInvalidIV = cryptoimpl.ErrInvalidIV

// ErrInvalidCipherText 表示无效的密文数据。ErrInvalidCipherText indicates invalid encrypted data.
var ErrInvalidCipherText = cryptoimpl.ErrInvalidCipherText

// MD5Hex 返回字符串的 MD5 小写十六进制摘要，对应 Hutool SecureUtil.md5。MD5Hex returns the MD5 digest of s in lower-case hex form.
func MD5Hex(s string) string { return cryptoimpl.MD5Hex([]byte(s)) }

// MD5HexBytes 返回字节切片的 MD5 小写十六进制摘要。MD5HexBytes returns the MD5 digest of data in lower-case hex form.
func MD5HexBytes(data []byte) string { return cryptoimpl.MD5Hex(data) }

// SHA1Hex 返回字符串的 SHA1 小写十六进制摘要。SHA1Hex returns the SHA1 digest of s in lower-case hex form.
func SHA1Hex(s string) string { return cryptoimpl.SHA1Hex([]byte(s)) }

// SHA256Hex 返回字符串的 SHA256 小写十六进制摘要。SHA256Hex returns the SHA256 digest of s in lower-case hex form.
func SHA256Hex(s string) string { return cryptoimpl.SHA256Hex([]byte(s)) }

// SHA512Hex 返回字符串的 SHA512 小写十六进制摘要。SHA512Hex returns the SHA512 digest of s in lower-case hex form.
func SHA512Hex(s string) string { return cryptoimpl.SHA512Hex([]byte(s)) }

// HMACMD5Hex 返回 HMAC-MD5 小写十六进制摘要。HMACMD5Hex returns HMAC-MD5 in lower-case hex form.
func HMACMD5Hex(key, data []byte) string { return cryptoimpl.HMACMD5Hex(key, data) }

// HMACSHA1Hex 返回 HMAC-SHA1 小写十六进制摘要。HMACSHA1Hex returns HMAC-SHA1 in lower-case hex form.
func HMACSHA1Hex(key, data []byte) string { return cryptoimpl.HMACSHA1Hex(key, data) }

// HMACSHA256Hex 返回 HMAC-SHA256 小写十六进制摘要。HMACSHA256Hex returns HMAC-SHA256 in lower-case hex form.
func HMACSHA256Hex(key, data []byte) string { return cryptoimpl.HMACSHA256Hex(key, data) }

// HMACSHA512Hex 返回 HMAC-SHA512 小写十六进制摘要。HMACSHA512Hex returns HMAC-SHA512 in lower-case hex form.
func HMACSHA512Hex(key, data []byte) string { return cryptoimpl.HMACSHA512Hex(key, data) }

// RandomBytes 返回 n 个密码学安全随机字节，对应 Hutool SecureUtil.generateRandomBytes。RandomBytes returns n cryptographically secure random bytes.
func RandomBytes(n int) ([]byte, error) { return cryptoimpl.RandomBytes(n) }

// GenerateAESKey 生成 AES 随机密钥，长度必须为 16、24 或 32 字节。GenerateAESKey returns a random AES key. Valid sizes are 16, 24, or 32 bytes.
func GenerateAESKey(size int) ([]byte, error) { return cryptoimpl.GenerateAESKey(size) }

// AESEncryptCBC 使用 AES-CBC 和 PKCS#7 填充加密明文数据。AESEncryptCBC encrypts plain data using AES-CBC with PKCS#7 padding.
func AESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCBC(plain, key, iv)
}

// AESDecryptCBC 使用 AES-CBC 和 PKCS#7 填充解密密文数据。AESDecryptCBC decrypts AES-CBC data using PKCS#7 padding.
func AESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCBC(cipherText, key, iv)
}

// AESEncryptGCM 使用 AES-GCM 加密明文数据。AESEncryptGCM encrypts plain data using AES-GCM.
func AESEncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptGCM(plain, key, nonce, additionalData)
}

// AESDecryptGCM 解密 AES-GCM 密文数据。AESDecryptGCM decrypts AES-GCM data.
func AESDecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptGCM(cipherText, key, nonce, additionalData)
}

// GenerateRSAKey 生成 RSA 私钥，对应 Hutool KeyUtil.generateKeyPair 的常见用法。GenerateRSAKey generates an RSA private key.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) { return cryptoimpl.GenerateRSAKey(bits) }

// RSAEncryptOAEP 使用 RSA-OAEP 和 SHA-256 加密数据。RSAEncryptOAEP encrypts data using RSA-OAEP with SHA-256.
func RSAEncryptOAEP(plain []byte, pub *rsa.PublicKey, label []byte) ([]byte, error) {
	return cryptoimpl.RSAEncryptOAEP(plain, pub, label)
}

// RSADecryptOAEP 使用 RSA-OAEP 和 SHA-256 解密数据。RSADecryptOAEP decrypts data using RSA-OAEP with SHA-256.
func RSADecryptOAEP(cipherText []byte, priv *rsa.PrivateKey, label []byte) ([]byte, error) {
	return cryptoimpl.RSADecryptOAEP(cipherText, priv, label)
}

// PrivateKeyToPEM 将 RSA 私钥编码为 PKCS#1 PEM。PrivateKeyToPEM encodes an RSA private key as PKCS#1 PEM.
func PrivateKeyToPEM(priv *rsa.PrivateKey) []byte { return cryptoimpl.PrivateKeyToPEM(priv) }

// PublicKeyToPEM 将 RSA 公钥编码为 PKIX PEM。PublicKeyToPEM encodes an RSA public key as PKIX PEM.
func PublicKeyToPEM(pub *rsa.PublicKey) ([]byte, error) { return cryptoimpl.PublicKeyToPEM(pub) }

// ParseRSAPrivateKeyPEM 解析 PKCS#1 或 PKCS#8 RSA 私钥 PEM。ParseRSAPrivateKeyPEM parses a PKCS#1 or PKCS#8 RSA private key PEM.
func ParseRSAPrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	return cryptoimpl.ParseRSAPrivateKeyPEM(data)
}

// ParseRSAPublicKeyPEM 解析 PKIX 或 PKCS#1 RSA 公钥 PEM。ParseRSAPublicKeyPEM parses a PKIX or PKCS#1 RSA public key PEM.
func ParseRSAPublicKeyPEM(data []byte) (*rsa.PublicKey, error) {
	return cryptoimpl.ParseRSAPublicKeyPEM(data)
}

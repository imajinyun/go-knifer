package vcrypto

import (
	stdcrypto "crypto"
	"crypto/rsa"
	"crypto/x509"
	"hash"

	cryptoimpl "github.com/imajinyun/go-knifer/internal/crypto"
)

// ErrInvalidKey indicates an invalid cryptographic key.
var ErrInvalidKey = cryptoimpl.ErrInvalidKey

// ErrInvalidIV indicates an invalid initialization vector.
var ErrInvalidIV = cryptoimpl.ErrInvalidIV

// ErrInvalidCipherText indicates invalid encrypted data.
var ErrInvalidCipherText = cryptoimpl.ErrInvalidCipherText

// MD5Hex returns the MD5 digest of s in lower-case hex form.
func MD5Hex(s string) string { return cryptoimpl.MD5Hex([]byte(s)) }

// MD5HexBytes returns the MD5 digest of data in lower-case hex form.
func MD5HexBytes(data []byte) string { return cryptoimpl.MD5Hex(data) }

// MD5 returns the MD5 digest bytes of data.
func MD5(data []byte) []byte { return cryptoimpl.MD5(data) }

// MD5Hex16 returns the middle 16 characters of the MD5 hex digest.
func MD5Hex16(data []byte) string { return cryptoimpl.MD5Hex16(data) }

// MD5HexTo16 returns the middle 16 characters of a 32-character MD5 hex digest.
func MD5HexTo16(md5Hex string) string { return cryptoimpl.MD5HexTo16(md5Hex) }

// SHA1Hex returns the SHA1 digest of s in lower-case hex form.
func SHA1Hex(s string) string { return cryptoimpl.SHA1Hex([]byte(s)) }

// SHA1 returns the SHA1 digest bytes of data.
func SHA1(data []byte) []byte { return cryptoimpl.SHA1(data) }

// SHA1HexBytes returns the SHA1 digest of data in lower-case hex form.
func SHA1HexBytes(data []byte) string { return cryptoimpl.SHA1Hex(data) }

// SHA224 returns the SHA224 digest bytes of data.
func SHA224(data []byte) []byte { return cryptoimpl.SHA224(data) }

// SHA224Hex returns the SHA224 digest of data in lower-case hex form.
func SHA224Hex(data []byte) string { return cryptoimpl.SHA224Hex(data) }

// SHA256Hex returns the SHA256 digest of s in lower-case hex form.
func SHA256Hex(s string) string { return cryptoimpl.SHA256Hex([]byte(s)) }

// SHA256 returns the SHA256 digest bytes of data.
func SHA256(data []byte) []byte { return cryptoimpl.SHA256(data) }

// SHA256HexBytes returns the SHA256 digest of data in lower-case hex form.
func SHA256HexBytes(data []byte) string { return cryptoimpl.SHA256Hex(data) }

// SHA384 returns the SHA384 digest bytes of data.
func SHA384(data []byte) []byte { return cryptoimpl.SHA384(data) }

// SHA384Hex returns the SHA384 digest of data in lower-case hex form.
func SHA384Hex(data []byte) string { return cryptoimpl.SHA384Hex(data) }

// SHA512Hex returns the SHA512 digest of s in lower-case hex form.
func SHA512Hex(s string) string { return cryptoimpl.SHA512Hex([]byte(s)) }

// SHA512 returns the SHA512 digest bytes of data.
func SHA512(data []byte) []byte { return cryptoimpl.SHA512(data) }

// SHA512HexBytes returns the SHA512 digest of data in lower-case hex form.
func SHA512HexBytes(data []byte) string { return cryptoimpl.SHA512Hex(data) }

// HMACHex returns HMAC digest in lower-case hex form using the given hash function.
func HMACHex(fn func() hash.Hash, key, data []byte) string { return cryptoimpl.HMACHex(fn, key, data) }

// HMACBytes returns HMAC digest bytes using the given hash function.
func HMACBytes(fn func() hash.Hash, key, data []byte) []byte {
	return cryptoimpl.HMACBytes(fn, key, data)
}

// HMACMD5Hex returns HMAC-MD5 in lower-case hex form.
func HMACMD5Hex(key, data []byte) string { return cryptoimpl.HMACMD5Hex(key, data) }

// HMACSHA1Hex returns HMAC-SHA1 in lower-case hex form.
func HMACSHA1Hex(key, data []byte) string { return cryptoimpl.HMACSHA1Hex(key, data) }

// HMACSHA256Hex returns HMAC-SHA256 in lower-case hex form.
func HMACSHA256Hex(key, data []byte) string { return cryptoimpl.HMACSHA256Hex(key, data) }

// HMACSHA512Hex returns HMAC-SHA512 in lower-case hex form.
func HMACSHA512Hex(key, data []byte) string { return cryptoimpl.HMACSHA512Hex(key, data) }

// HMACSHA384Hex returns HMAC-SHA384 in lower-case hex form.
func HMACSHA384Hex(key, data []byte) string { return cryptoimpl.HMACSHA384Hex(key, data) }

// HMACEqual compares two MAC values in constant time.
func HMACEqual(a, b []byte) bool { return cryptoimpl.HMACEqual(a, b) }

// ConstantTimeEqual compares two byte slices in constant time when lengths match.
func ConstantTimeEqual(a, b []byte) bool { return cryptoimpl.ConstantTimeEqual(a, b) }

// PBKDF2 derives a key from password and salt using PBKDF2.
func PBKDF2(password, salt []byte, iterations, keyLen int, fn func() hash.Hash) ([]byte, error) {
	return cryptoimpl.PBKDF2(password, salt, iterations, keyLen, fn)
}

// PBKDF2SHA1 derives a key using PBKDF2-HMAC-SHA1.
func PBKDF2SHA1(password, salt []byte, iterations, keyLen int) ([]byte, error) {
	return cryptoimpl.PBKDF2SHA1(password, salt, iterations, keyLen)
}

// PBKDF2SHA256 derives a key using PBKDF2-HMAC-SHA256.
func PBKDF2SHA256(password, salt []byte, iterations, keyLen int) ([]byte, error) {
	return cryptoimpl.PBKDF2SHA256(password, salt, iterations, keyLen)
}

// SignParams joins params by sorted key and returns the digest hex using digestHex.
func SignParams(params map[string]any, digestHex func([]byte) string, separator, keyValueSeparator string, ignoreNil bool, otherParams ...string) string {
	return cryptoimpl.SignParams(params, digestHex, separator, keyValueSeparator, ignoreNil, otherParams...)
}

// SignParamsMD5 signs sorted params with MD5.
func SignParamsMD5(params map[string]any, otherParams ...string) string {
	return cryptoimpl.SignParamsMD5(params, otherParams...)
}

// SignParamsSHA1 signs sorted params with SHA1.
func SignParamsSHA1(params map[string]any, otherParams ...string) string {
	return cryptoimpl.SignParamsSHA1(params, otherParams...)
}

// SignParamsSHA256 signs sorted params with SHA256.
func SignParamsSHA256(params map[string]any, otherParams ...string) string {
	return cryptoimpl.SignParamsSHA256(params, otherParams...)
}

// RandomBytes returns n cryptographically secure random bytes.
func RandomBytes(n int) ([]byte, error) { return cryptoimpl.RandomBytes(n) }

// GenerateAESKey returns a random AES key. Valid sizes are 16, 24, or 32 bytes.
func GenerateAESKey(size int) ([]byte, error) { return cryptoimpl.GenerateAESKey(size) }

// AESEncryptCBC encrypts plain data using AES-CBC with PKCS#7 padding.
func AESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCBC(plain, key, iv)
}

// AESDecryptCBC decrypts AES-CBC data using PKCS#7 padding.
func AESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCBC(cipherText, key, iv)
}

// AESEncryptECB encrypts plain data using AES-ECB with PKCS#7 padding.
func AESEncryptECB(plain, key []byte) ([]byte, error) { return cryptoimpl.AESEncryptECB(plain, key) }

// AESDecryptECB decrypts AES-ECB data using PKCS#7 padding.
func AESDecryptECB(cipherText, key []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptECB(cipherText, key)
}

// AESEncryptCTR encrypts or decrypts data using AES-CTR.
func AESEncryptCTR(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCTR(data, key, iv)
}

// AESDecryptCTR decrypts or encrypts data using AES-CTR.
func AESDecryptCTR(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCTR(data, key, iv)
}

// AESEncryptCFB encrypts data using AES-CFB.
func AESEncryptCFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptCFB(data, key, iv)
}

// AESDecryptCFB decrypts data using AES-CFB.
func AESDecryptCFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptCFB(data, key, iv)
}

// AESEncryptOFB encrypts or decrypts data using AES-OFB.
func AESEncryptOFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptOFB(data, key, iv)
}

// AESDecryptOFB decrypts or encrypts data using AES-OFB.
func AESDecryptOFB(data, key, iv []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptOFB(data, key, iv)
}

// AESEncryptGCM encrypts plain data using AES-GCM.
func AESEncryptGCM(plain, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESEncryptGCM(plain, key, nonce, additionalData)
}

// AESDecryptGCM decrypts AES-GCM data.
func AESDecryptGCM(cipherText, key, nonce, additionalData []byte) ([]byte, error) {
	return cryptoimpl.AESDecryptGCM(cipherText, key, nonce, additionalData)
}

// DESEncryptCBC encrypts plain data using DES-CBC with PKCS#7 padding.
func DESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	return cryptoimpl.DESEncryptCBC(plain, key, iv)
}

// DESDecryptCBC decrypts DES-CBC data using PKCS#7 padding.
func DESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	return cryptoimpl.DESDecryptCBC(cipherText, key, iv)
}

// TripleDESEncryptCBC encrypts plain data using 3DES-CBC with PKCS#7 padding.
func TripleDESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	return cryptoimpl.TripleDESEncryptCBC(plain, key, iv)
}

// TripleDESDecryptCBC decrypts 3DES-CBC data using PKCS#7 padding.
func TripleDESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	return cryptoimpl.TripleDESDecryptCBC(cipherText, key, iv)
}

// RC4Crypt encrypts or decrypts data using RC4.
func RC4Crypt(data, key []byte) ([]byte, error) { return cryptoimpl.RC4Crypt(data, key) }

// VigenereEncrypt encrypts printable ASCII text using the Vigenere helper algorithm.
func VigenereEncrypt(data, cipherKey string) (string, error) {
	return cryptoimpl.VigenereEncrypt(data, cipherKey)
}

// VigenereDecrypt decrypts text encrypted by VigenereEncrypt.
func VigenereDecrypt(data, cipherKey string) (string, error) {
	return cryptoimpl.VigenereDecrypt(data, cipherKey)
}

// XXTEAEncrypt encrypts data using XXTEA.
func XXTEAEncrypt(data, key []byte) []byte { return cryptoimpl.XXTEAEncrypt(data, key) }

// XXTEADecrypt decrypts data using XXTEA.
func XXTEADecrypt(data, key []byte) ([]byte, error) { return cryptoimpl.XXTEADecrypt(data, key) }

// GenerateRSAKey generates an RSA private key.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) { return cryptoimpl.GenerateRSAKey(bits) }

// RSAEncryptOAEP encrypts data using RSA-OAEP with SHA-256.
func RSAEncryptOAEP(plain []byte, pub *rsa.PublicKey, label []byte) ([]byte, error) {
	return cryptoimpl.RSAEncryptOAEP(plain, pub, label)
}

// RSADecryptOAEP decrypts data using RSA-OAEP with SHA-256.
func RSADecryptOAEP(cipherText []byte, priv *rsa.PrivateKey, label []byte) ([]byte, error) {
	return cryptoimpl.RSADecryptOAEP(cipherText, priv, label)
}

// RSAEncryptPKCS1v15 encrypts data using RSA PKCS#1 v1.5 padding.
func RSAEncryptPKCS1v15(plain []byte, pub *rsa.PublicKey) ([]byte, error) {
	return cryptoimpl.RSAEncryptPKCS1v15(plain, pub)
}

// RSADecryptPKCS1v15 decrypts data using RSA PKCS#1 v1.5 padding.
func RSADecryptPKCS1v15(cipherText []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return cryptoimpl.RSADecryptPKCS1v15(cipherText, priv)
}

// RSASignPKCS1v15 signs digest using RSA PKCS#1 v1.5.
func RSASignPKCS1v15(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte) ([]byte, error) {
	return cryptoimpl.RSASignPKCS1v15(priv, hash, digest)
}

// RSAVerifyPKCS1v15 verifies an RSA PKCS#1 v1.5 signature.
func RSAVerifyPKCS1v15(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte) error {
	return cryptoimpl.RSAVerifyPKCS1v15(pub, hash, digest, sig)
}

// RSASignPSS signs digest using RSA-PSS.
func RSASignPSS(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte) ([]byte, error) {
	return cryptoimpl.RSASignPSS(priv, hash, digest)
}

// RSAVerifyPSS verifies an RSA-PSS signature.
func RSAVerifyPSS(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte) error {
	return cryptoimpl.RSAVerifyPSS(pub, hash, digest, sig)
}

// SignSHA256WithRSA signs data using SHA256withRSA.
func SignSHA256WithRSA(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return cryptoimpl.SignSHA256WithRSA(data, priv)
}

// VerifySHA256WithRSA verifies SHA256withRSA signature.
func VerifySHA256WithRSA(data, sig []byte, pub *rsa.PublicKey) error {
	return cryptoimpl.VerifySHA256WithRSA(data, sig, pub)
}

// PrivateKeyToPEM encodes an RSA private key as PKCS#1 PEM.
func PrivateKeyToPEM(priv *rsa.PrivateKey) []byte { return cryptoimpl.PrivateKeyToPEM(priv) }

// PublicKeyToPEM encodes an RSA public key as PKIX PEM.
func PublicKeyToPEM(pub *rsa.PublicKey) ([]byte, error) { return cryptoimpl.PublicKeyToPEM(pub) }

// ParseRSAPrivateKeyPEM parses a PKCS#1 or PKCS#8 RSA private key PEM.
func ParseRSAPrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	return cryptoimpl.ParseRSAPrivateKeyPEM(data)
}

// ParseRSAPublicKeyPEM parses a PKIX or PKCS#1 RSA public key PEM.
func ParseRSAPublicKeyPEM(data []byte) (*rsa.PublicKey, error) {
	return cryptoimpl.ParseRSAPublicKeyPEM(data)
}

// PrivateKeyToPKCS8PEM encodes an RSA private key as PKCS#8 PEM.
func PrivateKeyToPKCS8PEM(priv *rsa.PrivateKey) ([]byte, error) {
	return cryptoimpl.PrivateKeyToPKCS8PEM(priv)
}

// PublicKeyToPKCS1PEM encodes an RSA public key as PKCS#1 PEM.
func PublicKeyToPKCS1PEM(pub *rsa.PublicKey) []byte { return cryptoimpl.PublicKeyToPKCS1PEM(pub) }

// ParseX509CertificatePEM parses an X.509 certificate from PEM data.
func ParseX509CertificatePEM(data []byte) (*x509.Certificate, error) {
	return cryptoimpl.ParseX509CertificatePEM(data)
}

// PublicKeyFromCertificatePEM extracts an RSA public key from an X.509 certificate PEM.
func PublicKeyFromCertificatePEM(data []byte) (*rsa.PublicKey, error) {
	return cryptoimpl.PublicKeyFromCertificatePEM(data)
}

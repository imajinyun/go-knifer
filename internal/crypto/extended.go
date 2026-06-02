package crypto

import (
	"bytes"
	stdcrypto "crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash"
	"sort"
)

const xxteaDelta uint32 = 0x9E3779B9

// MD5 returns the MD5 digest bytes of data.
func MD5(data []byte) []byte {
	sum := md5.Sum(data)
	return sum[:]
}

// SHA1 returns the SHA1 digest bytes of data.
func SHA1(data []byte) []byte {
	sum := sha1.Sum(data)
	return sum[:]
}

// SHA224 returns the SHA224 digest bytes of data.
func SHA224(data []byte) []byte {
	sum := sha256.Sum224(data)
	return sum[:]
}

// SHA256 returns the SHA256 digest bytes of data.
func SHA256(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// SHA384 returns the SHA384 digest bytes of data.
func SHA384(data []byte) []byte {
	sum := sha512.Sum384(data)
	return sum[:]
}

// SHA512 returns the SHA512 digest bytes of data.
func SHA512(data []byte) []byte {
	sum := sha512.Sum512(data)
	return sum[:]
}

// SHA384Hex returns the SHA384 digest of data in lower-case hex form.
func SHA384Hex(data []byte) string { return hex.EncodeToString(SHA384(data)) }

// SHA224Hex returns the SHA224 digest of data in lower-case hex form.
func SHA224Hex(data []byte) string { return hex.EncodeToString(SHA224(data)) }

// MD5Hex16 returns the middle 16 characters of the MD5 hex digest.
func MD5Hex16(data []byte) string { return MD5HexTo16(MD5Hex(data)) }

// MD5HexTo16 returns the middle 16 characters of a 32-character MD5 hex digest.
func MD5HexTo16(md5Hex string) string {
	if len(md5Hex) < 24 {
		return ""
	}
	return md5Hex[8:24]
}

// HMACBytes returns HMAC digest bytes using the given hash function.
func HMACBytes(fn func() hash.Hash, key, data []byte) []byte {
	h := hmac.New(fn, key)
	_, _ = h.Write(data)
	return h.Sum(nil)
}

// HMACSHA384Hex returns HMAC-SHA384 in lower-case hex form.
func HMACSHA384Hex(key, data []byte) string { return HMACHex(sha512.New384, key, data) }

// HMACEqual compares two MAC values in constant time.
func HMACEqual(a, b []byte) bool { return hmac.Equal(a, b) }

// ConstantTimeEqual compares two byte slices in constant time when lengths match.
func ConstantTimeEqual(a, b []byte) bool {
	return len(a) == len(b) && subtle.ConstantTimeCompare(a, b) == 1
}

// PBKDF2 derives a key from password and salt using PBKDF2.
func PBKDF2(password, salt []byte, iterations, keyLen int, fn func() hash.Hash) ([]byte, error) {
	if iterations <= 0 || keyLen <= 0 || fn == nil {
		return nil, ErrInvalidKey
	}
	h := fn()
	hLen := h.Size()
	nBlocks := (keyLen + hLen - 1) / hLen
	derived := make([]byte, 0, nBlocks*hLen)
	var blockIndex [4]byte
	for block := 1; block <= nBlocks; block++ {
		blockIndex[0] = byte(block >> 24)
		blockIndex[1] = byte(block >> 16)
		blockIndex[2] = byte(block >> 8)
		blockIndex[3] = byte(block)

		u := hmac.New(fn, password)
		_, _ = u.Write(salt)
		_, _ = u.Write(blockIndex[:])
		sum := u.Sum(nil)
		t := append([]byte(nil), sum...)

		for i := 1; i < iterations; i++ {
			u = hmac.New(fn, password)
			_, _ = u.Write(sum)
			sum = u.Sum(nil)
			for j := range t {
				t[j] ^= sum[j]
			}
		}
		derived = append(derived, t...)
	}
	return derived[:keyLen], nil
}

// PBKDF2SHA1 derives a key using PBKDF2-HMAC-SHA1.
func PBKDF2SHA1(password, salt []byte, iterations, keyLen int) ([]byte, error) {
	return PBKDF2(password, salt, iterations, keyLen, sha1.New)
}

// PBKDF2SHA256 derives a key using PBKDF2-HMAC-SHA256.
func PBKDF2SHA256(password, salt []byte, iterations, keyLen int) ([]byte, error) {
	return PBKDF2(password, salt, iterations, keyLen, sha256.New)
}

// SignParams joins params by sorted key and returns the digest hex using digestHex.
func SignParams(params map[string]any, digestHex func([]byte) string, separator, keyValueSeparator string, ignoreNil bool, otherParams ...string) string {
	keys := make([]string, 0, len(params))
	for key, value := range params {
		if ignoreNil && value == nil {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys)+len(otherParams))
	for _, key := range keys {
		value := params[key]
		parts = append(parts, key+keyValueSeparator+fmt.Sprint(value))
	}
	parts = append(parts, otherParams...)
	return digestHex([]byte(stringsJoin(parts, separator)))
}

// SignParamsMD5 signs sorted params with MD5.
func SignParamsMD5(params map[string]any, otherParams ...string) string {
	return SignParams(params, MD5Hex, "", "", true, otherParams...)
}

// SignParamsSHA1 signs sorted params with SHA1.
func SignParamsSHA1(params map[string]any, otherParams ...string) string {
	return SignParams(params, SHA1Hex, "", "", true, otherParams...)
}

// SignParamsSHA256 signs sorted params with SHA256.
func SignParamsSHA256(params map[string]any, otherParams ...string) string {
	return SignParams(params, SHA256Hex, "", "", true, otherParams...)
}

// AESEncryptECB encrypts plain data using AES-ECB with PKCS#7 padding.
func AESEncryptECB(plain, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return encryptECB(block, plain), nil
}

// AESDecryptECB decrypts AES-ECB data using PKCS#7 padding.
func AESDecryptECB(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return decryptECB(block, cipherText)
}

// AESEncryptCTR encrypts or decrypts data using AES-CTR.
func AESEncryptCTR(data, key, iv []byte) ([]byte, error) {
	return aesStream(data, key, iv, cipher.NewCTR)
}

// AESDecryptCTR decrypts or encrypts data using AES-CTR.
func AESDecryptCTR(data, key, iv []byte) ([]byte, error) { return AESEncryptCTR(data, key, iv) }

// AESEncryptCFB encrypts data using AES-CFB.
func AESEncryptCFB(data, key, iv []byte) ([]byte, error) {
	return aesCFB(data, key, iv, false)
}

// AESDecryptCFB decrypts data using AES-CFB.
func AESDecryptCFB(data, key, iv []byte) ([]byte, error) {
	return aesCFB(data, key, iv, true)
}

// AESEncryptOFB encrypts or decrypts data using AES-OFB.
func AESEncryptOFB(data, key, iv []byte) ([]byte, error) {
	return aesOFB(data, key, iv)
}

// AESDecryptOFB decrypts or encrypts data using AES-OFB.
func AESDecryptOFB(data, key, iv []byte) ([]byte, error) { return AESEncryptOFB(data, key, iv) }

// DESEncryptCBC encrypts plain data using DES-CBC with PKCS#7 padding.
func DESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return encryptCBC(block, plain, iv)
}

// DESDecryptCBC decrypts DES-CBC data using PKCS#7 padding.
func DESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return decryptCBC(block, cipherText, iv)
}

// TripleDESEncryptCBC encrypts plain data using 3DES-CBC with PKCS#7 padding.
func TripleDESEncryptCBC(plain, key, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return encryptCBC(block, plain, iv)
}

// TripleDESDecryptCBC decrypts 3DES-CBC data using PKCS#7 padding.
func TripleDESDecryptCBC(cipherText, key, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	return decryptCBC(block, cipherText, iv)
}

// RC4Crypt encrypts or decrypts data using RC4.
func RC4Crypt(data, key []byte) ([]byte, error) {
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	c.XORKeyStream(out, data)
	return out, nil
}

// VigenereEncrypt encrypts printable ASCII text using the Vigenere helper algorithm.
func VigenereEncrypt(data, cipherKey string) (string, error) {
	if cipherKey == "" {
		return "", ErrInvalidKey
	}
	dataRunes, keyRunes := []rune(data), []rune(cipherKey)
	out := make([]rune, len(dataRunes))
	for i, r := range dataRunes {
		k := keyRunes[i%len(keyRunes)]
		out[i] = (r+k-64)%95 + 32
	}
	return string(out), nil
}

// VigenereDecrypt decrypts text encrypted by VigenereEncrypt.
func VigenereDecrypt(data, cipherKey string) (string, error) {
	if cipherKey == "" {
		return "", ErrInvalidKey
	}
	dataRunes, keyRunes := []rune(data), []rune(cipherKey)
	out := make([]rune, len(dataRunes))
	for i, r := range dataRunes {
		k := keyRunes[i%len(keyRunes)]
		diff := r - k
		if diff >= 0 {
			out[i] = diff%95 + 32
		} else {
			out[i] = (diff+95)%95 + 32
		}
	}
	return string(out), nil
}

// XXTEAEncrypt encrypts data using XXTEA.
func XXTEAEncrypt(data, key []byte) []byte {
	if len(data) == 0 {
		return append([]byte(nil), data...)
	}
	return xxteaToByteArray(xxteaEncryptWords(xxteaToIntArray(data, true), xxteaToIntArray(xxteaFixKey(key), false)), false)
}

// XXTEADecrypt decrypts data using XXTEA.
func XXTEADecrypt(data, key []byte) ([]byte, error) {
	if len(data) == 0 {
		return append([]byte(nil), data...), nil
	}
	out := xxteaToByteArray(xxteaDecryptWords(xxteaToIntArray(data, false), xxteaToIntArray(xxteaFixKey(key), false)), true)
	if out == nil {
		return nil, ErrInvalidCipherText
	}
	return out, nil
}

// RSAEncryptPKCS1v15 encrypts data using RSA PKCS#1 v1.5 padding.
func RSAEncryptPKCS1v15(plain []byte, pub *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, plain)
}

// RSADecryptPKCS1v15 decrypts data using RSA PKCS#1 v1.5 padding.
func RSADecryptPKCS1v15(cipherText []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}

// RSASignPKCS1v15 signs digest using RSA PKCS#1 v1.5.
func RSASignPKCS1v15(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, priv, hash, digest)
}

// RSAVerifyPKCS1v15 verifies an RSA PKCS#1 v1.5 signature.
func RSAVerifyPKCS1v15(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte) error {
	return rsa.VerifyPKCS1v15(pub, hash, digest, sig)
}

// RSASignPSS signs digest using RSA-PSS.
func RSASignPSS(priv *rsa.PrivateKey, hash stdcrypto.Hash, digest []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, priv, hash, digest, nil)
}

// RSAVerifyPSS verifies an RSA-PSS signature.
func RSAVerifyPSS(pub *rsa.PublicKey, hash stdcrypto.Hash, digest, sig []byte) error {
	return rsa.VerifyPSS(pub, hash, digest, sig, nil)
}

// SignSHA256WithRSA signs data using SHA256withRSA.
func SignSHA256WithRSA(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
	digest := sha256.Sum256(data)
	return RSASignPKCS1v15(priv, stdcrypto.SHA256, digest[:])
}

// VerifySHA256WithRSA verifies SHA256withRSA signature.
func VerifySHA256WithRSA(data, sig []byte, pub *rsa.PublicKey) error {
	digest := sha256.Sum256(data)
	return RSAVerifyPKCS1v15(pub, stdcrypto.SHA256, digest[:], sig)
}

// PrivateKeyToPKCS8PEM encodes an RSA private key as PKCS#8 PEM.
func PrivateKeyToPKCS8PEM(priv *rsa.PrivateKey) ([]byte, error) {
	b, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b}), nil
}

// PublicKeyToPKCS1PEM encodes an RSA public key as PKCS#1 PEM.
func PublicKeyToPKCS1PEM(pub *rsa.PublicKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pub)})
}

// ParseX509CertificatePEM parses an X.509 certificate from PEM data.
func ParseX509CertificatePEM(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidKey
	}
	return x509.ParseCertificate(block.Bytes)
}

// PublicKeyFromCertificatePEM extracts an RSA public key from an X.509 certificate PEM.
func PublicKeyFromCertificatePEM(data []byte) (*rsa.PublicKey, error) {
	cert, err := ParseX509CertificatePEM(data)
	if err != nil {
		return nil, err
	}
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	return pub, nil
}

func aesBlockWithIV(key, iv []byte) (cipher.Block, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIV
	}
	return block, nil
}

func aesStream(data, key, iv []byte, newStream func(cipher.Block, []byte) cipher.Stream) ([]byte, error) {
	block, err := aesBlockWithIV(key, iv)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(data))
	newStream(block, iv).XORKeyStream(out, data)
	return out, nil
}

func aesCFB(data, key, iv []byte, decrypt bool) ([]byte, error) {
	block, err := aesBlockWithIV(key, iv)
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

func aesOFB(data, key, iv []byte) ([]byte, error) {
	block, err := aesBlockWithIV(key, iv)
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

func stringsJoin(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}
	if separator == "" {
		var b bytes.Buffer
		for _, part := range parts {
			b.WriteString(part)
		}
		return b.String()
	}
	var b bytes.Buffer
	for i, part := range parts {
		if i > 0 {
			b.WriteString(separator)
		}
		b.WriteString(part)
	}
	return b.String()
}

func xxteaEncryptWords(v, k []uint32) []uint32 {
	n := len(v) - 1
	if n < 1 {
		return v
	}
	z := v[n]
	q := 6 + 52/uint32(n+1)
	var sum uint32
	for q > 0 {
		q--
		sum += xxteaDelta
		e := (sum >> 2) & 3
		for p := 0; p < n; p++ {
			y := v[p+1]
			v[p] += xxteaMX(sum, y, z, uint32(p), e, k)
			z = v[p]
		}
		y := v[0]
		v[n] += xxteaMX(sum, y, z, uint32(n), e, k)
		z = v[n]
	}
	return v
}

func xxteaDecryptWords(v, k []uint32) []uint32 {
	n := len(v) - 1
	if n < 1 {
		return v
	}
	y := v[0]
	q := 6 + 52/uint32(n+1)
	sum := q * xxteaDelta
	for sum != 0 {
		e := (sum >> 2) & 3
		for p := n; p > 0; p-- {
			z := v[p-1]
			v[p] -= xxteaMX(sum, y, z, uint32(p), e, k)
			y = v[p]
		}
		z := v[n]
		v[0] -= xxteaMX(sum, y, z, 0, e, k)
		y = v[0]
		sum -= xxteaDelta
	}
	return v
}

func xxteaMX(sum, y, z, p, e uint32, k []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (k[p&3^e] ^ z))
}

func xxteaFixKey(key []byte) []byte {
	fixed := make([]byte, 16)
	copy(fixed, key)
	return fixed
}

func xxteaToIntArray(data []byte, includeLength bool) []uint32 {
	n := (len(data) + 3) >> 2
	if includeLength {
		n++
	}
	result := make([]uint32, n)
	for i, b := range data {
		result[i>>2] |= uint32(b) << ((i & 3) << 3)
	}
	if includeLength {
		result[n-1] = uint32(len(data))
	}
	return result
}

func xxteaToByteArray(data []uint32, includeLength bool) []byte {
	n := len(data) << 2
	if includeLength {
		m := int(data[len(data)-1])
		n -= 4
		if m < n-3 || m > n {
			return nil
		}
		n = m
	}
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		var word [4]byte
		binary.LittleEndian.PutUint32(word[:], data[i>>2])
		result[i] = word[i&3]
	}
	return result
}

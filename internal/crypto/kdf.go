package crypto

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"hash"
)

// PBKDF2 derives a key from password and salt using PBKDF2.
func PBKDF2(password, salt []byte, iterations, keyLen int, fn func() hash.Hash) ([]byte, error) {
	if iterations <= 0 || keyLen <= 0 || fn == nil {
		return nil, ErrInvalidKey
	}
	key, err := pbkdf2.Key(fn, string(password), salt, iterations, keyLen)
	if err != nil {
		return nil, ErrInvalidKey
	}
	return key, nil
}

// PBKDF2SHA256 derives a key using PBKDF2-HMAC-SHA256.
func PBKDF2SHA256(password, salt []byte, iterations, keyLen int) ([]byte, error) {
	return PBKDF2(password, salt, iterations, keyLen, sha256.New)
}

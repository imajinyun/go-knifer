package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/fnv"
)

// AdditiveHash calculates an additive hash modulo prime. Non-positive prime falls back to 31.
func AdditiveHash(s string, prime int) int {
	if prime <= 0 {
		prime = 31
	}
	h := len(s)
	for _, r := range s {
		h += int(r)
	}
	return h % prime
}

// FnvHash calculates a 32-bit FNV-1 hash.
func FnvHash(s string) uint32 {
	h := fnv.New32()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// MD5Hex calculates the MD5 digest and returns lowercase hex text.
func MD5Hex(s string) string {
	h := md5.Sum([]byte(s)) // #nosec G401 -- compatibility helper for non-security hashing.
	return hex.EncodeToString(h[:])
}

// SHA1Hex calculates the SHA-1 digest and returns lowercase hex text.
func SHA1Hex(s string) string {
	h := sha1.Sum([]byte(s)) // #nosec G401 -- compatibility helper for non-security hashing.
	return hex.EncodeToString(h[:])
}

// SHA256Hex calculates the SHA-256 digest and returns lowercase hex text.
func SHA256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

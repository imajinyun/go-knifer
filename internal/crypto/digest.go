package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// Digest returns digest bytes computed by newHash.
func Digest(data []byte, newHash func() hash.Hash) []byte {
	if newHash == nil {
		newHash = sha256.New
	}
	h := newHash()
	_, _ = h.Write(data)
	return h.Sum(nil)
}

// DigestHex returns the digest computed by newHash in lower-case hex form.
func DigestHex(data []byte, newHash func() hash.Hash) string {
	return hex.EncodeToString(Digest(data, newHash))
}

// SHA224 returns the SHA224 digest bytes of data.
func SHA224(data []byte) []byte {
	sum := sha256.Sum224(data)
	return sum[:]
}

// SHA224Hex returns the SHA224 digest of data in lower-case hex form.
func SHA224Hex(data []byte) string { return hex.EncodeToString(SHA224(data)) }

// SHA256 returns the SHA256 digest bytes of data.
func SHA256(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// SHA256Hex returns the SHA256 digest of data in lower-case hex form.
func SHA256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// SHA384 returns the SHA384 digest bytes of data.
func SHA384(data []byte) []byte {
	sum := sha512.Sum384(data)
	return sum[:]
}

// SHA384Hex returns the SHA384 digest of data in lower-case hex form.
func SHA384Hex(data []byte) string { return hex.EncodeToString(SHA384(data)) }

// SHA512 returns the SHA512 digest bytes of data.
func SHA512(data []byte) []byte {
	sum := sha512.Sum512(data)
	return sum[:]
}

// SHA512Hex returns the SHA512 digest of data in lower-case hex form.
func SHA512Hex(data []byte) string {
	sum := sha512.Sum512(data)
	return hex.EncodeToString(sum[:])
}

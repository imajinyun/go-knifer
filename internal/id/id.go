package id

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mathrand "math/rand"
	"time"
)

var defaultRand = mathrand.New(mathrand.NewSource(time.Now().UnixNano())) // #nosec G404 -- fallback only for IDs when crypto/rand is unavailable.

// SimpleUUID returns a 32-character UUID without hyphens.
func SimpleUUID() string {
	b := make([]byte, 16)
	if _, err := cryptorand.Read(b); err != nil {
		// Fallback path when crypto/rand is unavailable.
		binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UnixNano()))
		binary.BigEndian.PutUint64(b[8:], defaultRand.Uint64())
	}
	// version 4 / variant
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b)
}

// FastUUID returns a standard 8-4-4-4-12 UUID string.
func FastUUID() string {
	s := SimpleUUID()
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

// ObjectId returns a MongoDB-style 12-byte object id encoded as 24 hex characters.
// Layout: 4-byte Unix timestamp in seconds, 5 random bytes, and a 3-byte counter.
var objectIdCounter uint32

func ObjectId() string {
	now := uint32(time.Now().Unix()) // #nosec G115 -- ObjectId timestamp is intentionally stored in 4 bytes.
	rnd := make([]byte, 5)
	fillRandomBytes(rnd)
	c := nextCounter()
	b := make([]byte, 12)
	binary.BigEndian.PutUint32(b[0:4], now)
	copy(b[4:9], rnd)
	b[9] = byte(c >> 16)
	b[10] = byte(c >> 8)
	b[11] = byte(c)
	return hex.EncodeToString(b)
}

func nextCounter() uint32 {
	objectIdCounter++
	return objectIdCounter & 0x00ffffff
}

// NanoId returns a default 21-character NanoId using a URL-safe alphabet.
func NanoId() string { return NanoIdN(21) }

// NanoIdN returns a NanoId with the specified length.
func NanoIdN(n int) string {
	const alphabet = "_-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n <= 0 {
		return ""
	}
	mask := 63 // alphabet length is 64.
	step := (n*8 + 7) / 8
	out := make([]byte, 0, n)
	buf := make([]byte, step)
	for {
		fillRandomBytes(buf)
		for i := 0; i < step && len(out) < n; i++ {
			out = append(out, alphabet[buf[i]&byte(mask)])
		}
		if len(out) >= n {
			break
		}
	}
	return string(out[:n])
}

func fillRandomBytes(buf []byte) {
	if _, err := cryptorand.Read(buf); err != nil {
		for i := range buf {
			buf[i] = byte(defaultRand.Intn(256)) // #nosec G115 -- Intn(256) always fits in byte.
		}
	}
}

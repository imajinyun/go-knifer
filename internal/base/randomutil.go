package base

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mathrand "math/rand"
	"time"
)

// This file provides random-value helpers aligned with hutool-core RandomUtil.

// Character set constants used by random string helpers.
const (
	BaseNumber       = "0123456789"
	BaseChar         = "abcdefghijklmnopqrstuvwxyz"
	BaseCharNumber   = BaseChar + BaseNumber
	BaseCharNumberUC = BaseChar + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + BaseNumber
)

var defaultRand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

// RandomInt returns a random integer in [0, max). Non-positive max returns 0.
func RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return defaultRand.Intn(max)
}

// RandomIntRange returns a random integer in [min, max). If max <= min, min is returned.
func RandomIntRange(min, max int) int {
	if max <= min {
		return min
	}
	return min + defaultRand.Intn(max-min)
}

// RandomLong returns a non-negative random int64.
func RandomLong() int64 { return defaultRand.Int63() }

// RandomFloat returns a random float64 in [0.0, 1.0).
func RandomFloat() float64 { return defaultRand.Float64() }

// RandomBool returns a random boolean.
func RandomBool() bool { return defaultRand.Intn(2) == 0 }

// RandomBytes returns n cryptographically secure random bytes when possible.
func RandomBytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	buf := make([]byte, n)
	fillRandomBytes(buf)
	return buf
}

func fillRandomBytes(buf []byte) {
	if _, err := cryptorand.Read(buf); err != nil {
		// Fall back to math/rand when crypto/rand is unavailable.
		for i := range buf {
			buf[i] = byte(defaultRand.Intn(256))
		}
	}
}

// RandomString returns a random string from BaseCharNumber, using lowercase letters and digits.
func RandomString(n int) string { return RandomStringFrom(BaseCharNumber, n) }

// RandomNumbers returns a random numeric string.
func RandomNumbers(n int) string { return RandomStringFrom(BaseNumber, n) }

// RandomStringUpper returns a random string with lowercase letters, uppercase letters, and digits.
func RandomStringUpper(n int) string { return RandomStringFrom(BaseCharNumberUC, n) }

// RandomStringFrom builds a random string by sampling runes from charset.
func RandomStringFrom(charset string, n int) string {
	if n <= 0 || len(charset) == 0 {
		return ""
	}
	rs := []rune(charset)
	out := make([]rune, n)
	for i := 0; i < n; i++ {
		out[i] = rs[defaultRand.Intn(len(rs))]
	}
	return string(out)
}

// RandomEle returns a random element, or the zero value for an empty slice.
func RandomEle[T any](a []T) T {
	if len(a) == 0 {
		var zero T
		return zero
	}
	return a[defaultRand.Intn(len(a))]
}

// This section provides ID helpers aligned with hutool-core IdUtil.

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
	now := uint32(time.Now().Unix())
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

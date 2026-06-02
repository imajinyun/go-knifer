package hash

import (
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

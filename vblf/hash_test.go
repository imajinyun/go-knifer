package vblf_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vblf"
)

func TestFacadeHashFunctions(t *testing.T) {
	// smoke test: hash functions should return consistent values
	h1 := vblf.BloomRSHash("abc")
	h2 := vblf.BloomRSHash("abc")
	if h1 != h2 {
		t.Fatal("hash function should be deterministic")
	}
}

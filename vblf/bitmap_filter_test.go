package vblf_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vblf"
)

func TestFacadeBitMapBloomFilterWithOptions(t *testing.T) {
	bitmap := vblf.NewBitMapBloomFilterWithOptions(
		vblf.WithBitMapSize(5),
		vblf.WithBloomFilters(vblf.NewFNVFilter(1<<20), vblf.NewRSFilter(1<<20)),
	)
	if !bitmap.Add("world") || !bitmap.Contains("world") {
		t.Fatal("expected options-created bitmap filter to contain value")
	}
}

func TestFacadeBitMap(t *testing.T) {
	bm := vblf.NewIntMap(100)
	bm.Add(42)
	if !bm.Contains(42) {
		t.Fatal("expected bitmap to contain 42")
	}
}

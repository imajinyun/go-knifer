package bloomfilter

import "testing"

func TestBitMapBloomFilter(t *testing.T) {
	bf := NewBitMapBloomFilter(5)
	if !bf.Add("foo") {
		t.Fatal("add foo should return true")
	}
	if !bf.Contains("foo") {
		t.Fatal("should contain foo")
	}
	if bf.Add("foo") {
		t.Fatal("repeat add foo should return false")
	}
	if bf.Contains("not-in-filter-12345") {
		t.Fatal("should not contain unknown token")
	}
}

func TestBitMapBloomFilter_CustomFilters(t *testing.T) {
	bf := NewBitMapBloomFilterWithFilters(
		5,
		NewFNVFilter(1<<20),
		NewRSFilter(1<<20),
	)
	if !bf.Add("bar") {
		t.Fatal()
	}
	if !bf.Contains("bar") {
		t.Fatal()
	}
}

func TestBitMapBloomFilterWithOptions(t *testing.T) {
	bf := NewBitMapBloomFilterWithOptions(WithBitMapSize(5))
	if len(bf.filters) != 5 {
		t.Fatalf("default filter count = %d, want 5", len(bf.filters))
	}
	if !bf.Add("foo") || !bf.Contains("foo") {
		t.Fatal("options-created bitmap filter should add and contain value")
	}

	custom := NewBitMapBloomFilterWithOptions(
		WithBitMapSize(5),
		WithBloomFilters(NewFNVFilter(1<<20), NewRSFilter(1<<20)),
	)
	if len(custom.filters) != 2 {
		t.Fatalf("custom filter count = %d, want 2", len(custom.filters))
	}
	if !custom.Add("bar") || !custom.Contains("bar") {
		t.Fatal("custom options-created bitmap filter should add and contain value")
	}
}

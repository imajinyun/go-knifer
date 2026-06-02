package bloomfilter

// CreateBitSet creates a BitSet-based Bloom filter.
// See NewBitSetBloomFilter for detailed parameter semantics.
func CreateBitSet(c, n, k int) *BitSetBloomFilter { return NewBitSetBloomFilter(c, n, k) }

// CreateBitMap creates a BitMap-backed Bloom filter.
func CreateBitMap(m int) *BitMapBloomFilter { return NewBitMapBloomFilter(m) }

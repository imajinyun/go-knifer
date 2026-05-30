package bloomfilter

// CreateBitSet creates a BitSet-based Bloom filter corresponding to hutool BloomFilterUtil.createBitSet.
// See NewBitSetBloomFilter for detailed parameter semantics.
func CreateBitSet(c, n, k int) *BitSetBloomFilter { return NewBitSetBloomFilter(c, n, k) }

// CreateBitMap creates a BitMap-backed Bloom filter corresponding to hutool BloomFilterUtil.createBitMap.
func CreateBitMap(m int) *BitMapBloomFilter { return NewBitMapBloomFilter(m) }

package bloomfilter

// CreateBitSet 创建一个基于 BitSet 的布隆过滤器，对应 hutool BloomFilterUtil.createBitSet。
// 详细参数语义见 NewBitSetBloomFilter。
func CreateBitSet(c, n, k int) *BitSetBloomFilter { return NewBitSetBloomFilter(c, n, k) }

// CreateBitMap 创建 BitMap 实现的布隆过滤器，对应 hutool BloomFilterUtil.createBitMap。
func CreateBitMap(m int) *BitMapBloomFilter { return NewBitMapBloomFilter(m) }

package bloomfilter

// BitMapBloomFilter 基于多 Filter 组合的布隆过滤器，对应 hutool BitMapBloomFilter。
// 内部聚合多个 BloomFilter（默认使用 5 个不同哈希过滤器）。
type BitMapBloomFilter struct {
	filters []BloomFilter
}

// NewBitMapBloomFilter 使用默认的 5 个过滤器（DefaultFilter/ELFFilter/JSFilter/PJWFilter/SDBMFilter）。
//
// m: M 值（MB），决定底层 BitMap 的大小。最终位数 = m/5 * 1024 * 1024 * 8。
func NewBitMapBloomFilter(m int) *BitMapBloomFilter {
	mNum := int64(m) / 5
	size := mNum * 1024 * 1024 * 8
	return &BitMapBloomFilter{
		filters: []BloomFilter{
			NewDefaultFilter(size),
			NewELFFilter(size),
			NewJSFilter(size),
			NewPJWFilter(size),
			NewSDBMFilter(size),
		},
	}
}

// NewBitMapBloomFilterWithFilters 使用自定义的多个过滤器构造 BitMapBloomFilter。
// 与 hutool 行为一致，仍按 m 校验，但使用传入的 filters 替换默认过滤器集合。
func NewBitMapBloomFilterWithFilters(m int, filters ...BloomFilter) *BitMapBloomFilter {
	b := NewBitMapBloomFilter(m)
	if len(filters) > 0 {
		b.filters = filters
	}
	return b
}

// Add 实现 BloomFilter.Add。任一过滤器新增成功即视为新增。
func (b *BitMapBloomFilter) Add(str string) bool {
	flag := false
	for _, f := range b.filters {
		if f.Add(str) {
			flag = true
		}
	}
	return flag
}

// Contains 实现 BloomFilter.Contains。所有过滤器都判定包含才视为包含。
func (b *BitMapBloomFilter) Contains(str string) bool {
	for _, f := range b.filters {
		if !f.Contains(str) {
			return false
		}
	}
	return true
}

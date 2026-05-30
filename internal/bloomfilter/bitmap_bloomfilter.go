package bloomfilter

// BitMapBloomFilter is a Bloom filter composed from multiple filters, mirroring hutool BitMapBloomFilter.
// It aggregates several BloomFilter instances and uses five different hash filters by default.
type BitMapBloomFilter struct {
	filters []BloomFilter
}

// NewBitMapBloomFilter uses five default filters: Default, ELF, JS, PJW, and SDBM.
//
// m is the M value in MB and controls the underlying BitMap size. Final bits = m/5 * 1024 * 1024 * 8.
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

// NewBitMapBloomFilterWithFilters creates a BitMapBloomFilter with custom filters.
// It keeps hutool-compatible m validation while replacing the default filter set.
func NewBitMapBloomFilterWithFilters(m int, filters ...BloomFilter) *BitMapBloomFilter {
	b := NewBitMapBloomFilter(m)
	if len(filters) > 0 {
		b.filters = filters
	}
	return b
}

// Add implements BloomFilter.Add. The value is considered added if any filter changes.
func (b *BitMapBloomFilter) Add(str string) bool {
	flag := false
	for _, f := range b.filters {
		if f.Add(str) {
			flag = true
		}
	}
	return flag
}

// Contains implements BloomFilter.Contains. All filters must report containment.
func (b *BitMapBloomFilter) Contains(str string) bool {
	for _, f := range b.filters {
		if !f.Contains(str) {
			return false
		}
	}
	return true
}

package bloomfilter

// 机器位常量。
const (
	Machine32 = 32
	Machine64 = 64
)

// BitMap 用于将某个 int 或 long 值映射到一个数组中，从而判定该值是否存在。
type BitMap interface {
	// Add 加入值。
	Add(i int64)
	// Contains 判断是否包含值。
	Contains(i int64) bool
	// Remove 移除值。
	Remove(i int64)
}

// IntMap 32 位机器友好的 BitMap 实现。
type IntMap struct {
	ints []int32
}

// NewIntMap 构造指定容量的 IntMap，单位为 int32 槽数（每槽容纳 32 位）。
func NewIntMap(size int) *IntMap { return &IntMap{ints: make([]int32, size)} }

// Add 实现 BitMap.Add。
func (m *IntMap) Add(i int64) {
	r := int(i / Machine32)
	c := uint32(i & (Machine32 - 1))
	m.ints[r] |= 1 << c
}

// Contains 实现 BitMap.Contains。
func (m *IntMap) Contains(i int64) bool {
	r := int(i / Machine32)
	c := uint32(i & (Machine32 - 1))
	return ((uint32(m.ints[r]) >> c) & 1) == 1
}

// Remove 实现 BitMap.Remove。
func (m *IntMap) Remove(i int64) {
	r := int(i / Machine32)
	c := uint32(i & (Machine32 - 1))
	m.ints[r] &= ^(1 << c)
}

// LongMap 64 位机器友好的 BitMap 实现。
type LongMap struct {
	longs []int64
}

// NewLongMap 构造指定容量的 LongMap。
func NewLongMap(size int) *LongMap { return &LongMap{longs: make([]int64, size)} }

// Add 实现 BitMap.Add。
func (m *LongMap) Add(i int64) {
	r := int(i / Machine64)
	c := uint64(i & (Machine64 - 1))
	m.longs[r] |= 1 << c
}

// Contains 实现 BitMap.Contains。
func (m *LongMap) Contains(i int64) bool {
	r := int(i / Machine64)
	c := uint64(i & (Machine64 - 1))
	return ((uint64(m.longs[r]) >> c) & 1) == 1
}

// Remove 实现 BitMap.Remove。
func (m *LongMap) Remove(i int64) {
	r := int(i / Machine64)
	c := uint64(i & (Machine64 - 1))
	m.longs[r] &= ^(1 << c)
}

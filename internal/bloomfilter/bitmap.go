package bloomfilter

// Machine word size constants.
const (
	Machine32 = 32
	Machine64 = 64
)

// BitMap maps integer values into an array to test whether a value exists.
type BitMap interface {
	// Add inserts a value.
	Add(i int64)
	// Contains reports whether the value exists.
	Contains(i int64) bool
	// Remove deletes a value.
	Remove(i int64)
}

// IntMap is a 32-bit word based BitMap implementation.
type IntMap struct {
	ints []int32
}

// NewIntMap creates an IntMap with size int32 slots, each storing 32 bits.
func NewIntMap(size int) *IntMap { return &IntMap{ints: make([]int32, size)} }

// Add implements BitMap.Add.
func (m *IntMap) Add(i int64) {
	r := int(i / Machine32)
	c := uint32(i & (Machine32 - 1))
	m.ints[r] |= 1 << c
}

// Contains implements BitMap.Contains.
func (m *IntMap) Contains(i int64) bool {
	r := int(i / Machine32)
	c := uint32(i & (Machine32 - 1))
	return ((uint32(m.ints[r]) >> c) & 1) == 1
}

// Remove implements BitMap.Remove.
func (m *IntMap) Remove(i int64) {
	r := int(i / Machine32)
	c := uint32(i & (Machine32 - 1))
	m.ints[r] &= ^(1 << c)
}

// LongMap is a 64-bit word based BitMap implementation.
type LongMap struct {
	longs []int64
}

// NewLongMap creates a LongMap with size int64 slots.
func NewLongMap(size int) *LongMap { return &LongMap{longs: make([]int64, size)} }

// Add implements BitMap.Add.
func (m *LongMap) Add(i int64) {
	r := int(i / Machine64)
	c := uint64(i & (Machine64 - 1))
	m.longs[r] |= 1 << c
}

// Contains implements BitMap.Contains.
func (m *LongMap) Contains(i int64) bool {
	r := int(i / Machine64)
	c := uint64(i & (Machine64 - 1))
	return ((uint64(m.longs[r]) >> c) & 1) == 1
}

// Remove implements BitMap.Remove.
func (m *LongMap) Remove(i int64) {
	r := int(i / Machine64)
	c := uint64(i & (Machine64 - 1))
	m.longs[r] &= ^(1 << c)
}

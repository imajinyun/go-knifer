package hash

// This file implements a collection of well-known non-cryptographic string
// hash algorithms aligned with the Hutool HashUtil set. To keep behavior
// aligned with Java, all operations use int32 truncation semantics equivalent
// to a Java int.

// charsOf splits a string from a Java char perspective: ASCII maps directly to
// bytes, while non-ASCII characters map approximately to Unicode code points.
// This matches Java charAt for BMP characters.
func charsOf(s string) []rune { return []rune(s) }

// RsHash implements the RS algorithm.
func RsHash(str string) int32 {
	var b int32 = 378551
	var a int32 = 63689
	var hash int32 = 0
	for _, c := range charsOf(str) {
		hash = hash*a + c
		a *= b
	}
	return hash & 0x7FFFFFFF
}

// JsHash implements the JS algorithm.
func JsHash(str string) int32 {
	var hash int32 = 1315423911
	for _, c := range charsOf(str) {
		hash ^= (hash << 5) + c + (hash >> 2)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash & 0x7FFFFFFF
}

// PjwHash implements the PJW algorithm.
func PjwHash(str string) int32 {
	const bitsInUnsignedInt = 32
	const threeQuarters = (bitsInUnsignedInt * 3) / 4
	const oneEighth = bitsInUnsignedInt / 8
	shift := uint32(bitsInUnsignedInt - oneEighth)
	highBits := int32(uint32(0xFFFFFFFF) << shift) // #nosec G115 -- intentional bit reinterpretation to match Java int semantics.
	var hash int32
	var test int32
	for _, c := range charsOf(str) {
		hash = (hash << oneEighth) + c
		if test = hash & highBits; test != 0 {
			hash = (hash ^ (test >> threeQuarters)) & (^highBits)
		}
	}
	return hash & 0x7FFFFFFF
}

// ElfHash implements the ELF algorithm.
func ElfHash(str string) int32 {
	var hash int32
	var x int32
	var maskU uint32 = 0xF0000000
	mask := int32(maskU) // #nosec G115 -- intentional bit reinterpretation to match Java int semantics.
	for _, c := range charsOf(str) {
		hash = (hash << 4) + c
		if x = hash & mask; x != 0 {
			hash ^= x >> 24
			hash &= ^x
		}
	}
	return hash & 0x7FFFFFFF
}

// BkdrHash implements the BKDR algorithm.
func BkdrHash(str string) int32 {
	const seed int32 = 131
	var hash int32
	for _, c := range charsOf(str) {
		hash = hash*seed + c
	}
	return hash & 0x7FFFFFFF
}

// SdbmHash implements the SDBM algorithm.
func SdbmHash(str string) int32 {
	var hash int32
	for _, c := range charsOf(str) {
		hash = c + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}

// DjbHash implements the DJB algorithm.
func DjbHash(str string) int32 {
	var hash int32 = 5381
	for _, c := range charsOf(str) {
		hash = ((hash << 5) + hash) + c
	}
	return hash & 0x7FFFFFFF
}

// ApHash implements the AP algorithm.
func ApHash(str string) int32 {
	var hash int32
	for i, c := range charsOf(str) {
		if i&1 == 0 {
			hash ^= (hash << 7) ^ c ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ c ^ (hash >> 5))
		}
	}
	return hash
}

// FnvHashString implements the improved 32-bit FNV-1 algorithm for strings.
//
// This differs from FnvHash, which uses the standard library hash/fnv FNV-1.
func FnvHashString(data string) int32 {
	const p int32 = 16777619
	var seed uint32 = 2166136261
	hash := int32(seed) // #nosec G115 -- intentional bit reinterpretation to match Java int semantics.
	for _, c := range charsOf(data) {
		hash = (hash ^ c) * p
	}
	hash += hash << 13
	hash ^= int32(uint32(hash) >> 7) // #nosec G115 -- logical right shift to match Java >>> semantics.
	hash += hash << 3
	hash ^= int32(uint32(hash) >> 17) // #nosec G115 -- logical right shift to match Java >>> semantics.
	hash += hash << 5
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// HfHash implements the HF hash algorithm.
func HfHash(data string) int64 {
	var hash int64
	for i, c := range charsOf(data) {
		hash += int64(c) * 3 * int64(i)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// HfIpHash implements the HFIP hash algorithm.
func HfIpHash(data string) int64 {
	chars := charsOf(data)
	length := len(chars)
	var hash int64
	for i := 0; i < length; i++ {
		hash += int64(chars[i%4] ^ chars[i])
	}
	return hash
}

// TianlHash implements the TianL hash algorithm.
func TianlHash(str string) int64 {
	chars := charsOf(str)
	iLength := len(chars)
	if iLength == 0 {
		return 0
	}
	var hash int64
	if iLength <= 256 {
		hash = 16777216 * int64(iLength-1)
	} else {
		hash = 4278190080
	}
	process := func(i int, ch rune) {
		c := ch
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		hash += (3*int64(i)*int64(c)*int64(c) + 5*int64(i)*int64(c) + 7*int64(i) + 11*int64(c)) % 16777216
	}
	if iLength <= 96 {
		for i := 1; i <= iLength; i++ {
			process(i, chars[i-1])
		}
	} else {
		for i := 1; i <= 96; i++ {
			process(i, chars[i+iLength-96-1])
		}
	}
	if hash < 0 {
		hash *= -1
	}
	return hash
}

// JavaDefaultHash simulates Java String.hashCode.
func JavaDefaultHash(str string) int32 {
	var h int32
	for _, c := range charsOf(str) {
		h = 31*h + c
	}
	return h
}

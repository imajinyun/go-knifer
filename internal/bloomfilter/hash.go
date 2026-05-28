package bloomfilter

// 本文件实现 hutool HashUtil 中布隆过滤器所需的哈希函数子集。
// 为保持与 Java 端实现行为一致，所有运算均使用与 Java int(32 位有符号) 等价的 int32 截断语义。

// charsOf 将字符串按 Java char 视角拆解：ASCII 直接取字节，非 ASCII 取 UTF-16 code unit。
// 由于 hutool 在常见 ASCII 字符串场景下与 Java 行为一致，这里同样按 rune（即 Unicode 码点）展开，
// 对 BMP 内字符与 Java charAt 等价；对辅助平面字符仅做近似处理（与 hutool 行为足以等价用于布隆过滤器）。
func charsOf(s string) []rune { return []rune(s) }

// RsHash RS 算法。
func RsHash(str string) int32 {
	var b int32 = 378551
	var a int32 = 63689
	var hash int32 = 0
	for _, c := range charsOf(str) {
		hash = hash*a + int32(c)
		a = a * b
	}
	return hash & 0x7FFFFFFF
}

// JsHash JS 算法。
func JsHash(str string) int32 {
	var hash int32 = 1315423911
	for _, c := range charsOf(str) {
		hash ^= (hash << 5) + int32(c) + (hash >> 2)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash & 0x7FFFFFFF
}

// PjwHash PJW 算法。
func PjwHash(str string) int32 {
	const bitsInUnsignedInt = 32
	const threeQuarters = (bitsInUnsignedInt * 3) / 4
	const oneEighth = bitsInUnsignedInt / 8
	shift := uint32(bitsInUnsignedInt - oneEighth)
	highBits := int32(uint32(0xFFFFFFFF) << shift)
	var hash int32
	var test int32
	for _, c := range charsOf(str) {
		hash = (hash << oneEighth) + int32(c)
		if test = hash & highBits; test != 0 {
			hash = (hash ^ (test >> threeQuarters)) & (^highBits)
		}
	}
	return hash & 0x7FFFFFFF
}

// ElfHash ELF 算法。
func ElfHash(str string) int32 {
	var hash int32
	var x int32
	var maskU uint32 = 0xF0000000
	mask := int32(maskU)
	for _, c := range charsOf(str) {
		hash = (hash << 4) + int32(c)
		if x = hash & mask; x != 0 {
			hash ^= x >> 24
			hash &= ^x
		}
	}
	return hash & 0x7FFFFFFF
}

// BkdrHash BKDR 算法。
func BkdrHash(str string) int32 {
	const seed int32 = 131
	var hash int32
	for _, c := range charsOf(str) {
		hash = hash*seed + int32(c)
	}
	return hash & 0x7FFFFFFF
}

// SdbmHash SDBM 算法。
func SdbmHash(str string) int32 {
	var hash int32
	for _, c := range charsOf(str) {
		hash = int32(c) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}

// DjbHash DJB 算法。
func DjbHash(str string) int32 {
	var hash int32 = 5381
	for _, c := range charsOf(str) {
		hash = ((hash << 5) + hash) + int32(c)
	}
	return hash & 0x7FFFFFFF
}

// ApHash AP 算法。
func ApHash(str string) int32 {
	var hash int32
	for i, c := range charsOf(str) {
		if i&1 == 0 {
			hash ^= (hash << 7) ^ int32(c) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ int32(c) ^ (hash >> 5))
		}
	}
	return hash
}

// FnvHashString 改进的 32 位 FNV-1 算法（字符串）。
func FnvHashString(data string) int32 {
	const p int32 = 16777619
	var seed uint32 = 2166136261
	hash := int32(seed)
	for _, c := range charsOf(data) {
		hash = (hash ^ int32(c)) * p
	}
	hash += hash << 13
	hash ^= int32(uint32(hash) >> 7)
	hash += hash << 3
	hash ^= int32(uint32(hash) >> 17)
	hash += hash << 5
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// HfHash HF Hash 算法。
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

// HfIpHash HFIP Hash 算法。
func HfIpHash(data string) int64 {
	chars := charsOf(data)
	length := len(chars)
	var hash int64
	for i := 0; i < length; i++ {
		hash += int64(chars[i%4] ^ chars[i])
	}
	return hash
}

// TianlHash TianL Hash 算法。
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

// JavaDefaultHash 模拟 Java String.hashCode 的算法。
func JavaDefaultHash(str string) int32 {
	var h int32
	for _, c := range charsOf(str) {
		h = 31*h + int32(c)
	}
	return h
}

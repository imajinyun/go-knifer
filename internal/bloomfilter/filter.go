package bloomfilter

import "fmt"

// BloomFilter 布隆过滤器接口，对应 hutool BloomFilter。
type BloomFilter interface {
	// Contains 判断字符串是否在过滤器中存在（存在误差）。
	Contains(str string) bool
	// Add 添加字符串到过滤器，若已存在返回 false，否则添加并返回 true。
	Add(str string) bool
}

// HashFunc 用于 FuncFilter 的字符串哈希函数。
type HashFunc func(str string) int64

// FuncFilter 基于自定义哈希函数的布隆过滤器，对应 hutool FuncFilter / AbstractFilter。
type FuncFilter struct {
	bm       BitMap
	size     int64
	hashFunc HashFunc
}

// DefaultMachineNum FuncFilter 默认机器位（与 hutool 一致使用 32 位）。
var DefaultMachineNum = Machine32

// NewFuncFilter 使用默认机器位构造 FuncFilter。
func NewFuncFilter(maxValue int64, hashFunc HashFunc) *FuncFilter {
	return NewFuncFilterWithMachineNum(maxValue, DefaultMachineNum, hashFunc)
}

// NewFuncFilterWithMachineNum 使用指定机器位构造 FuncFilter。
func NewFuncFilterWithMachineNum(maxValue int64, machineNum int, hashFunc HashFunc) *FuncFilter {
	if maxValue < 1 || maxValue > 0x7FFFFFFF {
		panic(fmt.Sprintf("maxValue must be between 1 and %d", int64(0x7FFFFFFF)))
	}
	capacity := int((maxValue + int64(machineNum) - 1) / int64(machineNum))
	var bm BitMap
	switch machineNum {
	case Machine32:
		bm = NewIntMap(capacity)
	case Machine64:
		bm = NewLongMap(capacity)
	default:
		panic("Error Machine number!")
	}
	return &FuncFilter{bm: bm, size: maxValue, hashFunc: hashFunc}
}

// hash 调用底层哈希函数并对 size 取模、取绝对值。
func (f *FuncFilter) hash(str string) int64 {
	v := f.hashFunc(str) % f.size
	if v < 0 {
		v = -v
	}
	return v
}

// Contains 实现 BloomFilter.Contains。
func (f *FuncFilter) Contains(str string) bool { return f.bm.Contains(f.hash(str)) }

// Add 实现 BloomFilter.Add。
func (f *FuncFilter) Add(str string) bool {
	h := f.hash(str)
	if f.bm.Contains(h) {
		return false
	}
	f.bm.Add(h)
	return true
}

//============= 基于具体哈希算法的便捷过滤器 =============

// NewDefaultFilter 默认布隆过滤器（Java String.hashCode）。
func NewDefaultFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(JavaDefaultHash(s)) })
}

// NewELFFilter ELF 哈希过滤器。
func NewELFFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(ElfHash(s)) })
}

// NewFNVFilter FNV 哈希过滤器。
func NewFNVFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(FnvHashString(s)) })
}

// NewHfFilter HF 哈希过滤器。
func NewHfFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, HfHash)
}

// NewHfIpFilter HFIP 哈希过滤器。
func NewHfIpFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, HfIpHash)
}

// NewJSFilter JS 哈希过滤器。
func NewJSFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(JsHash(s)) })
}

// NewPJWFilter PJW 哈希过滤器。
func NewPJWFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(PjwHash(s)) })
}

// NewRSFilter RS 哈希过滤器。
func NewRSFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(RsHash(s)) })
}

// NewSDBMFilter SDBM 哈希过滤器。
func NewSDBMFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, func(s string) int64 { return int64(SdbmHash(s)) })
}

// NewTianlFilter TianL 哈希过滤器。
func NewTianlFilter(maxValue int64) *FuncFilter {
	return NewFuncFilter(maxValue, TianlHash)
}

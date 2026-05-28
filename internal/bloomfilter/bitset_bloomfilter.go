package bloomfilter

import (
	"bufio"
	"io"
	"math"
	"os"
)

// BitSetBloomFilter 使用大小固定的位集合实现的布隆过滤器，对应 hutool BitSetBloomFilter。
// 算法使用固定顺序，只需指定个数即可。
type BitSetBloomFilter struct {
	bits               []uint64 // 模拟 BitSet
	bitSetSize         int
	addedElements      int
	hashFunctionNumber int
}

// NewBitSetBloomFilter 构造布隆过滤器，过滤器容量为 c*k 个 bit。
//
// c: 当前过滤器预先开辟的最大包含记录数（通常要比预计存入的记录多一倍）。
// n: 当前过滤器预计所要包含的记录。
// k: 哈希函数的个数（取值 1~8）。
func NewBitSetBloomFilter(c, n, k int) *BitSetBloomFilter {
	if c <= 0 {
		panic("Parameter c must be positive")
	}
	if n <= 0 {
		panic("Parameter n must be positive")
	}
	if k < 1 || k > 8 {
		panic("hashFunctionNumber must be between 1 and 8")
	}
	size := c * k
	return &BitSetBloomFilter{
		bits:               make([]uint64, (size+63)/64),
		bitSetSize:         size,
		addedElements:      n,
		hashFunctionNumber: k,
	}
}

func (b *BitSetBloomFilter) setBit(pos int) { b.bits[pos>>6] |= 1 << uint(pos&63) }

func (b *BitSetBloomFilter) getBit(pos int) bool {
	return (b.bits[pos>>6]>>uint(pos&63))&1 == 1
}

// InitFromFile 通过文件初始化过滤器，逐行 add。
func (b *BitSetBloomFilter) InitFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			// 去除尾部换行
			for len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
				line = line[:len(line)-1]
			}
			b.Add(line)
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// Add 加入字符串，若已存在返回 false。
func (b *BitSetBloomFilter) Add(str string) bool {
	if b.Contains(str) {
		return false
	}
	positions := b.createHashes(str, b.hashFunctionNumber)
	for _, v := range positions {
		pos := absInt(v % int32(b.bitSetSize))
		b.setBit(int(pos))
	}
	return true
}

// Contains 判断是否包含。
func (b *BitSetBloomFilter) Contains(str string) bool {
	positions := b.createHashes(str, b.hashFunctionNumber)
	for _, v := range positions {
		pos := absInt(v % int32(b.bitSetSize))
		if !b.getBit(int(pos)) {
			return false
		}
	}
	return true
}

// FalsePositiveProbability 当前过滤器误判率：(1 - e^(-k * n / m)) ^ k。
func (b *BitSetBloomFilter) FalsePositiveProbability() float64 {
	return math.Pow(1-math.Exp(-float64(b.hashFunctionNumber)*float64(b.addedElements)/float64(b.bitSetSize)),
		float64(b.hashFunctionNumber))
}

// createHashes 多哈希。
func (b *BitSetBloomFilter) createHashes(str string, hashNumber int) []int32 {
	out := make([]int32, hashNumber)
	for i := 0; i < hashNumber; i++ {
		out[i] = bitSetHash(str, i)
	}
	return out
}

// bitSetHash 与 hutool BitSetBloomFilter.hash 一致。
func bitSetHash(str string, k int) int32 {
	switch k {
	case 0:
		return RsHash(str)
	case 1:
		return JsHash(str)
	case 2:
		return ElfHash(str)
	case 3:
		return BkdrHash(str)
	case 4:
		return ApHash(str)
	case 5:
		return DjbHash(str)
	case 6:
		return SdbmHash(str)
	case 7:
		return PjwHash(str)
	default:
		return 0
	}
}

func absInt(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

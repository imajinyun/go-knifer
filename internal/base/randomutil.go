package base

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mathrand "math/rand"
	"time"
)

// 对应 hutool-core RandomUtil。

// 字符集合常量。
const (
	BaseNumber       = "0123456789"
	BaseChar         = "abcdefghijklmnopqrstuvwxyz"
	BaseCharNumber   = BaseChar + BaseNumber
	BaseCharNumberUC = BaseChar + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + BaseNumber
)

var defaultRand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

// RandomInt 返回 [0, max) 的随机整数；max<=0 返回 0。
func RandomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return defaultRand.Intn(max)
}

// RandomIntRange 返回 [min, max) 的随机整数。
func RandomIntRange(min, max int) int {
	if max <= min {
		return min
	}
	return min + defaultRand.Intn(max-min)
}

// RandomLong 返回非负随机 int64。
func RandomLong() int64 { return defaultRand.Int63() }

// RandomFloat 返回 [0.0, 1.0) 的随机浮点数。
func RandomFloat() float64 { return defaultRand.Float64() }

// RandomBool 返回随机布尔。
func RandomBool() bool { return defaultRand.Intn(2) == 0 }

// RandomBytes 返回指定长度的加密安全随机字节。
func RandomBytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	buf := make([]byte, n)
	if _, err := cryptorand.Read(buf); err != nil {
		// 退回到 math/rand
		for i := range buf {
			buf[i] = byte(defaultRand.Intn(256))
		}
	}
	return buf
}

// RandomString 从 BaseCharNumber 中随机产生指定长度字符串（小写字母+数字）。
func RandomString(n int) string { return RandomStringFrom(BaseCharNumber, n) }

// RandomNumbers 仅使用数字的随机字符串。
func RandomNumbers(n int) string { return RandomStringFrom(BaseNumber, n) }

// RandomStringUpper 大小写字母+数字。
func RandomStringUpper(n int) string { return RandomStringFrom(BaseCharNumberUC, n) }

// RandomStringFrom 在指定字符集合中随机抽取构建字符串。
func RandomStringFrom(charset string, n int) string {
	if n <= 0 || len(charset) == 0 {
		return ""
	}
	rs := []rune(charset)
	out := make([]rune, n)
	for i := 0; i < n; i++ {
		out[i] = rs[defaultRand.Intn(len(rs))]
	}
	return string(out)
}

// RandomEle 随机选取一个元素。
func RandomEle[T any](a []T) T {
	if len(a) == 0 {
		var zero T
		return zero
	}
	return a[defaultRand.Intn(len(a))]
}

// 对应 hutool-core IdUtil。

// SimpleUUID 32 位 UUID（无连字符）。
func SimpleUUID() string {
	b := make([]byte, 16)
	if _, err := cryptorand.Read(b); err != nil {
		// 退化方案
		binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UnixNano()))
		binary.BigEndian.PutUint64(b[8:], defaultRand.Uint64())
	}
	// version 4 / variant
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b)
}

// FastUUID 标准 8-4-4-4-12 UUID。
func FastUUID() string {
	s := SimpleUUID()
	return s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

// ObjectId 12 字节 MongoDB 风格的 ObjectId（24 hex 字符）：4 字节秒级时间戳 + 5 字节随机 + 3 字节计数。
var objectIdCounter uint32

func ObjectId() string {
	now := uint32(time.Now().Unix())
	rnd := make([]byte, 5)
	_, _ = cryptorand.Read(rnd)
	c := nextCounter()
	b := make([]byte, 12)
	binary.BigEndian.PutUint32(b[0:4], now)
	copy(b[4:9], rnd)
	b[9] = byte(c >> 16)
	b[10] = byte(c >> 8)
	b[11] = byte(c)
	return hex.EncodeToString(b)
}

func nextCounter() uint32 {
	objectIdCounter++
	return objectIdCounter & 0x00ffffff
}

// NanoId 生成默认 21 位 NanoId（URL 安全字符集）。
func NanoId() string { return NanoIdN(21) }

// NanoIdN 生成指定长度 NanoId。
func NanoIdN(n int) string {
	const alphabet = "_-0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n <= 0 {
		return ""
	}
	mask := 63 // alphabet 长度 64
	step := (n*8 + 7) / 8
	out := make([]byte, 0, n)
	buf := make([]byte, step)
	for {
		_, _ = cryptorand.Read(buf)
		for i := 0; i < step && len(out) < n; i++ {
			out = append(out, alphabet[buf[i]&byte(mask)])
		}
		if len(out) >= n {
			break
		}
	}
	return string(out[:n])
}

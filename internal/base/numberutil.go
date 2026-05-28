package base

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// 对应 hutool-core NumberUtil。

// NumberAdd 高精度加法（避免浮点误差）。
func NumberAdd(a, b float64) float64 {
	x := big.NewFloat(a)
	y := big.NewFloat(b)
	r, _ := new(big.Float).Add(x, y).Float64()
	return r
}

// NumberSub 高精度减法。
func NumberSub(a, b float64) float64 {
	x := big.NewFloat(a)
	y := big.NewFloat(b)
	r, _ := new(big.Float).Sub(x, y).Float64()
	return r
}

// NumberMul 高精度乘法。
func NumberMul(a, b float64) float64 {
	x := big.NewFloat(a)
	y := big.NewFloat(b)
	r, _ := new(big.Float).Mul(x, y).Float64()
	return r
}

// NumberDiv 除法保留 scale 位小数（HALF_UP 模式）。scale<0 不舍入。
func NumberDiv(a, b float64, scale int) float64 {
	if b == 0 {
		return 0
	}
	r := a / b
	if scale < 0 {
		return r
	}
	return Round(r, scale)
}

// Round HALF_UP 四舍五入到指定小数位。
func Round(v float64, scale int) float64 {
	if scale < 0 {
		scale = 0
	}
	pow := math.Pow(10, float64(scale))
	if v >= 0 {
		return math.Floor(v*pow+0.5) / pow
	}
	return -math.Floor(-v*pow+0.5) / pow
}

// IsNumber 是否为合法数字（支持整数、小数、负号）。
func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsInteger 是否为整数（仅含可选负号 + 数字）。
func IsInteger(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// IsDigits 是否全是数字字符（无符号）。
func IsDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Min 最小值。
func Min[T Ordered](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}
	m := nums[0]
	for _, v := range nums[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// Max 最大值。
func Max[T Ordered](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}
	m := nums[0]
	for _, v := range nums[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// Sum 求和。
func Sum[T Number](nums ...T) T {
	var s T
	for _, v := range nums {
		s += v
	}
	return s
}

// Avg 平均值（float64）。
func Avg[T Number](nums ...T) float64 {
	if len(nums) == 0 {
		return 0
	}
	var s float64
	for _, v := range nums {
		s += float64(v)
	}
	return s / float64(len(nums))
}

// Range 生成 [start, end) 步长 step 的序列；step 为 0 时视作 1（或 -1）。
func Range(start, end, step int) []int {
	if step == 0 {
		if end >= start {
			step = 1
		} else {
			step = -1
		}
	}
	out := make([]int, 0)
	if step > 0 {
		for i := start; i < end; i += step {
			out = append(out, i)
		}
	} else {
		for i := start; i > end; i += step {
			out = append(out, i)
		}
	}
	return out
}

// Equals 浮点数是否相等（精度 1e-9）。
func Equals(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// DecimalFormat 按 fmt 风格格式化（如 "%.2f"）。
func DecimalFormat(format string, v float64) string {
	if format == "" || !strings.Contains(format, "%") {
		// 类似 "0.00" 这种格式，统一转 fmt
		dot := strings.Index(format, ".")
		if dot >= 0 {
			n := len(format) - dot - 1
			return strconv.FormatFloat(v, 'f', n, 64)
		}
		return strconv.FormatFloat(v, 'f', 0, 64)
	}
	return ""
}

// 类型约束。

// Number 数值类型集合。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Ordered 有序类型集合。
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

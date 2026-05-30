package base

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

// This file provides numeric helpers aligned with hutool-core NumberUtil.

// NumberAdd performs high-precision addition to reduce floating-point rounding surprises.
func NumberAdd(a, b float64) float64 {
	x := big.NewFloat(a)
	y := big.NewFloat(b)
	r, _ := new(big.Float).Add(x, y).Float64()
	return r
}

// NumberSub performs high-precision subtraction.
func NumberSub(a, b float64) float64 {
	x := big.NewFloat(a)
	y := big.NewFloat(b)
	r, _ := new(big.Float).Sub(x, y).Float64()
	return r
}

// NumberMul performs high-precision multiplication.
func NumberMul(a, b float64) float64 {
	x := big.NewFloat(a)
	y := big.NewFloat(b)
	r, _ := new(big.Float).Mul(x, y).Float64()
	return r
}

// NumberDiv divides a by b and rounds to scale decimal places with HALF_UP semantics.
// A negative scale disables rounding; division by zero returns 0 for compatibility.
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

// Round rounds v to scale decimal places with HALF_UP semantics.
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

// IsNumber reports whether s is a valid number accepted by strconv.ParseFloat.
func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsInteger reports whether s is a valid base-10 integer.
func IsInteger(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// IsDigits reports whether s contains only unsigned ASCII digits.
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

// Min returns the minimum value, or the zero value when no values are provided.
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

// Max returns the maximum value, or the zero value when no values are provided.
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

// Sum returns the sum of all values.
func Sum[T Number](nums ...T) T {
	var s T
	for _, v := range nums {
		s += v
	}
	return s
}

// Avg returns the arithmetic mean as float64, or 0 for empty input.
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

// Range returns a half-open integer sequence [start, end) using step.
// A zero step is normalized to 1 or -1 based on the direction.
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

// Equals compares two floats using a fixed 1e-9 tolerance.
func Equals(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// DecimalFormat formats v with simple decimal patterns such as "0.00".
func DecimalFormat(format string, v float64) string {
	if format == "" || !strings.Contains(format, "%") {
		// Convert simple patterns like "0.00" to fixed decimal precision.
		dot := strings.Index(format, ".")
		if dot >= 0 {
			n := len(format) - dot - 1
			return strconv.FormatFloat(v, 'f', n, 64)
		}
		return strconv.FormatFloat(v, 'f', 0, 64)
	}
	return ""
}

// Generic numeric constraints.

// Number is the set of supported numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Ordered is the set of supported ordered types.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

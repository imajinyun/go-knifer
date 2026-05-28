package cron

// PartMatcher 字段匹配器接口，对应 hutool 的 PartMatcher。
type PartMatcher interface {
	// Match 判断给定值是否匹配。
	Match(v int) bool
	// NextAfter 返回不小于 v 的下一个匹配值，超过最大值时回绕到最小值。
	NextAfter(v int) int
}

// alwaysTrueMatcher 对应 AlwaysTrueMatcher，匹配一切。
type alwaysTrueMatcher struct{}

// AlwaysTrueMatcher 单例。
var AlwaysTrueMatcher PartMatcher = alwaysTrueMatcher{}

func (alwaysTrueMatcher) Match(int) bool      { return true }
func (alwaysTrueMatcher) NextAfter(v int) int { return v }

// boolArrayMatcher 对应 BoolArrayMatcher，使用布尔数组保存可匹配值。
type boolArrayMatcher struct {
	values   []bool
	minValue int
	maxValue int
}

// newBoolArrayMatcher 根据值集合创建 matcher。values 中的值会被拷贝。
func newBoolArrayMatcher(values []int) *boolArrayMatcher {
	if len(values) == 0 {
		return &boolArrayMatcher{}
	}
	maxV := values[0]
	for _, v := range values {
		if v > maxV {
			maxV = v
		}
	}
	arr := make([]bool, maxV+1)
	minV := -1
	for _, v := range values {
		if v < 0 || v > maxV {
			continue
		}
		arr[v] = true
		if minV == -1 || v < minV {
			minV = v
		}
	}
	return &boolArrayMatcher{values: arr, minValue: minV, maxValue: maxV}
}

func (m *boolArrayMatcher) Match(v int) bool {
	if v < 0 || v >= len(m.values) {
		return false
	}
	return m.values[v]
}

func (m *boolArrayMatcher) NextAfter(v int) int {
	if v > m.maxValue {
		return m.minValue
	}
	if v < 0 {
		v = 0
	}
	for i := v; i <= m.maxValue; i++ {
		if m.values[i] {
			return i
		}
	}
	return m.minValue
}

// MinValue 返回最小匹配值。
func (m *boolArrayMatcher) MinValue() int { return m.minValue }

// MaxValue 返回最大匹配值。
func (m *boolArrayMatcher) MaxValue() int { return m.maxValue }

// dayOfMonthMatcher 对应 DayOfMonthMatcher，支持 L（最后一天）。
type dayOfMonthMatcher struct {
	*boolArrayMatcher
}

const lastDayOfMonthSentinel = 32

func newDayOfMonthMatcher(values []int) *dayOfMonthMatcher {
	return &dayOfMonthMatcher{boolArrayMatcher: newBoolArrayMatcher(values)}
}

// IsLast 判断表达式是否包含 L。
func (m *dayOfMonthMatcher) IsLast() bool {
	return lastDayOfMonthSentinel < len(m.values) && m.values[lastDayOfMonthSentinel]
}

// MatchDay 判断 day 是否匹配（含 L 处理）。
func (m *dayOfMonthMatcher) MatchDay(day, month int, leap bool) bool {
	if m.boolArrayMatcher.Match(day) {
		return true
	}
	if m.IsLast() && day == lastDayOfMonth(month, leap) {
		return true
	}
	return false
}

// yearValueMatcher 对应 YearValueMatcher。
type yearValueMatcher struct {
	values []int // 已排序去重
}

func newYearValueMatcher(values []int) *yearValueMatcher {
	// 排序去重
	seen := make(map[int]struct{}, len(values))
	out := make([]int, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	// 简单插入排序
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return &yearValueMatcher{values: out}
}

func (m *yearValueMatcher) Match(v int) bool {
	for _, x := range m.values {
		if x == v {
			return true
		}
	}
	return false
}

// NextAfter 返回不小于 v 的最小年份，若全部小于 v 返回 -1。
func (m *yearValueMatcher) NextAfter(v int) int {
	for _, x := range m.values {
		if x >= v {
			return x
		}
	}
	return -1
}

// lastDayOfMonth 返回某月最后一天。
func lastDayOfMonth(month int, leap bool) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if leap {
			return 29
		}
		return 28
	}
	return 31
}

// isLeapYear 判断是否闰年。
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

package cron

// PartMatcher is the field matcher interface aligned with the utility toolkit PartMatcher.
type PartMatcher interface {
	// Match reports whether the given value matches.
	Match(v int) bool
	// NextAfter returns the next matched value not less than v, wrapping to the minimum after the maximum.
	NextAfter(v int) int
}

// alwaysTrueMatcher is aligned with AlwaysTrueMatcher and matches everything.
type alwaysTrueMatcher struct{}

// AlwaysTrueMatcher is the singleton matcher that matches everything.
var AlwaysTrueMatcher PartMatcher = alwaysTrueMatcher{}

func (alwaysTrueMatcher) Match(int) bool      { return true }
func (alwaysTrueMatcher) NextAfter(v int) int { return v }

// boolArrayMatcher is aligned with BoolArrayMatcher and stores matched values in a bool array.
type boolArrayMatcher struct {
	values   []bool
	minValue int
	maxValue int
}

// newBoolArrayMatcher creates a matcher from a value set; values are copied.
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

// MinValue returns the minimum matched value.
func (m *boolArrayMatcher) MinValue() int { return m.minValue }

// MaxValue returns the maximum matched value.
func (m *boolArrayMatcher) MaxValue() int { return m.maxValue }

// dayOfMonthMatcher is aligned with DayOfMonthMatcher and supports L for the last day.
type dayOfMonthMatcher struct {
	*boolArrayMatcher
}

const lastDayOfMonthSentinel = 32

func newDayOfMonthMatcher(values []int) *dayOfMonthMatcher {
	return &dayOfMonthMatcher{boolArrayMatcher: newBoolArrayMatcher(values)}
}

// IsLast reports whether the expression contains L.
func (m *dayOfMonthMatcher) IsLast() bool {
	return lastDayOfMonthSentinel < len(m.values) && m.values[lastDayOfMonthSentinel]
}

// MatchDay reports whether day matches, including L handling.
func (m *dayOfMonthMatcher) MatchDay(day, month int, leap bool) bool {
	if m.Match(day) {
		return true
	}
	if m.IsLast() && day == lastDayOfMonth(month, leap) {
		return true
	}
	return false
}

// yearValueMatcher is aligned with YearValueMatcher.
type yearValueMatcher struct {
	values []int // sorted and deduplicated
}

func newYearValueMatcher(values []int) *yearValueMatcher {
	// Deduplicate values.
	seen := make(map[int]struct{}, len(values))
	out := make([]int, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	// Simple insertion sort.
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

// NextAfter returns the smallest year not less than v, or -1 when all years are smaller.
func (m *yearValueMatcher) NextAfter(v int) int {
	for _, x := range m.values {
		if x >= v {
			return x
		}
	}
	return -1
}

// lastDayOfMonth returns the last day of a month.
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

// isLeapYear reports whether year is a leap year.
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

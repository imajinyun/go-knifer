package cron

import (
	"strconv"
	"strings"
)

// monthAlias 月份别名（1-based）。
var monthAlias = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// weekAlias 星期别名（0=sun ~ 6=sat）。
var weekAlias = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
}

// parsePart 解析单字段表达式，返回对应的 PartMatcher。
func parsePart(part Part, value string) (PartMatcher, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, NewCronError("empty cron field")
	}
	values, err := parsePartValues(part, value)
	if err != nil {
		return nil, err
	}
	switch part {
	case PartDayOfMonth:
		return newDayOfMonthMatcher(values), nil
	case PartYear:
		return newYearValueMatcher(values), nil
	default:
		// 当字段为 *（且没有 step 或具体值），可使用 AlwaysTrueMatcher，但 BoolArrayMatcher 也能正常工作。
		// 这里使用 BoolArrayMatcher 以保持简单一致。
		return newBoolArrayMatcher(values), nil
	}
}

// parsePartValues 解析字段，返回所有可匹配的值。* 或 ? 视为 AlwaysTrue 但调用方需要先判断。
func parsePartValues(part Part, value string) ([]int, error) {
	// 列表
	tokens := strings.Split(value, ",")
	var result []int
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			return nil, NewCronError("empty list element in %q", value)
		}
		vals, err := parseRangeStep(part, tok)
		if err != nil {
			return nil, err
		}
		result = append(result, vals...)
	}
	return result, nil
}

// parseRangeStep 解析形如 a, a-b, *, a/n, *​/n, a-b/n。
func parseRangeStep(part Part, value string) ([]int, error) {
	step := 1
	rangeStr := value
	if i := strings.Index(value, "/"); i >= 0 {
		rangeStr = value[:i]
		stepStr := value[i+1:]
		s, err := strconv.Atoi(stepStr)
		if err != nil || s <= 0 {
			return nil, NewCronError("invalid step %q", stepStr)
		}
		step = s
	}
	return parseRange(part, rangeStr, step)
}

// parseRange 解析 a / a-b / * / ? 等。
func parseRange(part Part, value string, step int) ([]int, error) {
	if value == "*" || value == "?" {
		// 全部范围按 step 取值
		return enumerate(part.Min(), part.Max(), step), nil
	}

	// 单值（包含 L、月份/星期别名、数字、负数）
	if !strings.Contains(value, "-") || (strings.HasPrefix(value, "-") && strings.Count(value, "-") == 1) {
		// 注意 "-3" 表示从最大值往回数 3
		v, err := parseSingle(part, value)
		if err != nil {
			return nil, err
		}
		if step == 1 {
			return []int{v}, nil
		}
		// a/step：从 a 到 max
		return enumerate(v, part.Max(), step), nil
	}

	// 范围 a-b
	idx := strings.Index(value, "-")
	// 处理类似 "-3-5" 这种形态：第一个字符是 '-' 表示负数起点
	if value[0] == '-' {
		idx = strings.Index(value[1:], "-")
		if idx >= 0 {
			idx++
		}
	}
	if idx < 0 {
		v, err := parseSingle(part, value)
		if err != nil {
			return nil, err
		}
		return []int{v}, nil
	}
	beginStr := value[:idx]
	endStr := value[idx+1:]
	begin, err := parseSingle(part, beginStr)
	if err != nil {
		return nil, err
	}
	end, err := parseSingle(part, endStr)
	if err != nil {
		return nil, err
	}
	if begin <= end {
		return enumerate(begin, end, step), nil
	}
	// 反向区间 b-a：等价于 [b..max] ∪ [min..a]
	first := enumerate(begin, part.Max(), step)
	second := enumerate(part.Min(), end, step)
	return append(first, second...), nil
}

// parseSingle 解析单个数字 / 别名 / L / 负数。
func parseSingle(part Part, value string) (int, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return 0, NewCronError("empty value")
	}
	// 别名
	lower := strings.ToLower(v)
	if lower == "l" {
		return part.Max(), nil
	}
	if part == PartMonth {
		if n, ok := monthAlias[lower]; ok {
			return n, nil
		}
	}
	if part == PartDayOfWeek {
		if n, ok := weekAlias[lower]; ok {
			return n, nil
		}
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, NewCronError("invalid number %q", v)
	}
	// 星期 7 视为 0
	if part == PartDayOfWeek && n == 7 {
		n = 0
	}
	// 负数视为相对最大值的回绕
	if n < 0 {
		n += part.Max()
	}
	if err := part.CheckValue(n); err != nil {
		return 0, err
	}
	return n, nil
}

// enumerate 在 [begin, end] 内按 step 枚举。
func enumerate(begin, end, step int) []int {
	if step <= 0 {
		step = 1
	}
	if begin > end {
		return nil
	}
	out := make([]int, 0, (end-begin)/step+1)
	for i := begin; i <= end; i += step {
		out = append(out, i)
	}
	return out
}

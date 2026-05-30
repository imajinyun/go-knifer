package cron

import (
	"strconv"
	"strings"
)

// monthAlias maps month aliases to one-based month values.
var monthAlias = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
	"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// weekAlias maps weekday aliases from 0 for Sunday to 6 for Saturday.
var weekAlias = map[string]int{
	"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
}

// parsePart parses a single field expression and returns the corresponding PartMatcher.
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
		// For a plain * field without a step or concrete values, AlwaysTrueMatcher could be used.
		// BoolArrayMatcher also works correctly here and keeps the implementation simple and consistent.
		return newBoolArrayMatcher(values), nil
	}
}

// parsePartValues parses a field and returns all matchable values; callers handle plain * or ? first.
func parsePartValues(part Part, value string) ([]int, error) {
	// List.
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

// parseRangeStep parses forms such as a, a-b, *, a/n, */n, and a-b/n.
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

// parseRange parses forms such as a, a-b, *, and ?.
func parseRange(part Part, value string, step int) ([]int, error) {
	if value == "*" || value == "?" {
		// Enumerate the whole range by step.
		return enumerate(part.Min(), part.Max(), step), nil
	}

	// Single value, including L, month/weekday aliases, numbers, and negative numbers.
	if !strings.Contains(value, "-") || (strings.HasPrefix(value, "-") && strings.Count(value, "-") == 1) {
		// Note that "-3" means counting back 3 from the maximum value.
		v, err := parseSingle(part, value)
		if err != nil {
			return nil, err
		}
		if step == 1 {
			return []int{v}, nil
		}
		// a/step: enumerate from a to max.
		return enumerate(v, part.Max(), step), nil
	}

	// Range a-b.
	idx := strings.Index(value, "-")
	// Handle forms like "-3-5", where the first '-' indicates a negative start value.
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
	// Reversed range b-a is equivalent to [b..max] ∪ [min..a].
	first := enumerate(begin, part.Max(), step)
	second := enumerate(part.Min(), end, step)
	return append(first, second...), nil
}

// parseSingle parses a single number, alias, L, or negative value.
func parseSingle(part Part, value string) (int, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return 0, NewCronError("empty value")
	}
	// Aliases.
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
	// Weekday 7 is treated as 0.
	if part == PartDayOfWeek && n == 7 {
		n = 0
	}
	// Negative numbers wrap relative to the maximum value.
	if n < 0 {
		n += part.Max()
	}
	if err := part.CheckValue(n); err != nil {
		return 0, err
	}
	return n, nil
}

// enumerate enumerates values in [begin, end] by step.
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

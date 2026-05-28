package cron

import (
	"strings"
	"time"
)

// patternMatcher 对应 hutool 的 PatternMatcher，由 7 个字段 matcher 组成。
type patternMatcher struct {
	matchers [7]PartMatcher
}

// newPatternMatcher 解析单个（不含 |）cron 表达式为 patternMatcher。
func newPatternMatcher(expr string) (*patternMatcher, error) {
	parts := strings.Fields(expr)
	switch len(parts) {
	case 5:
		// minute hour dom month dow
		parts = append([]string{"0"}, parts...)
		parts = append(parts, "*")
	case 6:
		// second minute hour dom month dow
		parts = append(parts, "*")
	case 7:
		// 完整
	default:
		return nil, NewCronError("invalid cron expression %q (parts=%d)", expr, len(parts))
	}
	pm := &patternMatcher{}
	partOrder := []Part{PartSecond, PartMinute, PartHour, PartDayOfMonth, PartMonth, PartDayOfWeek, PartYear}
	for i, p := range partOrder {
		raw := parts[i]
		if raw == "*" || raw == "?" {
			pm.matchers[p] = AlwaysTrueMatcher
			continue
		}
		m, err := parsePart(p, raw)
		if err != nil {
			return nil, err
		}
		pm.matchers[p] = m
	}
	return pm, nil
}

// match 判断给定的 fields 是否匹配。
// fields 顺序：[second, minute, hour, dayOfMonth, month, dayOfWeek, year]。
// 当 second < 0 时跳过秒匹配。
func (pm *patternMatcher) match(fields [7]int) bool {
	if fields[PartSecond] >= 0 {
		if !pm.matchers[PartSecond].Match(fields[PartSecond]) {
			return false
		}
	}
	if !pm.matchers[PartMinute].Match(fields[PartMinute]) {
		return false
	}
	if !pm.matchers[PartHour].Match(fields[PartHour]) {
		return false
	}
	// 日：DayOfMonth 需结合 month 与闰年判断 L
	if dom, ok := pm.matchers[PartDayOfMonth].(*dayOfMonthMatcher); ok {
		if !dom.MatchDay(fields[PartDayOfMonth], fields[PartMonth], isLeapYear(fields[PartYear])) {
			return false
		}
	} else if !pm.matchers[PartDayOfMonth].Match(fields[PartDayOfMonth]) {
		return false
	}
	if !pm.matchers[PartMonth].Match(fields[PartMonth]) {
		return false
	}
	if !pm.matchers[PartDayOfWeek].Match(fields[PartDayOfWeek]) {
		return false
	}
	if !pm.matchers[PartYear].Match(fields[PartYear]) {
		return false
	}
	return true
}

// Pattern 对应 hutool 的 CronPattern：可以由多个用 | 分隔的子表达式组成，任一匹配即认为匹配。
type Pattern struct {
	raw      string
	matchers []*patternMatcher
}

// NewPattern 解析 cron 表达式（支持 | 分隔多组表达式）。
func NewPattern(expr string) (*Pattern, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, NewCronError("empty cron expression")
	}
	groups := strings.Split(expr, "|")
	matchers := make([]*patternMatcher, 0, len(groups))
	for _, g := range groups {
		g = strings.TrimSpace(g)
		if g == "" {
			return nil, NewCronError("empty sub-expression in %q", expr)
		}
		pm, err := newPatternMatcher(g)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, pm)
	}
	return &Pattern{raw: expr, matchers: matchers}, nil
}

// MustNewPattern 与 NewPattern 类似，解析失败时 panic。
func MustNewPattern(expr string) *Pattern {
	p, err := NewPattern(expr)
	if err != nil {
		panic(err)
	}
	return p
}

// Raw 返回原始表达式。
func (p *Pattern) Raw() string { return p.raw }

// Match 判断给定时间是否匹配。matchSecond 为 false 时忽略秒字段。
func (p *Pattern) Match(t time.Time, matchSecond bool) bool {
	fields := timeToFields(t, matchSecond)
	for _, pm := range p.matchers {
		if pm.match(fields) {
			return true
		}
	}
	return false
}

// timeToFields 将 time.Time 转为 [second, minute, hour, dayOfMonth, month, dayOfWeek, year]。
// 周字段：周日 = 0 ~ 周六 = 6；月份从 1 开始；如果不匹配秒，second 设置为 -1。
func timeToFields(t time.Time, matchSecond bool) [7]int {
	sec := t.Second()
	if !matchSecond {
		sec = -1
	}
	dow := int(t.Weekday()) // time.Sunday == 0
	return [7]int{
		sec,
		t.Minute(),
		t.Hour(),
		t.Day(),
		int(t.Month()),
		dow,
		t.Year(),
	}
}

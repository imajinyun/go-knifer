package cron

import (
	"strconv"
	"strings"
	"time"
)

type patternConfig struct {
	parseInt func(string) (int, error)
}

// PatternOption customizes cron pattern parsing per call.
type PatternOption func(*patternConfig)

// WithPatternIntParser sets the integer parser used by NewPatternWithOptions.
func WithPatternIntParser(parser func(string) (int, error)) PatternOption {
	return func(c *patternConfig) {
		if parser != nil {
			c.parseInt = parser
		}
	}
}

func applyPatternOptions(opts []PatternOption) patternConfig {
	cfg := patternConfig{parseInt: strconv.Atoi}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.parseInt == nil {
		cfg.parseInt = strconv.Atoi
	}
	return cfg
}

// patternMatcher is aligned with the utility toolkit PatternMatcher and consists of seven field matchers.
type patternMatcher struct {
	matchers [7]PartMatcher
}

// newPatternMatcher parses one cron expression without | into a patternMatcher.
func newPatternMatcher(expr string, cfg patternConfig) (*patternMatcher, error) {
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
		// Full expression.
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
		m, err := parsePart(p, raw, cfg)
		if err != nil {
			return nil, err
		}
		pm.matchers[p] = m
	}
	return pm, nil
}

// match reports whether the given fields match.
// fields order: [second, minute, hour, dayOfMonth, month, dayOfWeek, year].
// When second < 0, second matching is skipped.
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
	// Day-of-month needs month and leap-year context to evaluate L.
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

// Pattern is aligned with the utility toolkit CronPattern and may contain multiple | separated sub-expressions.
type Pattern struct {
	raw      string
	matchers []*patternMatcher
}

// NewPattern parses a cron expression and supports multiple | separated groups.
func NewPattern(expr string) (*Pattern, error) {
	return NewPatternWithOptions(expr)
}

// NewPatternWithOptions parses a cron expression with custom parser providers and supports multiple | separated groups.
func NewPatternWithOptions(expr string, opts ...PatternOption) (*Pattern, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, NewCronError("empty cron expression")
	}
	cfg := applyPatternOptions(opts)
	groups := strings.Split(expr, "|")
	matchers := make([]*patternMatcher, 0, len(groups))
	for _, g := range groups {
		g = strings.TrimSpace(g)
		if g == "" {
			return nil, NewCronError("empty sub-expression in %q", expr)
		}
		pm, err := newPatternMatcher(g, cfg)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, pm)
	}
	return &Pattern{raw: expr, matchers: matchers}, nil
}

// MustNewPattern is like NewPattern but panics when parsing fails.
func MustNewPattern(expr string) *Pattern {
	p, err := NewPatternWithOptions(expr)
	if err != nil {
		panic(err)
	}
	return p
}

// MustNewPatternWithOptions is like NewPatternWithOptions but panics when parsing fails.
func MustNewPatternWithOptions(expr string, opts ...PatternOption) *Pattern {
	p, err := NewPatternWithOptions(expr, opts...)
	if err != nil {
		panic(err)
	}
	return p
}

// Raw returns the original expression.
func (p *Pattern) Raw() string { return p.raw }

// Match reports whether the given time matches; when matchSecond is false, seconds are ignored.
func (p *Pattern) Match(t time.Time, matchSecond bool) bool {
	fields := timeToFields(t, matchSecond)
	for _, pm := range p.matchers {
		if pm.match(fields) {
			return true
		}
	}
	return false
}

// timeToFields converts time.Time to [second, minute, hour, dayOfMonth, month, dayOfWeek, year].
// Weekday uses Sunday = 0 through Saturday = 6; month starts at 1; second is -1 when ignored.
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

package base

import (
	"errors"
	"strings"
	"time"
)

// 对应 hutool-core DateUtil。

// 常用日期格式。
const (
	NormPattern         = "2006-01-02 15:04:05"
	NormDatePattern     = "2006-01-02"
	NormTimePattern     = "15:04:05"
	NormDatetimePattern = NormPattern
	PureDatePattern     = "20060102"
	PureDatetimePattern = "20060102150405"
	HTTPPattern         = time.RFC1123
	UTCPattern          = "2006-01-02T15:04:05Z"
)

// Now 当前时间。
func Now() time.Time { return time.Now() }

// Today 当前日期 0 点。
func Today() time.Time { return BeginOfDay(time.Now()) }

// FormatDate 按指定 layout 格式化（layout 使用 Go 的参考时间格式）。
func FormatDate(t time.Time, layout string) string {
	if layout == "" {
		layout = NormPattern
	}
	return t.Format(layout)
}

// FormatDateNorm 标准格式（yyyy-MM-dd HH:mm:ss）。
func FormatDateNorm(t time.Time) string { return t.Format(NormPattern) }

// FormatDateOnly 仅日期（yyyy-MM-dd）。
func FormatDateOnly(t time.Time) string { return t.Format(NormDatePattern) }

// FormatTimeOnly 仅时间（HH:mm:ss）。
func FormatTimeOnly(t time.Time) string { return t.Format(NormTimePattern) }

// ParseDate 自动识别若干常见格式解析；按本地时区。
func ParseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("empty date string")
	}
	patterns := []string{
		NormPattern,
		NormDatePattern,
		NormTimePattern,
		PureDatetimePattern,
		PureDatePattern,
		UTCPattern,
		time.RFC3339,
		time.RFC1123,
		"2006/01/02 15:04:05",
		"2006/01/02",
		"2006-01-02T15:04:05",
	}
	for _, p := range patterns {
		if t, err := time.ParseInLocation(p, s, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unsupported date format: " + s)
}

// ParseDateLayout 按指定 layout 解析。
func ParseDateLayout(s, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, s, time.Local)
}

// BeginOfDay 当天 0 点。
func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 当天 23:59:59.999999999。
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfMonth 月初 0 点。
func BeginOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 月末 23:59:59.999999999。
func EndOfMonth(t time.Time) time.Time {
	first := BeginOfMonth(t)
	return EndOfDay(first.AddDate(0, 1, -1))
}

// BeginOfYear 年初。
func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 年末。
func EndOfYear(t time.Time) time.Time {
	return EndOfDay(time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location()))
}

// OffsetDay 偏移天数。
func OffsetDay(t time.Time, days int) time.Time { return t.AddDate(0, 0, days) }

// OffsetMonth 偏移月份。
func OffsetMonth(t time.Time, months int) time.Time { return t.AddDate(0, months, 0) }

// OffsetYear 偏移年份。
func OffsetYear(t time.Time, years int) time.Time { return t.AddDate(years, 0, 0) }

// OffsetHour 偏移小时。
func OffsetHour(t time.Time, hours int) time.Time { return t.Add(time.Duration(hours) * time.Hour) }

// OffsetMinute 偏移分钟。
func OffsetMinute(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// OffsetSecond 偏移秒。
func OffsetSecond(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// BetweenDays 两个时间之间的天数（绝对值）。
func BetweenDays(a, b time.Time) int {
	d := b.Sub(a) / (24 * time.Hour)
	if d < 0 {
		d = -d
	}
	return int(d)
}

// IsSameDay 是否同一天（按本地时区）。
func IsSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

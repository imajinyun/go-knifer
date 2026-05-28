package log

import (
	"fmt"
	"strings"
)

// formatTemplate 兼容 hutool 风格的 "{}" 占位符。
// 当 format 中包含 "{}" 时，按位置依次替换为参数；否则使用 fmt.Sprintf。
//
//	formatTemplate("hello {}", "world")  // "hello world"
//	formatTemplate("a=%d", 1)            // "a=1"
func formatTemplate(format string, args ...any) string {
	if format == "" {
		return concatArgs(args...)
	}
	if strings.Contains(format, "{}") {
		return replacePlaceholders(format, args...)
	}
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

// replacePlaceholders 按 {} 顺序替换参数。
func replacePlaceholders(format string, args ...any) string {
	var b strings.Builder
	b.Grow(len(format))
	idx := 0
	for i := 0; i < len(format); i++ {
		if i+1 < len(format) && format[i] == '{' && format[i+1] == '}' {
			if idx < len(args) {
				b.WriteString(fmt.Sprint(args[idx]))
				idx++
			} else {
				b.WriteString("{}")
			}
			i++
			continue
		}
		b.WriteByte(format[i])
	}
	return b.String()
}

// concatArgs 将参数拼接为单一字符串（空格分隔），与 fmt.Sprint 行为一致。
func concatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	return fmt.Sprint(args...)
}

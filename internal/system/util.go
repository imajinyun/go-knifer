package system

import (
	"fmt"
	"strings"
)

// appendLine 与 hutool SystemUtil.append 对应，将形如 "Caption: value\n" 写入 builder。
// 当 value 为空字符串时，使用 [n/a] 占位。
func appendLine(b *strings.Builder, caption string, value any) {
	s := toStr(value)
	if s == "" {
		s = "[n/a]"
	}
	b.WriteString(caption)
	b.WriteString(s)
	b.WriteByte('\n')
}

// toStr 将任意值转为字符串。
func toStr(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%v", x)
	}
}

// readableSize 将字节数转为可读字符串，例如 "1.2 GB"。
func readableSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), units[exp])
}

// addSuffixIfNot 如果 s 不以 suffix 结尾，则追加 suffix。
func addSuffixIfNot(s, suffix string) string {
	if s == "" {
		return s
	}
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

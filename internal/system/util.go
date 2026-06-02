package system

import (
	"fmt"
	"strings"
)

// appendLine writes a line in the form "Caption: value\n" to the builder.
// It uses [n/a] as a placeholder when value is an empty string.
func appendLine(b *strings.Builder, caption string, value any) {
	s := toStr(value)
	if s == "" {
		s = "[n/a]"
	}
	b.WriteString(caption)
	b.WriteString(s)
	b.WriteByte('\n')
}

// toStr converts any value to a string.
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

// readableSize converts bytes to a human-readable string, such as "1.2 GB".
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

// addSuffixIfNot appends suffix when s does not already end with it.
func addSuffixIfNot(s, suffix string) string {
	if s == "" {
		return s
	}
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

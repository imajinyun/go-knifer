package json

import "strings"

// FormatJSONStr 对应 hutool 的 JSONStrFormatter.format，
// 对原始 JSON 字符串进行格式化（4 空格缩进），不构造对象树。
func FormatJSONStr(raw string) string {
	const indentStr = "    "
	var sb strings.Builder
	level := 0
	inString := false
	prevEsc := false

	writeIndent := func() {
		for i := 0; i < level; i++ {
			sb.WriteString(indentStr)
		}
	}

	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if inString {
			sb.WriteByte(c)
			if c == '\\' && !prevEsc {
				prevEsc = true
				continue
			}
			if c == '"' && !prevEsc {
				inString = false
			}
			prevEsc = false
			continue
		}
		switch c {
		case '"':
			inString = true
			sb.WriteByte(c)
		case '{', '[':
			sb.WriteByte(c)
			// 空容器原样输出
			if i+1 < len(raw) && (raw[i+1] == '}' || raw[i+1] == ']') {
				continue
			}
			level++
			sb.WriteByte('\n')
			writeIndent()
		case '}', ']':
			// 移除尾随空白
			s := strings.TrimRight(sb.String(), " ")
			sb.Reset()
			sb.WriteString(s)
			if !strings.HasSuffix(s, "{") && !strings.HasSuffix(s, "[") {
				level--
				sb.WriteByte('\n')
				writeIndent()
			}
			sb.WriteByte(c)
		case ',':
			sb.WriteByte(c)
			sb.WriteByte('\n')
			writeIndent()
		case ':':
			sb.WriteByte(c)
			sb.WriteByte(' ')
		case ' ', '\t', '\n', '\r':
			// 忽略原有空白
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

package conf

import (
	"bufio"
	"strings"
)

// ParseTOML parses common TOML key-value and section syntax into grouped configuration.
func ParseTOML(content string) (*Conf, error) {
	c := New()
	group := defaultGroup
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(stripInlineComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = strings.TrimSpace(line[1 : len(line)-1])
			if group == "" || strings.HasPrefix(group, "[") || strings.HasSuffix(group, "]") {
				return nil, invalidInputf("invalid toml section at line %d: %s", lineNo, line)
			}
			c.ensureGroup(group)
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			return nil, invalidInputf("invalid toml line %d: %s", lineNo, line)
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if key == "" {
			return nil, invalidInputf("empty toml key at line %d", lineNo)
		}
		c.SetByGroup(group, strings.Trim(key, `"'`), normalizeScalar(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, wrapConfigParse("scan toml content", err)
	}
	return c, nil
}

func stripInlineComment(line string) string {
	inSingle, inDouble := false, false
	for i, r := range line {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return line[:i]
			}
		}
	}
	return line
}

func normalizeScalar(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		inner := strings.TrimSpace(value[1 : len(value)-1])
		if inner == "" {
			return ""
		}
		parts := strings.Split(inner, ",")
		for i := range parts {
			parts[i] = normalizeScalar(parts[i])
		}
		return strings.Join(parts, ",")
	}
	return unquote(value)
}

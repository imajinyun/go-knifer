package base

import (
	"strings"
	"unicode"
)

// This file provides naming-case conversion helpers aligned with hutool-core NamingCase.

// ToCamelCase converts separators to lower camel case, for example hello_world -> helloWorld.
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}
	if !strings.ContainsAny(s, "_- ") {
		// Already camel/Pascal-like; only lower-case the first rune.
		r := []rune(s)
		r[0] = unicode.ToLower(r[0])
		return string(r)
	}
	var b strings.Builder
	b.Grow(len(s))
	upper := false
	first := true
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			upper = true
			continue
		}
		if first {
			b.WriteRune(unicode.ToLower(r))
			first = false
			continue
		}
		if upper {
			b.WriteRune(unicode.ToUpper(r))
			upper = false
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// ToPascalCase converts separators to upper camel case, for example hello_world -> HelloWorld.
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}
	c := ToCamelCase(s)
	r := []rune(c)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// ToUnderlineCase converts camel case or separators to snake case, for example HelloWorld -> hello_world.
func ToUnderlineCase(s string) string { return toSeparated(s, '_') }

// ToKebabCase converts camel case or separators to kebab case, for example HelloWorld -> hello-world.
func ToKebabCase(s string) string { return toSeparated(s, '-') }

func toSeparated(s string, sep rune) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 4)
	rs := []rune(s)
	for i, r := range rs {
		if r == '_' || r == '-' || r == ' ' {
			b.WriteRune(sep)
			continue
		}
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rs[i-1]
				// Add a separator before acronym boundaries or after lower-case/digit runs.
				prevIsLower := unicode.IsLower(prev) || unicode.IsDigit(prev)
				nextIsLower := i+1 < len(rs) && unicode.IsLower(rs[i+1])
				if prevIsLower || (unicode.IsUpper(prev) && nextIsLower) {
					b.WriteRune(sep)
				}
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

package json

import "strings"

type formatConfig struct {
	indent        string
	spaceAfterKey bool
}

// FormatOption customizes raw JSON string formatting.
type FormatOption func(*formatConfig)

// WithFormatIndent sets the indentation string used by FormatJSONStrWithOptions.
func WithFormatIndent(indent string) FormatOption {
	return func(c *formatConfig) { c.indent = indent }
}

// WithFormatIndentWidth sets indentation to n spaces.
func WithFormatIndentWidth(n int) FormatOption {
	return func(c *formatConfig) {
		if n < 0 {
			n = 0
		}
		c.indent = strings.Repeat(" ", n)
	}
}

// WithFormatSpaceAfterKey controls whether a space is written after ':'.
func WithFormatSpaceAfterKey(space bool) FormatOption {
	return func(c *formatConfig) { c.spaceAfterKey = space }
}

func applyFormatOptions(opts []FormatOption) formatConfig {
	cfg := formatConfig{indent: "    ", spaceAfterKey: true}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// FormatJSONStr matches the utility JSONStrFormatter.format,
// formatting a raw JSON string with 4-space indentation without building an object tree.
func FormatJSONStr(raw string) string {
	return FormatJSONStrWithOptions(raw)
}

// FormatJSONStrWithOptions formats raw JSON using custom formatting options.
func FormatJSONStrWithOptions(raw string, opts ...FormatOption) string {
	cfg := applyFormatOptions(opts)
	var sb strings.Builder
	level := 0
	inString := false
	prevEsc := false

	writeIndent := func() {
		for i := 0; i < level; i++ {
			sb.WriteString(cfg.indent)
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
			// Keep empty containers unchanged.
			if i+1 < len(raw) && (raw[i+1] == '}' || raw[i+1] == ']') {
				continue
			}
			level++
			sb.WriteByte('\n')
			writeIndent()
		case '}', ']':
			// Remove trailing whitespace.
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
			if cfg.spaceAfterKey {
				sb.WriteByte(' ')
			}
		case ' ', '\t', '\n', '\r':
			// Ignore original whitespace.
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

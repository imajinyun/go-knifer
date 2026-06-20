package cli

import (
	"fmt"
	"strings"
)

// ColorMode controls ANSI color output.
type ColorMode int

const (
	// ColorAuto enables color output for callers that opt into automatic behavior.
	ColorAuto ColorMode = iota
	// ColorAlways always emits ANSI color escape sequences.
	ColorAlways
	// ColorNever never emits ANSI color escape sequences.
	ColorNever
)

// Color names supported ANSI foreground colors.
type Color string

const (
	// Red is the ANSI red foreground color.
	Red Color = "31"
	// Green is the ANSI green foreground color.
	Green Color = "32"
	// Yellow is the ANSI yellow foreground color.
	Yellow Color = "33"
	// Blue is the ANSI blue foreground color.
	Blue Color = "34"
	// Bold is the ANSI bold text attribute.
	Bold Color = "1"
)

type colorConfig struct {
	mode ColorMode
}

// ColorOption customizes color rendering.
type ColorOption func(*colorConfig)

// WithColorMode sets ANSI color behavior.
func WithColorMode(mode ColorMode) ColorOption {
	return func(c *colorConfig) { c.mode = mode }
}

func applyColorOptions(opts []ColorOption) colorConfig {
	cfg := colorConfig{mode: ColorAuto}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// Colorize wraps text in ANSI escape codes when enabled.
func Colorize(text string, color Color, opts ...ColorOption) string {
	cfg := applyColorOptions(opts)
	if cfg.mode != ColorAlways || color == "" {
		return text
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, text)
}

// RenderHelp returns deterministic help text for cmd.
func RenderHelp(cmd *Command, opts ...ColorOption) string {
	if cmd == nil {
		return ""
	}
	usage := cmd.Usage
	if usage == "" {
		usage = cmd.Name
		if len(cmd.Children) > 0 {
			usage += " <command>"
		}
	}
	var b strings.Builder
	b.WriteString(Colorize("Usage:", Bold, opts...))
	b.WriteByte(' ')
	b.WriteString(usage)
	b.WriteByte('\n')
	if cmd.Summary != "" {
		b.WriteByte('\n')
		b.WriteString(cmd.Summary)
		b.WriteByte('\n')
	}
	if len(cmd.Children) > 0 {
		b.WriteString("\nCommands:\n")
		for _, child := range cmd.Children {
			if child == nil {
				continue
			}
			b.WriteString("  ")
			b.WriteString(child.Name)
			if child.Summary != "" {
				b.WriteString("\t")
				b.WriteString(child.Summary)
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}

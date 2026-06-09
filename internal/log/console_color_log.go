package log

import (
	"fmt"
	"sync"
)

// ANSI color codes.
const (
	ansiReset   = "\033[0m"
	ansiDefault = "\033[39m"
	ansiBlack   = "\033[30m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiWhite   = "\033[37m"
)

// ColorFactory returns the ANSI color code for a level.
type ColorFactory func(level Level) string

var (
	colorFactoryMu sync.RWMutex
	colorFactory   ColorFactory = defaultColorFactory
)

func defaultColorFactory(level Level) string {
	switch level {
	case LevelDebug, LevelInfo:
		return ansiGreen
	case LevelWarn:
		return ansiYellow
	case LevelError, LevelFatal:
		return ansiRed
	case LevelTrace:
		return ansiMagenta
	default:
		return ansiDefault
	}
}

// SetColorFactory customizes the color factory.
func SetColorFactory(f ColorFactory) {
	if f == nil {
		return
	}
	colorFactoryMu.Lock()
	defer colorFactoryMu.Unlock()
	colorFactory = f
}

// getColorFactory returns the current color factory.
func getColorFactory() ColorFactory {
	colorFactoryMu.RLock()
	defer colorFactoryMu.RUnlock()
	return colorFactory
}

// ConsoleColorLog matches the utility ConsoleColorLog and prints logs with ANSI colors.
type ConsoleColorLog struct {
	*ConsoleLog
	colorFactory ColorFactory
}

// NewConsoleColorLog creates a colored console log instance.
func NewConsoleColorLog(name string) *ConsoleColorLog {
	return NewConsoleColorLogWithOptions(name)
}

// NewConsoleColorLogWithOptions creates a colored console log instance and applies constructor options.
func NewConsoleColorLogWithOptions(name string, opts ...ConsoleLogOption) *ConsoleColorLog {
	base := NewConsoleLogWithOptions(name, opts...)
	factory := base.colorFactory
	if factory == nil {
		factory = getColorFactory()
	}
	c := &ConsoleColorLog{ConsoleLog: base, colorFactory: factory}
	// Replace Core with the colored implementation.
	base.Core = c.write
	return c
}

// write is the colored write implementation.
func (c *ConsoleColorLog) write(level Level, err error, format string, args ...any) {
	msg := renderLogMessage(format, args...)
	factory := c.colorFactory
	if factory == nil {
		factory = defaultColorFactory
	}
	color := factory(level)
	line := fmt.Sprintf(
		"%s[%s]%s %s[%-5s]%s %s%s%s: %s",
		ansiWhite, c.now().Format(c.layout()), ansiReset,
		color, level.String(), ansiReset,
		ansiCyan, c.name, ansiReset,
		msg,
	)
	if err != nil {
		line = line + " | error: " + err.Error()
	}
	w := c.targetWriter(level)
	_, _ = fmt.Fprintln(w, line)
}

package log

// Level is the log level matching the utility toolkit log.level.Level.
// Lower values are more detailed: All < Trace < Debug < Info < Warn < Error < Fatal < Off.
type Level int

const (
	// LevelAll allows all levels and is the lowest filtering threshold.
	LevelAll Level = iota
	// LevelTrace trace level.
	LevelTrace
	// LevelDebug debug level.
	LevelDebug
	// LevelInfo info level.
	LevelInfo
	// LevelWarn warn level.
	LevelWarn
	// LevelError error level.
	LevelError
	// LevelFatal fatal level.
	LevelFatal
	// LevelOff disables logging and is the highest filtering threshold.
	LevelOff
)

// String returns the human-readable level name.
func (l Level) String() string {
	switch l {
	case LevelAll:
		return "ALL"
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	case LevelOff:
		return "OFF"
	default:
		return "UNKNOWN"
	}
}

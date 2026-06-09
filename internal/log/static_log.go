package log

// StaticLog provides static convenience methods using the default Log instance named "static".
//
// Because Go has no Java-style CallerUtil, this uses a fixed Log name.
// Use Get(name) for a dedicated Log when caller-based naming is needed.

const staticLogName = "static"

func staticLogger() Log {
	return Get(staticLogName)
}

func staticLoggerWithOptions(opts ...LoggerOption) Log {
	return GetWithOptions(staticLogName, opts...)
}

// Trace prints a trace-level log.
func Trace(args ...any) { staticLogger().Trace(args...) }

// Tracef prints a trace-level log with a template and supports "{}" placeholders.
func Tracef(format string, args ...any) { staticLogger().Tracef(format, args...) }

// Debug prints a debug-level log.
func Debug(args ...any) { staticLogger().Debug(args...) }

// Debugf prints a debug-level log with a template.
func Debugf(format string, args ...any) { staticLogger().Debugf(format, args...) }

// Info prints an info-level log.
func Info(args ...any) { staticLogger().Info(args...) }

// Infof prints an info-level log with a template.
func Infof(format string, args ...any) { staticLogger().Infof(format, args...) }

// Warn prints a warn-level log.
func Warn(args ...any) { staticLogger().Warn(args...) }

// Warnf prints a warn-level log with a template.
func Warnf(format string, args ...any) { staticLogger().Warnf(format, args...) }

// ErrorLog prints an error-level log; named to avoid colliding with the built-in error type.
func ErrorLog(args ...any) { staticLogger().Error(args...) }

// Errorf prints an error-level log with a template.
func Errorf(format string, args ...any) { staticLogger().Errorf(format, args...) }

// LogAt prints a log at the specified level.
func LogAt(level Level, format string, args ...any) {
	staticLogger().Log(level, format, args...)
}

// LogAtE prints a log at the specified level with an error object.
func LogAtE(level Level, err error, format string, args ...any) {
	staticLogger().LogE(level, err, format, args...)
}

// TraceWithOptions prints trace-level output through a per-call logger configuration.
func TraceWithOptions(opts []LoggerOption, args ...any) {
	staticLoggerWithOptions(opts...).Trace(args...)
}

// TracefWithOptions prints formatted trace-level output through a per-call logger configuration.
func TracefWithOptions(opts []LoggerOption, format string, args ...any) {
	staticLoggerWithOptions(opts...).Tracef(format, args...)
}

// DebugWithOptions prints debug-level output through a per-call logger configuration.
func DebugWithOptions(opts []LoggerOption, args ...any) {
	staticLoggerWithOptions(opts...).Debug(args...)
}

// DebugfWithOptions prints formatted debug-level output through a per-call logger configuration.
func DebugfWithOptions(opts []LoggerOption, format string, args ...any) {
	staticLoggerWithOptions(opts...).Debugf(format, args...)
}

// InfoWithOptions prints info-level output through a per-call logger configuration.
func InfoWithOptions(opts []LoggerOption, args ...any) {
	staticLoggerWithOptions(opts...).Info(args...)
}

// InfofWithOptions prints formatted info-level output through a per-call logger configuration.
func InfofWithOptions(opts []LoggerOption, format string, args ...any) {
	staticLoggerWithOptions(opts...).Infof(format, args...)
}

// WarnWithOptions prints warn-level output through a per-call logger configuration.
func WarnWithOptions(opts []LoggerOption, args ...any) {
	staticLoggerWithOptions(opts...).Warn(args...)
}

// WarnfWithOptions prints formatted warn-level output through a per-call logger configuration.
func WarnfWithOptions(opts []LoggerOption, format string, args ...any) {
	staticLoggerWithOptions(opts...).Warnf(format, args...)
}

// ErrorLogWithOptions prints error-level output through a per-call logger configuration.
func ErrorLogWithOptions(opts []LoggerOption, args ...any) {
	staticLoggerWithOptions(opts...).Error(args...)
}

// ErrorfWithOptions prints formatted error-level output through a per-call logger configuration.
func ErrorfWithOptions(opts []LoggerOption, format string, args ...any) {
	staticLoggerWithOptions(opts...).Errorf(format, args...)
}

// LogAtWithOptions logs output at the provided level through a per-call logger configuration.
func LogAtWithOptions(opts []LoggerOption, level Level, format string, args ...any) {
	staticLoggerWithOptions(opts...).Log(level, format, args...)
}

// LogAtEWithOptions logs output at the provided level with an error through a per-call logger configuration.
func LogAtEWithOptions(opts []LoggerOption, level Level, err error, format string, args ...any) {
	staticLoggerWithOptions(opts...).LogE(level, err, format, args...)
}

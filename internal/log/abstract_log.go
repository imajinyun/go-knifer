package log

// AbstractLog provides a convenience implementation of the Log interface; subtypes only need to implement LogCore, namely LogE.
// avoids repeated boilerplate in each concrete implementation by delegating Log interface methods to LogE.
type AbstractLog struct {
	// Core is provided by the concrete implementation to print logs at the specified level.
	// AbstractLog skips the call when IsEnabled(level) is false.
	Core func(level Level, err error, format string, args ...any)
	// IsEnabledFn is provided by the concrete implementation to decide whether the specified level is enabled.
	IsEnabledFn func(level Level) bool
}

// IsEnabled checks through the injected function.
func (a *AbstractLog) IsEnabled(level Level) bool {
	if a.IsEnabledFn == nil {
		return true
	}
	return a.IsEnabledFn(level)
}

// IsTraceEnabled and related convenience methods.
func (a *AbstractLog) IsTraceEnabled() bool { return a.IsEnabled(LevelTrace) }
func (a *AbstractLog) IsDebugEnabled() bool { return a.IsEnabled(LevelDebug) }
func (a *AbstractLog) IsInfoEnabled() bool  { return a.IsEnabled(LevelInfo) }
func (a *AbstractLog) IsWarnEnabled() bool  { return a.IsEnabled(LevelWarn) }
func (a *AbstractLog) IsErrorEnabled() bool { return a.IsEnabled(LevelError) }

// LogE is the common entry point for printing logs with errors by level.
func (a *AbstractLog) LogE(level Level, err error, format string, args ...any) {
	if !a.IsEnabled(level) {
		return
	}
	if a.Core != nil {
		a.Core(level, err, format, args...)
	}
}

// Log is the common entry point for printing logs by level.
func (a *AbstractLog) Log(level Level, format string, args ...any) {
	a.LogE(level, nil, format, args...)
}

// Trace logs at the trace level; the following methods are level shortcuts.
func (a *AbstractLog) Trace(args ...any)                 { a.Log(LevelTrace, "", args...) }
func (a *AbstractLog) Tracef(format string, args ...any) { a.Log(LevelTrace, format, args...) }
func (a *AbstractLog) Debug(args ...any)                 { a.Log(LevelDebug, "", args...) }
func (a *AbstractLog) Debugf(format string, args ...any) { a.Log(LevelDebug, format, args...) }
func (a *AbstractLog) Info(args ...any)                  { a.Log(LevelInfo, "", args...) }
func (a *AbstractLog) Infof(format string, args ...any)  { a.Log(LevelInfo, format, args...) }
func (a *AbstractLog) Warn(args ...any)                  { a.Log(LevelWarn, "", args...) }
func (a *AbstractLog) Warnf(format string, args ...any)  { a.Log(LevelWarn, format, args...) }
func (a *AbstractLog) Error(args ...any)                 { a.Log(LevelError, "", args...) }
func (a *AbstractLog) Errorf(format string, args ...any) { a.Log(LevelError, format, args...) }

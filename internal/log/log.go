package log

// Log matches the utility Log interface and provides unified logging methods.
//
// Unlike Java, Go has no native overloads, so methods with formatting templates use the *f suffix and are compatible with fmt.Sprintf:
//   - Trace(args ...any)        / Tracef(format string, args ...any)
//   - Debug / Debugf
//   - Info  / Infof
//   - Warn  / Warnf
//   - Error / Errorf
//
// Additional methods:
//   - Log(level, format, args...) common logging entry point;
//   - LogE(level, err, format, args...) includes an error object.
type Log interface {
	// GetName returns the log name, usually a class, package, or type name.
	GetName() string

	// IsEnabled reports whether the specified level is enabled.
	IsEnabled(level Level) bool

	// IsTraceEnabled / IsDebugEnabled / ... and related convenience checks.
	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool

	Trace(args ...any)
	Tracef(format string, args ...any)

	Debug(args ...any)
	Debugf(format string, args ...any)

	Info(args ...any)
	Infof(format string, args ...any)

	Warn(args ...any)
	Warnf(format string, args ...any)

	Error(args ...any)
	Errorf(format string, args ...any)

	// Log prints a log at the specified level.
	Log(level Level, format string, args ...any)
	// LogE prints a log at the specified level with an error object.
	LogE(level Level, err error, format string, args ...any)
}

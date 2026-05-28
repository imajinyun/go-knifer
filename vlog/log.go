package vlog

import logx "github.com/imajinyun/go-knifer/internal/log"

// LogLevel is the logging level used by the logging facade.
type LogLevel = logx.Level

// Level is the logging level used by the logging facade.
type Level = logx.Level

// Log is the common logging interface.
type Log = logx.Log

// LogFactory creates loggers by name.
type LogFactory = logx.LogFactory

// LogFactoryFunc adapts a function into a LogFactory.
type LogFactoryFunc = logx.LogFactoryFunc

// ConsoleLog is the default console logger implementation.
type ConsoleLog = logx.ConsoleLog

// ConsoleColorLog is a console logger with ANSI colors.
type ConsoleColorLog = logx.ConsoleColorLog

// ColorFactory maps log levels to ANSI colors.
type ColorFactory = logx.ColorFactory

// AbstractLog provides a reusable Log implementation skeleton.
type AbstractLog = logx.AbstractLog

const (
	// LogLevelAll enables every log level.
	LogLevelAll LogLevel = logx.LevelAll
	// LogLevelTrace is the trace level.
	LogLevelTrace LogLevel = logx.LevelTrace
	// LogLevelDebug is the debug level.
	LogLevelDebug LogLevel = logx.LevelDebug
	// LogLevelInfo is the info level.
	LogLevelInfo LogLevel = logx.LevelInfo
	// LogLevelWarn is the warn level.
	LogLevelWarn LogLevel = logx.LevelWarn
	// LogLevelError is the error level.
	LogLevelError LogLevel = logx.LevelError
	// LogLevelFatal is the fatal level.
	LogLevelFatal LogLevel = logx.LevelFatal
	// LogLevelOff disables logging.
	LogLevelOff LogLevel = logx.LevelOff
)

// NewConsoleLog creates a console logger by name.
func NewConsoleLog(name string) *ConsoleLog { return logx.NewConsoleLog(name) }

// NewConsoleColorLog creates a colored console logger by name.
func NewConsoleColorLog(name string) *ConsoleColorLog { return logx.NewConsoleColorLog(name) }

// Logger returns a cached logger by name.
func Logger(name string) Log { return logx.Get(name) }

// DefaultLogger returns the default logger.
func DefaultLogger() Log { return logx.GetDefault() }

// SetLogFactory sets the global logging factory.
func SetLogFactory(factory LogFactory) { logx.SetFactory(factory) }

// SetLogLevel sets the console logging threshold.
func SetLogLevel(level LogLevel) { logx.SetConsoleLevel(level) }

// GetLogLevel returns the console logging threshold.
func GetLogLevel() LogLevel { return logx.GetConsoleLevel() }

// SetLogColorFactory customizes console log colors.
func SetLogColorFactory(f ColorFactory) { logx.SetColorFactory(f) }

// Trace logs trace-level output through the static logger.
func Trace(args ...any) { logx.Trace(args...) }

// Tracef logs formatted trace-level output through the static logger.
func Tracef(format string, args ...any) { logx.Tracef(format, args...) }

// Debug logs debug-level output through the static logger.
func Debug(args ...any) { logx.Debug(args...) }

// Debugf logs formatted debug-level output through the static logger.
func Debugf(format string, args ...any) { logx.Debugf(format, args...) }

// Info logs info-level output through the static logger.
func Info(args ...any) { logx.Info(args...) }

// Infof logs formatted info-level output through the static logger.
func Infof(format string, args ...any) { logx.Infof(format, args...) }

// Warn logs warn-level output through the static logger.
func Warn(args ...any) { logx.Warn(args...) }

// Warnf logs formatted warn-level output through the static logger.
func Warnf(format string, args ...any) { logx.Warnf(format, args...) }

// ErrorLog logs error-level output through the static logger.
func ErrorLog(args ...any) { logx.ErrorLog(args...) }

// Errorf logs formatted error-level output through the static logger.
func Errorf(format string, args ...any) { logx.Errorf(format, args...) }

// LogAt logs output at the provided level through the static logger.
func LogAt(level LogLevel, format string, args ...any) { logx.LogAt(level, format, args...) }

// LogAtE logs output at the provided level with an error through the static logger.
func LogAtE(level LogLevel, err error, format string, args ...any) {
	logx.LogAtE(level, err, format, args...)
}

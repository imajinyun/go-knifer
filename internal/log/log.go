package log

// Log 对应 the utility Log 接口，提供统一的日志方法。
//
// 与 Java 不同，Go 中没有原生重载，因此采用 *f 后缀区分带格式化模板的方法（与 fmt.Sprintf 兼容）：
//   - Trace(args ...any)        / Tracef(format string, args ...any)
//   - Debug / Debugf
//   - Info  / Infof
//   - Warn  / Warnf
//   - Error / Errorf
//
// 此外提供：
//   - Log(level, format, args...) 通用日志入口；
//   - LogE(level, err, format, args...) 携带错误对象。
type Log interface {
	// GetName 返回日志名称（通常是类/包/类型名）。
	GetName() string

	// IsEnabled 判断指定级别是否开启。
	IsEnabled(level Level) bool

	// IsTraceEnabled / IsDebugEnabled / ... 等便捷判断。
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

	// Log 打印指定级别的日志。
	Log(level Level, format string, args ...any)
	// LogE 打印指定级别的日志并附带错误对象。
	LogE(level Level, err error, format string, args ...any)
}

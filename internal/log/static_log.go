package log

// StaticLog 提供静态便捷方法，使用名为 "static" 的默认 Log 实例。
//
// 由于 Go 没有 Java 那样的 CallerUtil，这里使用固定名称的 Log。
// 若需要按调用者定位，请使用 Get(name) 获取专用 Log。

const staticLogName = "static"

func staticLogger() Log {
	return Get(staticLogName)
}

// Trace 打印 trace 级别日志。
func Trace(args ...any) { staticLogger().Trace(args...) }

// Tracef 按模板打印 trace 级别日志，支持 "{}" 占位符。
func Tracef(format string, args ...any) { staticLogger().Tracef(format, args...) }

// Debug 打印 debug 级别日志。
func Debug(args ...any) { staticLogger().Debug(args...) }

// Debugf 按模板打印 debug 级别日志。
func Debugf(format string, args ...any) { staticLogger().Debugf(format, args...) }

// Info 打印 info 级别日志。
func Info(args ...any) { staticLogger().Info(args...) }

// Infof 按模板打印 info 级别日志。
func Infof(format string, args ...any) { staticLogger().Infof(format, args...) }

// Warn 打印 warn 级别日志。
func Warn(args ...any) { staticLogger().Warn(args...) }

// Warnf 按模板打印 warn 级别日志。
func Warnf(format string, args ...any) { staticLogger().Warnf(format, args...) }

// ErrorLog 打印 error 级别日志（命名避免与内置 error 类型混淆）。
func ErrorLog(args ...any) { staticLogger().Error(args...) }

// Errorf 按模板打印 error 级别日志。
func Errorf(format string, args ...any) { staticLogger().Errorf(format, args...) }

// LogAt 在指定级别打印日志。
func LogAt(level Level, format string, args ...any) {
	staticLogger().Log(level, format, args...)
}

// LogAtE 在指定级别打印日志，并附带错误对象。
func LogAtE(level Level, err error, format string, args ...any) {
	staticLogger().LogE(level, err, format, args...)
}

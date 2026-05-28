package log

// Level 日志等级，对应 hutool log.level.Level。
// 数值越小越详细：All < Trace < Debug < Info < Warn < Error < Fatal < Off。
type Level int

const (
	// LevelAll 所有级别（最低，过滤所有都通过）。
	LevelAll Level = iota
	// LevelTrace 跟踪级别。
	LevelTrace
	// LevelDebug 调试级别。
	LevelDebug
	// LevelInfo 信息级别。
	LevelInfo
	// LevelWarn 警告级别。
	LevelWarn
	// LevelError 错误级别。
	LevelError
	// LevelFatal 致命级别。
	LevelFatal
	// LevelOff 关闭日志（最高，过滤所有都不通过）。
	LevelOff
)

// String 返回日志级别的可读名称。
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

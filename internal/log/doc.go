// Package log 对应 hutool-log，提供统一的日志接口与默认控制台实现。
//
// 主要类型：
//   - Level：日志级别（Trace/Debug/Info/Warn/Error/Fatal/All/Off）。
//   - Log：日志统一接口（Trace/Debug/Info/Warn/Error 等方法）。
//   - LogFactory：根据名称构造 Log。
//   - StaticLog：静态便捷方法，可直接打印日志。
//   - ConsoleLog / ConsoleColorLog：默认控制台实现（可彩色输出）。
package log

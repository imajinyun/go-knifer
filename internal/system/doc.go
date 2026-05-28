// Package gksystem 对应 hutool-system，提供运行时、操作系统、用户、主机等系统信息工具。
//
// 由于 Go 没有 JVM 概念，本包将原 hutool 中的 JvmInfo/JvmSpecInfo/JavaInfo/JavaRuntimeInfo
// 等概念合并为 GoInfo（语言/运行时信息）与 RuntimeInfo（内存与协程信息）。
package system

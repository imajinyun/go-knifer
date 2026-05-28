package system

import (
	"runtime"
	"strings"
)

// GoInfo 对应 hutool 中的 JavaInfo / JavaSpecInfo / JavaRuntimeInfo / JvmInfo / JvmSpecInfo。
// 由于 Go 没有 JVM 概念，这里聚合 Go 运行时与编译信息。
type GoInfo struct {
	Version     string // 例如 go1.22.3
	Compiler    string // gc / gccgo
	GOROOT      string
	GOOS        string
	GOARCH      string
	NumCPU      int
	NumCgoCalls int64
}

// NewGoInfo 构造 GoInfo。
func NewGoInfo() *GoInfo {
	return &GoInfo{
		Version:     runtime.Version(),
		Compiler:    runtime.Compiler,
		GOROOT:      runtime.GOROOT(),
		GOOS:        runtime.GOOS,
		GOARCH:      runtime.GOARCH,
		NumCPU:      runtime.NumCPU(),
		NumCgoCalls: runtime.NumCgoCall(),
	}
}

// GetVersion 取得 Go 版本。
func (g *GoInfo) GetVersion() string { return g.Version }

// GetCompiler 取得编译器名称。
func (g *GoInfo) GetCompiler() string { return g.Compiler }

// GetGOROOT 取得 GOROOT 路径。
func (g *GoInfo) GetGOROOT() string { return g.GOROOT }

// GetGOOS 取得目标 OS。
func (g *GoInfo) GetGOOS() string { return g.GOOS }

// GetGOARCH 取得目标架构。
func (g *GoInfo) GetGOARCH() string { return g.GOARCH }

// GetNumCPU 取得可用 CPU 数。
func (g *GoInfo) GetNumCPU() int { return g.NumCPU }

// String 实现 fmt.Stringer。
func (g *GoInfo) String() string {
	var b strings.Builder
	appendLine(&b, "Go Version:    ", g.Version)
	appendLine(&b, "Go Compiler:   ", g.Compiler)
	appendLine(&b, "GOROOT:        ", g.GOROOT)
	appendLine(&b, "GOOS:          ", g.GOOS)
	appendLine(&b, "GOARCH:        ", g.GOARCH)
	appendLine(&b, "NumCPU:        ", g.NumCPU)
	appendLine(&b, "NumCgoCall:    ", g.NumCgoCalls)
	return b.String()
}

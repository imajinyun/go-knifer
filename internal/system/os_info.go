package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// OsInfo 对应 hutool OsInfo，代表当前操作系统信息。
type OsInfo struct {
	Name          string
	Arch          string
	Version       string
	FileSeparator string
	LineSeparator string
	PathSeparator string
}

// NewOsInfo 构造当前 OS 信息。
func NewOsInfo() *OsInfo {
	return &OsInfo{
		Name:          runtime.GOOS,
		Arch:          runtime.GOARCH,
		Version:       readOsVersion(),
		FileSeparator: string(filepath.Separator),
		LineSeparator: lineSeparator(),
		PathSeparator: string(os.PathListSeparator),
	}
}

// GetName 取得 OS 名称（GOOS）。
func (o *OsInfo) GetName() string { return o.Name }

// GetArch 取得 OS 架构（GOARCH）。
func (o *OsInfo) GetArch() string { return o.Arch }

// GetVersion 取得 OS 版本。
func (o *OsInfo) GetVersion() string { return o.Version }

// GetFileSeparator 文件路径分隔符。
func (o *OsInfo) GetFileSeparator() string { return o.FileSeparator }

// GetLineSeparator 行分隔符。
func (o *OsInfo) GetLineSeparator() string { return o.LineSeparator }

// GetPathSeparator 环境路径分隔符。
func (o *OsInfo) GetPathSeparator() string { return o.PathSeparator }

// IsLinux 判断是否 Linux。
func (o *OsInfo) IsLinux() bool { return o.Name == "linux" }

// IsMac 判断是否 macOS（Darwin）。
func (o *OsInfo) IsMac() bool { return o.Name == "darwin" }

// IsMacOsX 等同于 IsMac，保持与 hutool 命名一致。
func (o *OsInfo) IsMacOsX() bool { return o.IsMac() }

// IsWindows 判断是否 Windows。
func (o *OsInfo) IsWindows() bool { return o.Name == "windows" }

// IsAix 判断是否 AIX。
func (o *OsInfo) IsAix() bool { return o.Name == "aix" }

// IsSolaris 判断是否 Solaris。
func (o *OsInfo) IsSolaris() bool { return o.Name == "solaris" }

// IsFreeBSD 判断是否 FreeBSD。
func (o *OsInfo) IsFreeBSD() bool { return o.Name == "freebsd" }

// String 实现 fmt.Stringer。
func (o *OsInfo) String() string {
	var b strings.Builder
	appendLine(&b, "OS Arch:        ", o.Arch)
	appendLine(&b, "OS Name:        ", o.Name)
	appendLine(&b, "OS Version:     ", o.Version)
	appendLine(&b, "File Separator: ", o.FileSeparator)
	appendLine(&b, "Line Separator: ", o.LineSeparator)
	appendLine(&b, "Path Separator: ", o.PathSeparator)
	return b.String()
}

// lineSeparator 根据当前操作系统返回行分隔符。
func lineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// readOsVersion 通过环境变量 OSTYPE 或常用方式探测 OS 版本，
// 在无法获取时返回空字符串。Go 标准库没有统一 API，因此尽力而为。
func readOsVersion() string {
	if v := os.Getenv("OSVERSION"); v != "" {
		return v
	}
	if v := os.Getenv("OSTYPE"); v != "" {
		return v
	}
	return strings.TrimSpace(runtime.GOOS)
}

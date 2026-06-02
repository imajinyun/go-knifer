package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// OsInfo describes current operating system information.
type OsInfo struct {
	Name          string
	Arch          string
	Version       string
	FileSeparator string
	LineSeparator string
	PathSeparator string
}

// NewOsInfo creates current OS information.
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

// GetName returns the OS name (GOOS).
func (o *OsInfo) GetName() string { return o.Name }

// GetArch returns the OS architecture (GOARCH).
func (o *OsInfo) GetArch() string { return o.Arch }

// GetVersion returns the OS version.
func (o *OsInfo) GetVersion() string { return o.Version }

// GetFileSeparator returns the file path separator.
func (o *OsInfo) GetFileSeparator() string { return o.FileSeparator }

// GetLineSeparator returns the line separator.
func (o *OsInfo) GetLineSeparator() string { return o.LineSeparator }

// GetPathSeparator returns the environment path separator.
func (o *OsInfo) GetPathSeparator() string { return o.PathSeparator }

// IsLinux reports whether the OS is Linux.
func (o *OsInfo) IsLinux() bool { return o.Name == "linux" }

// IsMac reports whether the OS is macOS (Darwin).
func (o *OsInfo) IsMac() bool { return o.Name == "darwin" }

// IsMacOsX is equivalent to IsMac.
func (o *OsInfo) IsMacOsX() bool { return o.IsMac() }

// IsWindows reports whether the OS is Windows.
func (o *OsInfo) IsWindows() bool { return o.Name == "windows" }

// IsAix reports whether the OS is AIX.
func (o *OsInfo) IsAix() bool { return o.Name == "aix" }

// IsSolaris reports whether the OS is Solaris.
func (o *OsInfo) IsSolaris() bool { return o.Name == "solaris" }

// IsFreeBSD reports whether the OS is FreeBSD.
func (o *OsInfo) IsFreeBSD() bool { return o.Name == "freebsd" }

// String implements fmt.Stringer.
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

// lineSeparator returns the line separator for the current OS.
func lineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// readOsVersion detects the OS version from environment variables or common fallbacks.
// The Go standard library has no unified API for this, so this is best-effort.
func readOsVersion() string {
	if v := os.Getenv("OSVERSION"); v != "" {
		return v
	}
	if v := os.Getenv("OSTYPE"); v != "" {
		return v
	}
	return strings.TrimSpace(runtime.GOOS)
}

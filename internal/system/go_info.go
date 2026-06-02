package system

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GoInfo describes Go runtime and compilation metadata.
// It aggregates the metadata that would otherwise be split across Java/JVM info types.
type GoInfo struct {
	Version     string // For example, go1.22.3.
	Compiler    string // gc / gccgo
	GOROOT      string
	GOOS        string
	GOARCH      string
	NumCPU      int
	NumCgoCalls int64
}

// NewGoInfo creates GoInfo.
func NewGoInfo() *GoInfo {
	return &GoInfo{
		Version:     runtime.Version(),
		Compiler:    runtime.Compiler,
		GOROOT:      goRoot(),
		GOOS:        runtime.GOOS,
		GOARCH:      runtime.GOARCH,
		NumCPU:      runtime.NumCPU(),
		NumCgoCalls: runtime.NumCgoCall(),
	}
}

// GetVersion returns the Go version.
func (g *GoInfo) GetVersion() string { return g.Version }

// GetCompiler returns the compiler name.
func (g *GoInfo) GetCompiler() string { return g.Compiler }

// GetGOROOT returns the GOROOT path.
func (g *GoInfo) GetGOROOT() string { return g.GOROOT }

// GetGOOS returns the target OS.
func (g *GoInfo) GetGOOS() string { return g.GOOS }

// GetGOARCH returns the target architecture.
func (g *GoInfo) GetGOARCH() string { return g.GOARCH }

// GetNumCPU returns the number of available CPUs.
func (g *GoInfo) GetNumCPU() int { return g.NumCPU }

// String implements fmt.Stringer.
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

func goRoot() string {
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err == nil {
		if root := strings.TrimSpace(string(out)); root != "" {
			return root
		}
	}
	return os.Getenv("GOROOT")
}

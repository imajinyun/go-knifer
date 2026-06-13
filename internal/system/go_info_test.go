package system

import (
	"runtime"
	"strings"
	"testing"
)

func TestGoInfo(t *testing.T) {
	g := GetGoInfo()
	if g == nil {
		t.Fatal("GoInfo 不应为 nil")
	}
	if g.GetVersion() != runtime.Version() {
		t.Errorf("Go Version 不一致: %s vs %s", g.GetVersion(), runtime.Version())
	}
	if g.GetCompiler() != runtime.Compiler {
		t.Errorf("Compiler 不一致")
	}
	if g.GetNumCPU() != runtime.NumCPU() {
		t.Errorf("NumCPU 不一致")
	}
	if !strings.Contains(g.String(), "Go Version:") {
		t.Errorf("GoInfo.String 缺少 caption")
	}
}

func TestGoInfoWithOptions(t *testing.T) {
	g := NewGoInfoWithOptions(
		WithGoVersionFunc(func() string { return "go-option" }),
		WithGoCompilerFunc(func() string { return "compiler-option" }),
		WithGoRootFunc(func() string { return "/go/root" }),
		WithGoOSFunc(func() string { return "plan9" }),
		WithGoArchFunc(func() string { return "wasm" }),
		WithGoNumCPUFunc(func() int { return 9 }),
		WithGoNumCgoCallFunc(func() int64 { return 10 }),
	)
	if g.GetVersion() != "go-option" || g.GetCompiler() != "compiler-option" || g.GetGOROOT() != "/go/root" || g.GetGOOS() != "plan9" || g.GetGOARCH() != "wasm" || g.GetNumCPU() != 9 || g.NumCgoCalls != 10 {
		t.Fatalf("NewGoInfoWithOptions = %#v", g)
	}
}

func TestSystemInfoGettersWithOptions(t *testing.T) {
	g := GetGoInfoWithOptions(WithGoVersionFunc(func() string { return "go-getter" }))
	if g.GetVersion() != "go-getter" {
		t.Fatalf("GetGoInfoWithOptions version = %q", g.GetVersion())
	}

	o := GetOsInfoWithOptions(WithOSNameFunc(func() string { return "linux" }))
	if o.GetName() != "linux" {
		t.Fatalf("GetOsInfoWithOptions name = %q", o.GetName())
	}
}

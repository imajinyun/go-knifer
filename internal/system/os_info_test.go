package system

import (
	"runtime"
	"testing"
)

func TestOsInfo(t *testing.T) {
	o := GetOsInfo()
	if o == nil {
		t.Fatal("OsInfo 不应为 nil")
	}
	if o.GetName() != runtime.GOOS {
		t.Errorf("OS Name: 期望 %s 实际 %s", runtime.GOOS, o.GetName())
	}
	if o.GetArch() != runtime.GOARCH {
		t.Errorf("OS Arch: 期望 %s 实际 %s", runtime.GOARCH, o.GetArch())
	}
	switch runtime.GOOS {
	case "darwin":
		if !o.IsMac() || !o.IsMacOsX() {
			t.Errorf("darwin 应识别为 Mac")
		}
	case "linux":
		if !o.IsLinux() {
			t.Errorf("linux 应识别为 Linux")
		}
	case "windows":
		if !o.IsWindows() {
			t.Errorf("windows 应识别为 Windows")
		}
	}
	if o.GetFileSeparator() == "" || o.GetPathSeparator() == "" || o.GetLineSeparator() == "" {
		t.Errorf("分隔符不应为空: %+v", o)
	}
}

func TestOsInfoWithOptions(t *testing.T) {
	o := NewOsInfoWithOptions(
		WithOSNameFunc(func() string { return "linux" }),
		WithOSArchFunc(func() string { return "arm64" }),
		WithOSVersionFunc(func() string { return "test-version" }),
		WithOSFileSeparatorFunc(func() string { return "/" }),
		WithOSLineSeparatorFunc(func() string { return "\n" }),
		WithOSPathSeparatorFunc(func() string { return ":" }),
	)
	if o.GetName() != "linux" || o.GetArch() != "arm64" || o.GetVersion() != "test-version" || o.GetFileSeparator() != "/" || o.GetLineSeparator() != "\n" || o.GetPathSeparator() != ":" {
		t.Fatalf("NewOsInfoWithOptions = %#v", o)
	}
	if !o.IsLinux() || o.IsWindows() {
		t.Fatalf("NewOsInfoWithOptions OS helpers = %#v", o)
	}

	o = NewOsInfoWithOptions(
		WithOSNameFunc(func() string { return "windows" }),
		WithOSEnvLookupFunc(func(string) string { return "" }),
	)
	if o.GetVersion() != "windows" || o.GetLineSeparator() != "\r\n" {
		t.Fatalf("OS providers should drive version and line separator: %#v", o)
	}
}

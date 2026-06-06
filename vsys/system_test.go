package vsys_test

import (
	"bytes"
	"net"
	"os"
	"os/user"
	"runtime"
	"testing"

	"github.com/imajinyun/go-knifer/vsys"
)

func TestFacadeHostInfo(t *testing.T) {
	info := vsys.SystemHostInfo()
	if info == nil {
		t.Fatal("expected non-nil host info")
	}

	_, ipNet, err := net.ParseCIDR("10.1.2.3/24")
	if err != nil {
		t.Fatal(err)
	}
	ipNet.IP = net.ParseIP("10.1.2.3")
	info = vsys.SysHostInfoWithOptions(
		vsys.WithHostNameFunc(func() (string, error) { return "facade-host", nil }),
		vsys.WithHostInterfaceAddrsFunc(func() ([]net.Addr, error) { return []net.Addr{ipNet}, nil }),
	)
	if info.GetName() != "facade-host" || info.GetAddress() != "10.1.2.3" {
		t.Fatalf("SysHostInfoWithOptions = %#v", info)
	}

	info = vsys.NewHostInfoWithOptions(vsys.WithHostAddressFunc(func() string { return "198.51.100.2" }))
	if info.GetAddress() != "198.51.100.2" {
		t.Fatalf("NewHostInfoWithOptions address = %q", info.GetAddress())
	}
}

func TestFacadeOsInfo(t *testing.T) {
	info := vsys.SystemOsInfo()
	if info == nil {
		t.Fatal("expected non-nil os info")
	}
	info = vsys.SysOsInfoWithOptions(vsys.WithOSNameFunc(func() string { return "linux" }))
	if info.GetName() != "linux" {
		t.Fatalf("SysOsInfoWithOptions name = %q", info.GetName())
	}
}

func TestFacadeUserInfo(t *testing.T) {
	info := vsys.SystemUserInfo()
	if info == nil {
		t.Fatal("expected non-nil user info")
	}
}

func TestFacadeUserInfoOptions(t *testing.T) {
	info := vsys.SysUserInfoWithOptions(
		vsys.WithCurrentUserFunc(func() (*user.User, error) {
			return &user.User{Username: "facade-user", HomeDir: "/home/facade"}, nil
		}),
		vsys.WithWorkingDirFunc(func() (string, error) { return "/work/facade", nil }),
		vsys.WithTempDirFunc(func() string { return "/tmp/facade" }),
		vsys.WithUserEnvLookup(func(key string) string {
			if key == "LANG" {
				return "zh_CN.UTF-8"
			}
			return ""
		}),
	)
	sep := string(os.PathSeparator)
	if info.GetName() != "facade-user" || info.GetHomeDir() != "/home/facade"+sep || info.GetCurrentDir() != "/work/facade"+sep || info.GetTempDir() != "/tmp/facade"+sep {
		t.Fatalf("SystemUserInfoWithOptions = %#v", info)
	}
	if info.GetLanguage() != "zh" || info.GetCountry() != "CN" {
		t.Fatalf("SysUserInfoWithOptions locale = %s/%s", info.GetLanguage(), info.GetCountry())
	}

	info = vsys.NewUserInfoWithOptions(vsys.WithCurrentUserFunc(func() (*user.User, error) {
		return &user.User{Username: "new-user", HomeDir: "/home/new"}, nil
	}))
	if info.GetName() != "new-user" {
		t.Fatalf("NewUserInfoWithOptions name = %q", info.GetName())
	}
}

func TestFacadeGoInfo(t *testing.T) {
	info := vsys.SystemGoInfo()
	if info == nil {
		t.Fatal("expected non-nil go info")
	}
	info = vsys.SysGoInfoWithOptions(vsys.WithGoVersionFunc(func() string { return "go-sys" }))
	if info.GetVersion() != "go-sys" {
		t.Fatalf("SysGoInfoWithOptions version = %q", info.GetVersion())
	}
	info = vsys.NewGoInfoWithOptions(
		vsys.WithGoVersionFunc(func() string { return "go-facade" }),
		vsys.WithGoCompilerFunc(func() string { return "compiler-facade" }),
		vsys.WithGoRootFunc(func() string { return "/go/facade" }),
		vsys.WithGoOSFunc(func() string { return "linux" }),
		vsys.WithGoArchFunc(func() string { return "arm64" }),
		vsys.WithGoNumCPUFunc(func() int { return 8 }),
		vsys.WithGoNumCgoCallFunc(func() int64 { return 11 }),
	)
	if info.GetVersion() != "go-facade" || info.GetCompiler() != "compiler-facade" || info.GetGOROOT() != "/go/facade" || info.GetGOOS() != "linux" || info.GetGOARCH() != "arm64" || info.GetNumCPU() != 8 || info.NumCgoCalls != 11 {
		t.Fatalf("NewGoInfoWithOptions = %#v", info)
	}
}

func TestFacadeOsInfoOptions(t *testing.T) {
	info := vsys.NewOsInfoWithOptions(
		vsys.WithOSNameFunc(func() string { return "windows" }),
		vsys.WithOSArchFunc(func() string { return "amd64" }),
		vsys.WithOSVersionFunc(func() string { return "11" }),
		vsys.WithOSFileSeparatorFunc(func() string { return "\\" }),
		vsys.WithOSLineSeparatorFunc(func() string { return "\r\n" }),
		vsys.WithOSPathSeparatorFunc(func() string { return ";" }),
	)
	if info.GetName() != "windows" || info.GetArch() != "amd64" || info.GetVersion() != "11" || info.GetFileSeparator() != "\\" || info.GetLineSeparator() != "\r\n" || info.GetPathSeparator() != ";" {
		t.Fatalf("NewOsInfoWithOptions = %#v", info)
	}
	if !info.IsWindows() || info.IsLinux() {
		t.Fatalf("NewOsInfoWithOptions OS helpers = %#v", info)
	}
}

func TestFacadeRuntimeInfo(t *testing.T) {
	info := vsys.SystemRuntimeInfo()
	if info == nil {
		t.Fatal("expected non-nil runtime info")
	}
	info = vsys.SysRuntimeInfoWithOptions(
		vsys.WithReadMemStatsFunc(func(stats *runtime.MemStats) {
			stats.Sys = 4096
			stats.HeapSys = 1024
		}),
		vsys.WithNumGoroutineFunc(func() int { return 5 }),
	)
	if info.GetMaxMemory() != 4096 || info.GetTotalMemory() != 1024 || info.GetGoroutineCount() != 5 {
		t.Fatalf("SysRuntimeInfoWithOptions = %#v", info)
	}

	info = vsys.NewRuntimeInfoWithOptions(vsys.WithReadMemStatsFunc(func(stats *runtime.MemStats) { stats.Sys = 8192 }))
	if info.GetMaxMemory() != 8192 {
		t.Fatalf("NewRuntimeInfoWithOptions max = %d", info.GetMaxMemory())
	}
}

func TestFacadePID(t *testing.T) {
	pid := vsys.CurrentPID()
	if pid <= 0 {
		t.Fatalf("expected positive pid, got %d", pid)
	}
	if got := vsys.CurrentPIDWithOptions(vsys.WithPIDFunc(func() int { return 99 })); got != 99 {
		t.Fatalf("CurrentPIDWithOptions = %d", got)
	}
}

func TestFacadeMemory(t *testing.T) {
	total := vsys.TotalMemory()
	free := vsys.FreeMemory()
	max := vsys.MaxMemory()
	if total == 0 && free == 0 && max == 0 {
		t.Fatal("expected at least one memory metric to be non-zero")
	}
	opt := vsys.WithReadMemStatsFunc(func(stats *runtime.MemStats) {
		stats.Sys = 300
		stats.HeapSys = 200
		stats.HeapIdle = 100
	})
	if vsys.MaxMemoryWithOptions(opt) != 300 || vsys.TotalMemoryWithOptions(opt) != 200 || vsys.FreeMemoryWithOptions(opt) != 100 {
		t.Fatal("expected memory option providers to be used")
	}
}

func TestFacadeGoroutineCount(t *testing.T) {
	count := vsys.TotalGoroutineCount()
	if count < 1 {
		t.Fatalf("expected at least 1 goroutine, got %d", count)
	}
	if got := vsys.TotalGoroutineCountWithOptions(vsys.WithProcessNumGoroutineFunc(func() int { return 12 })); got != 12 {
		t.Fatalf("TotalGoroutineCountWithOptions = %d", got)
	}
}

func TestFacadeEnv(t *testing.T) {
	_ = os.Setenv("GO_KNIFER_TEST_KEY", "test_value")
	defer os.Unsetenv("GO_KNIFER_TEST_KEY")

	if got := vsys.Env("GO_KNIFER_TEST_KEY"); got != "test_value" {
		t.Fatalf("expected 'test_value', got %q", got)
	}
	if got := vsys.EnvOrDefault("GO_KNIFER_TEST_MISSING", "default"); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
	if got := vsys.EnvInt("GO_KNIFER_TEST_KEY", 0); got != 0 {
		t.Fatalf("expected 0 for non-int env, got %d", got)
	}
	if got := vsys.EnvBool("GO_KNIFER_TEST_KEY", false); got != false {
		t.Fatalf("expected false for non-bool env, got %v", got)
	}

	lookup := vsys.WithEnvLookupFunc(func(key string) (string, bool) {
		switch key {
		case "A":
			return "value", true
		case "N":
			return "13", true
		case "B":
			return "true", true
		default:
			return "", false
		}
	})
	var warning bytes.Buffer
	if got := vsys.EnvWithOptions("A", lookup); got != "value" {
		t.Fatalf("EnvWithOptions = %q", got)
	}
	if got := vsys.GetWithOptions("MISSING", false, lookup, vsys.WithEnvWarningWriter(&warning)); got != "" || warning.Len() == 0 {
		t.Fatalf("GetWithOptions missing = %q warning=%q", got, warning.String())
	}
	if got := vsys.EnvOrDefaultWithOptions("MISSING", "def", lookup); got != "def" {
		t.Fatalf("EnvOrDefaultWithOptions = %q", got)
	}
	if got := vsys.EnvIntWithOptions("N", 0, lookup); got != 13 {
		t.Fatalf("EnvIntWithOptions = %d", got)
	}
	if got := vsys.EnvBoolWithOptions("B", false, lookup); !got {
		t.Fatalf("EnvBoolWithOptions = %v", got)
	}
}

func TestFacadeDumpSystemInfo(t *testing.T) {
	var buf bytes.Buffer
	vsys.DumpSystemInfoTo(&buf)
	if buf.Len() == 0 {
		t.Fatal("expected non-empty system info dump")
	}
}

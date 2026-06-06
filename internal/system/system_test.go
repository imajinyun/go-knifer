package system

import (
	"bytes"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"
	"testing"
)

func TestHostInfo(t *testing.T) {
	h := GetHostInfo()
	if h == nil {
		t.Fatal("HostInfo 不应为 nil")
	}
	if h.GetName() == "" {
		t.Errorf("Host Name 不应为空")
	}
	if !strings.Contains(h.String(), "Host Name:") {
		t.Errorf("HostInfo.String 缺少 caption: %s", h.String())
	}
}

func TestHostInfoWithOptions(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("10.0.0.2/24")
	if err != nil {
		t.Fatal(err)
	}
	ipNet.IP = net.ParseIP("10.0.0.2")
	h := NewHostInfoWithOptions(
		WithHostNameFunc(func() (string, error) { return "option-host", nil }),
		WithHostInterfaceAddrsFunc(func() ([]net.Addr, error) { return []net.Addr{ipNet}, nil }),
	)
	if h.GetName() != "option-host" || h.GetAddress() != "10.0.0.2" {
		t.Fatalf("NewHostInfoWithOptions = %#v", h)
	}

	h = GetHostInfoWithOptions(WithHostAddressFunc(func() string { return "192.0.2.10" }))
	if h.GetAddress() != "192.0.2.10" {
		t.Fatalf("GetHostInfoWithOptions address = %q", h.GetAddress())
	}
}

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

func TestUserInfo(t *testing.T) {
	u := GetUserInfo()
	if u == nil {
		t.Fatal("UserInfo 不应为 nil")
	}
	if u.GetCurrentDir() == "" {
		t.Errorf("CurrentDir 不应为空")
	}
	if u.GetTempDir() == "" {
		t.Errorf("TempDir 不应为空")
	}
	if !strings.HasSuffix(u.GetCurrentDir(), string(os.PathSeparator)) {
		t.Errorf("CurrentDir 应以路径分隔符结尾: %q", u.GetCurrentDir())
	}
	if !strings.Contains(u.String(), "User Name:") {
		t.Errorf("UserInfo.String 缺少 caption")
	}
}

func TestUserInfoWithOptions(t *testing.T) {
	u := NewUserInfoWithOptions(
		WithCurrentUserFunc(func() (*user.User, error) {
			return &user.User{Username: "option-user", HomeDir: "/home/option"}, nil
		}),
		WithWorkingDirFunc(func() (string, error) { return "/work/option", nil }),
		WithTempDirFunc(func() string { return "/tmp/option" }),
		WithUserEnvLookup(func(key string) string {
			if key == "LANG" {
				return "zh_CN.UTF-8"
			}
			return ""
		}),
	)
	sep := string(os.PathSeparator)
	if u.GetName() != "option-user" || u.GetHomeDir() != "/home/option"+sep || u.GetCurrentDir() != "/work/option"+sep || u.GetTempDir() != "/tmp/option"+sep {
		t.Fatalf("NewUserInfoWithOptions paths = %#v", u)
	}
	if u.GetLanguage() != "zh" || u.GetCountry() != "CN" {
		t.Fatalf("NewUserInfoWithOptions locale = %s/%s", u.GetLanguage(), u.GetCountry())
	}

	fallback := GetUserInfoWithOptions(
		WithCurrentUserFunc(func() (*user.User, error) { return nil, os.ErrNotExist }),
		WithWorkingDirFunc(func() (string, error) { return "/cwd/fallback", nil }),
		WithTempDirFunc(func() string { return "/tmp/fallback" }),
		WithUserEnvLookup(func(key string) string {
			switch key {
			case "USER":
				return "env-user"
			case "HOME":
				return "/home/env"
			case "LC_ALL":
				return "en_US.UTF-8"
			default:
				return ""
			}
		}),
	)
	if fallback.GetName() != "env-user" || fallback.GetHomeDir() != "/home/env"+sep || fallback.GetLanguage() != "en" || fallback.GetCountry() != "US" {
		t.Fatalf("GetUserInfoWithOptions fallback = %#v", fallback)
	}
}

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

func TestResetInfoCacheClearsSingletons(t *testing.T) {
	ResetInfoCache()
	firstHost := GetHostInfo()
	firstOS := GetOsInfo()
	firstUser := GetUserInfo()
	firstGo := GetGoInfo()
	firstRuntime := GetRuntimeInfo()
	ResetInfoCache()
	if got := GetHostInfo(); got == nil || got == firstHost {
		t.Fatalf("GetHostInfo after reset = %p, first %p", got, firstHost)
	}
	if got := GetOsInfo(); got == nil || got == firstOS {
		t.Fatalf("GetOsInfo after reset = %p, first %p", got, firstOS)
	}
	if got := GetUserInfo(); got == nil || got == firstUser {
		t.Fatalf("GetUserInfo after reset = %p, first %p", got, firstUser)
	}
	if got := GetGoInfo(); got == nil || got == firstGo {
		t.Fatalf("GetGoInfo after reset = %p, first %p", got, firstGo)
	}
	if got := GetRuntimeInfo(); got == nil || got == firstRuntime {
		t.Fatalf("GetRuntimeInfo after reset = %p, first %p", got, firstRuntime)
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

func TestRuntimeInfo(t *testing.T) {
	r := GetRuntimeInfo()
	if r == nil {
		t.Fatal("RuntimeInfo 不应为 nil")
	}
	if r.GetGoroutineCount() <= 0 {
		t.Errorf("Goroutine 数应大于 0")
	}
	if r.GetMaxMemory() == 0 {
		t.Errorf("MaxMemory 不应为 0")
	}
	if !strings.Contains(r.String(), "Goroutine Count:") {
		t.Errorf("RuntimeInfo.String 缺少 caption")
	}
}

func TestRuntimeInfoWithOptions(t *testing.T) {
	readCalls := 0
	r := NewRuntimeInfoWithOptions(
		WithReadMemStatsFunc(func(stats *runtime.MemStats) {
			readCalls++
			stats.Sys = 1024
			stats.HeapSys = 512
			stats.HeapIdle = 128
			stats.HeapInuse = 256
		}),
		WithNumGoroutineFunc(func() int { return 7 }),
	)
	if readCalls != 1 || r.GetMaxMemory() != 1024 || r.GetTotalMemory() != 512 || r.GetFreeMemory() != 128 || r.GetUsableMemory() != 768 || r.GetGoroutineCount() != 7 {
		t.Fatalf("NewRuntimeInfoWithOptions = %#v calls=%d", r, readCalls)
	}
	r.Refresh()
	if readCalls != 2 {
		t.Fatalf("Refresh read calls = %d", readCalls)
	}

	r = GetRuntimeInfoWithOptions(WithReadMemStatsFunc(func(stats *runtime.MemStats) { stats.Sys = 2048 }))
	if r.GetMaxMemory() != 2048 {
		t.Fatalf("GetRuntimeInfoWithOptions max = %d", r.GetMaxMemory())
	}
}

func TestGetCurrentPID(t *testing.T) {
	if GetCurrentPID() != os.Getpid() {
		t.Errorf("PID 不一致")
	}
	if got := GetCurrentPIDWithOptions(WithPIDFunc(func() int { return 4242 })); got != 4242 {
		t.Fatalf("GetCurrentPIDWithOptions = %d", got)
	}
}

func TestGetEnv(t *testing.T) {
	t.Setenv("GKSYSTEM_TEST_KEY", "abc")
	if v := Get("GKSYSTEM_TEST_KEY", true); v != "abc" {
		t.Errorf("Get 应返回 abc，实际 %q", v)
	}
	if v := GetOrDefault("GKSYSTEM_TEST_NOT_EXIST", "def"); v != "def" {
		t.Errorf("GetOrDefault 默认值未生效: %q", v)
	}

	t.Setenv("GKSYSTEM_TEST_INT", "42")
	if n := GetInt("GKSYSTEM_TEST_INT", 0); n != 42 {
		t.Errorf("GetInt: 期望 42，实际 %d", n)
	}
	if n := GetInt("GKSYSTEM_TEST_INT_INVALID", 7); n != 7 {
		t.Errorf("GetInt 无效值应返回默认: %d", n)
	}

	t.Setenv("GKSYSTEM_TEST_BOOL", "true")
	if b := GetBool("GKSYSTEM_TEST_BOOL", false); !b {
		t.Errorf("GetBool 应为 true")
	}
}

func TestGetEnvWithOptions(t *testing.T) {
	lookup := func(key string) (string, bool) {
		switch key {
		case "STRING":
			return "value", true
		case "INT":
			return "12", true
		case "BOOL":
			return "true", true
		case "EMPTY":
			return "", true
		default:
			return "", false
		}
	}
	var warning bytes.Buffer
	if got := GetWithOptions("STRING", true, WithEnvLookupFunc(lookup)); got != "value" {
		t.Fatalf("GetWithOptions = %q", got)
	}
	if got := GetWithOptions("MISSING", false, WithEnvLookupFunc(lookup), WithEnvWarningWriter(&warning)); got != "" || !strings.Contains(warning.String(), "MISSING") {
		t.Fatalf("GetWithOptions missing = %q warning=%q", got, warning.String())
	}
	if got := GetOrDefaultWithOptions("EMPTY", "def", WithEnvLookupFunc(lookup)); got != "def" {
		t.Fatalf("GetOrDefaultWithOptions empty = %q", got)
	}
	if got := GetIntWithOptions("INT", 0, WithEnvLookupFunc(lookup)); got != 12 {
		t.Fatalf("GetIntWithOptions = %d", got)
	}
	if got := GetBoolWithOptions("BOOL", false, WithEnvLookupFunc(lookup)); !got {
		t.Fatalf("GetBoolWithOptions = %v", got)
	}
}

func TestTotalThreadCount(t *testing.T) {
	if GetTotalThreadCount() <= 0 {
		t.Errorf("总协程数应大于 0")
	}
	if got := GetTotalThreadCountWithOptions(WithProcessNumGoroutineFunc(func() int { return 6 })); got != 6 {
		t.Fatalf("GetTotalThreadCountWithOptions = %d", got)
	}
}

func TestDumpSystemInfo(t *testing.T) {
	var buf bytes.Buffer
	DumpSystemInfoTo(&buf)
	out := buf.String()
	for _, kw := range []string{"Go Version:", "OS Name:", "User Name:", "Host Name:", "Goroutine Count:"} {
		if !strings.Contains(out, kw) {
			t.Errorf("Dump 输出缺少 %q：\n%s", kw, out)
		}
	}
}

func TestReadableSize(t *testing.T) {
	cases := []struct {
		in   uint64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
	}
	for _, c := range cases {
		got := readableSize(c.in)
		if got != c.want {
			t.Errorf("readableSize(%d): 期望 %q 实际 %q", c.in, c.want, got)
		}
	}
}

func TestParseLocale(t *testing.T) {
	lang, country := parseLocale("zh_CN.UTF-8")
	if lang != "zh" || country != "CN" {
		t.Errorf("parseLocale(zh_CN.UTF-8) 错误: %s/%s", lang, country)
	}
	lang, country = parseLocale("")
	if lang != "" || country != "" {
		t.Errorf("空 locale 应返回空")
	}
	lang, country = parseLocale("en")
	if lang != "en" || country != "" {
		t.Errorf("parseLocale(en) 错误: %s/%s", lang, country)
	}
}

func TestFixPath(t *testing.T) {
	if fixPath("") != "" {
		t.Errorf("空字符串应保持空")
	}
	sep := string(os.PathSeparator)
	if fixPath("/tmp"+sep) != "/tmp"+sep {
		t.Errorf("已带后缀不应再追加")
	}
	if fixPath("/tmp") != "/tmp"+sep {
		t.Errorf("应追加分隔符")
	}
}

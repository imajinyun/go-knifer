package system

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// Singleton caches avoid repeated collection.
var (
	infoMu      sync.Mutex
	hostOnce    sync.Once
	hostInfo    *HostInfo
	osOnce      sync.Once
	osInfo      *OsInfo
	userOnce    sync.Once
	userInfo    *UserInfo
	goOnce      sync.Once
	goInfo      *GoInfo
	runtimeOnce sync.Once
	runtimeRef  *RuntimeInfo
)

// ResetInfoCache clears cached singleton system information.
func ResetInfoCache() {
	infoMu.Lock()
	defer infoMu.Unlock()
	hostOnce = sync.Once{}
	hostInfo = nil
	osOnce = sync.Once{}
	osInfo = nil
	userOnce = sync.Once{}
	userInfo = nil
	goOnce = sync.Once{}
	goInfo = nil
	runtimeOnce = sync.Once{}
	runtimeRef = nil
}

type processConfig struct {
	pid          func() int
	numGoroutine func() int
}

// ProcessOption customizes process/runtime scalar helpers per call.
type ProcessOption func(*processConfig)

// WithPIDFunc sets the function used to collect the current process id.
func WithPIDFunc(fn func() int) ProcessOption {
	return func(c *processConfig) {
		if fn != nil {
			c.pid = fn
		}
	}
}

// WithProcessNumGoroutineFunc sets the function used by process scalar helpers to collect goroutine count.
func WithProcessNumGoroutineFunc(fn func() int) ProcessOption {
	return func(c *processConfig) {
		if fn != nil {
			c.numGoroutine = fn
		}
	}
}

func applyProcessOptions(opts []ProcessOption) processConfig {
	cfg := processConfig{pid: os.Getpid, numGoroutine: runtime.NumGoroutine}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.pid == nil {
		cfg.pid = os.Getpid
	}
	if cfg.numGoroutine == nil {
		cfg.numGoroutine = runtime.NumGoroutine
	}
	return cfg
}

type envConfig struct {
	lookup        func(string) (string, bool)
	parseInt      func(string) (int, error)
	parseBool     func(string) (bool, error)
	warningWriter io.Writer
}

// DumpOption customizes system information dumping per call.
type DumpOption func(*dumpConfig)

type dumpConfig struct {
	hostOpts    []HostInfoOption
	osOpts      []OsInfoOption
	userOpts    []UserInfoOption
	goOpts      []GoInfoOption
	runtimeOpts []RuntimeInfoOption
}

// WithDumpHostOptions sets host information providers used by DumpSystemInfoWithOptions.
func WithDumpHostOptions(opts ...HostInfoOption) DumpOption {
	return func(c *dumpConfig) { c.hostOpts = append([]HostInfoOption(nil), opts...) }
}

// WithDumpOsOptions sets OS information providers used by DumpSystemInfoWithOptions.
func WithDumpOsOptions(opts ...OsInfoOption) DumpOption {
	return func(c *dumpConfig) { c.osOpts = append([]OsInfoOption(nil), opts...) }
}

// WithDumpUserOptions sets user information providers used by DumpSystemInfoWithOptions.
func WithDumpUserOptions(opts ...UserInfoOption) DumpOption {
	return func(c *dumpConfig) { c.userOpts = append([]UserInfoOption(nil), opts...) }
}

// WithDumpGoOptions sets Go runtime metadata providers used by DumpSystemInfoWithOptions.
func WithDumpGoOptions(opts ...GoInfoOption) DumpOption {
	return func(c *dumpConfig) { c.goOpts = append([]GoInfoOption(nil), opts...) }
}

// WithDumpRuntimeOptions sets runtime information providers used by DumpSystemInfoWithOptions.
func WithDumpRuntimeOptions(opts ...RuntimeInfoOption) DumpOption {
	return func(c *dumpConfig) { c.runtimeOpts = append([]RuntimeInfoOption(nil), opts...) }
}

func applyDumpOptions(opts []DumpOption) dumpConfig {
	cfg := dumpConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// EnvOption customizes environment helpers per call.
type EnvOption func(*envConfig)

// WithEnvLookupFunc sets the function used to look up environment variables.
func WithEnvLookupFunc(fn func(string) (string, bool)) EnvOption {
	return func(c *envConfig) {
		if fn != nil {
			c.lookup = fn
		}
	}
}

// WithEnvWarningWriter sets the writer used for missing-variable warnings.
func WithEnvWarningWriter(w io.Writer) EnvOption {
	return func(c *envConfig) {
		if w != nil {
			c.warningWriter = w
		}
	}
}

// WithEnvIntParser sets the parser used by GetIntWithOptions.
func WithEnvIntParser(parser func(string) (int, error)) EnvOption {
	return func(c *envConfig) {
		if parser != nil {
			c.parseInt = parser
		}
	}
}

// WithEnvBoolParser sets the parser used by GetBoolWithOptions.
func WithEnvBoolParser(parser func(string) (bool, error)) EnvOption {
	return func(c *envConfig) {
		if parser != nil {
			c.parseBool = parser
		}
	}
}

func applyEnvOptions(opts []EnvOption) envConfig {
	cfg := envConfig{lookup: os.LookupEnv, parseInt: strconv.Atoi, parseBool: strconv.ParseBool, warningWriter: os.Stderr}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.lookup == nil {
		cfg.lookup = os.LookupEnv
	}
	if cfg.parseInt == nil {
		cfg.parseInt = strconv.Atoi
	}
	if cfg.parseBool == nil {
		cfg.parseBool = strconv.ParseBool
	}
	if cfg.warningWriter == nil {
		cfg.warningWriter = os.Stderr
	}
	return cfg
}

// GetHostInfo returns cached host information.
func GetHostInfo() *HostInfo {
	infoMu.Lock()
	defer infoMu.Unlock()
	hostOnce.Do(func() { hostInfo = NewHostInfo() })
	return hostInfo
}

// GetHostInfoWithOptions returns uncached host information collected with per-call options.
func GetHostInfoWithOptions(opts ...HostInfoOption) *HostInfo {
	return NewHostInfoWithOptions(opts...)
}

// GetOsInfo returns cached OS information.
func GetOsInfo() *OsInfo {
	infoMu.Lock()
	defer infoMu.Unlock()
	osOnce.Do(func() { osInfo = NewOsInfo() })
	return osInfo
}

// GetOsInfoWithOptions returns uncached OS information collected with per-call options.
func GetOsInfoWithOptions(opts ...OsInfoOption) *OsInfo {
	return NewOsInfoWithOptions(opts...)
}

// GetUserInfo returns cached user information.
func GetUserInfo() *UserInfo {
	infoMu.Lock()
	defer infoMu.Unlock()
	userOnce.Do(func() { userInfo = NewUserInfo() })
	return userInfo
}

// GetUserInfoWithOptions returns uncached user information collected with per-call options.
func GetUserInfoWithOptions(opts ...UserInfoOption) *UserInfo {
	return NewUserInfoWithOptions(opts...)
}

// GetGoInfo returns cached Go runtime metadata.
func GetGoInfo() *GoInfo {
	infoMu.Lock()
	defer infoMu.Unlock()
	goOnce.Do(func() { goInfo = NewGoInfo() })
	return goInfo
}

// GetGoInfoWithOptions returns uncached Go runtime metadata collected with per-call options.
func GetGoInfoWithOptions(opts ...GoInfoOption) *GoInfo {
	return NewGoInfoWithOptions(opts...)
}

// GetRuntimeInfo returns runtime memory information and refreshes it on each call.
func GetRuntimeInfo() *RuntimeInfo {
	infoMu.Lock()
	defer infoMu.Unlock()
	runtimeOnce.Do(func() { runtimeRef = NewRuntimeInfo() })
	return runtimeRef.Refresh()
}

// GetRuntimeInfoWithOptions returns uncached runtime information collected with per-call options.
func GetRuntimeInfoWithOptions(opts ...RuntimeInfoOption) *RuntimeInfo {
	return NewRuntimeInfoWithOptions(opts...)
}

// GetCurrentPID returns the current process PID.
func GetCurrentPID() int {
	return GetCurrentPIDWithOptions()
}

// GetCurrentPIDWithOptions returns the current process PID using custom providers.
func GetCurrentPIDWithOptions(opts ...ProcessOption) int {
	return applyProcessOptions(opts).pid()
}

// GetTotalMemory returns total memory requested from the OS by the current Go program.
func GetTotalMemory() uint64 {
	return GetRuntimeInfo().GetTotalMemory()
}

// GetTotalMemoryWithOptions returns total memory using custom runtime providers.
func GetTotalMemoryWithOptions(opts ...RuntimeInfoOption) uint64 {
	return GetRuntimeInfoWithOptions(opts...).GetTotalMemory()
}

// GetFreeMemory returns idle memory in the current Go program.
func GetFreeMemory() uint64 {
	return GetRuntimeInfo().GetFreeMemory()
}

// GetFreeMemoryWithOptions returns idle memory using custom runtime providers.
func GetFreeMemoryWithOptions(opts ...RuntimeInfoOption) uint64 {
	return GetRuntimeInfoWithOptions(opts...).GetFreeMemory()
}

// GetMaxMemory returns the detected memory upper bound for the current Go program.
func GetMaxMemory() uint64 {
	return GetRuntimeInfo().GetMaxMemory()
}

// GetMaxMemoryWithOptions returns the memory upper bound using custom runtime providers.
func GetMaxMemoryWithOptions(opts ...RuntimeInfoOption) uint64 {
	return GetRuntimeInfoWithOptions(opts...).GetMaxMemory()
}

// GetTotalThreadCount returns the total goroutine count.
func GetTotalThreadCount() int {
	return GetTotalThreadCountWithOptions()
}

// GetTotalThreadCountWithOptions returns the goroutine count using custom providers.
func GetTotalThreadCountWithOptions(opts ...ProcessOption) int {
	return applyProcessOptions(opts).numGoroutine()
}

// Get returns an environment variable by key.
// If quiet is false and the variable is missing, it prints a warning to stderr.
func Get(key string, quiet bool) string {
	return GetWithOptions(key, quiet)
}

// GetWithOptions returns an environment variable by key using custom providers.
// If quiet is false and the variable is missing, it prints a warning to the configured writer.
func GetWithOptions(key string, quiet bool, opts ...EnvOption) string {
	cfg := applyEnvOptions(opts)
	v, ok := cfg.lookup(key)
	if !ok && !quiet {
		_, _ = fmt.Fprintf(cfg.warningWriter, "[gksystem] env %q not found\n", key)
	}
	return v
}

// GetOrDefault returns an environment variable, or def when it is missing or empty.
func GetOrDefault(key, def string) string {
	return GetOrDefaultWithOptions(key, def)
}

// GetOrDefaultWithOptions returns an environment variable, or def when it is missing or empty, using custom providers.
func GetOrDefaultWithOptions(key, def string, opts ...EnvOption) string {
	v, ok := applyEnvOptions(opts).lookup(key)
	if !ok || v == "" {
		return def
	}
	return v
}

// GetInt returns an environment variable as an int, or def on conversion failure.
func GetInt(key string, def int) int {
	return GetIntWithOptions(key, def)
}

// GetIntWithOptions returns an environment variable as an int, or def on conversion failure, using custom providers.
func GetIntWithOptions(key string, def int, opts ...EnvOption) int {
	cfg := applyEnvOptions(opts)
	v, ok := cfg.lookup(key)
	if !ok {
		return def
	}
	n, err := cfg.parseInt(v)
	if err != nil {
		return def
	}
	return n
}

// GetBool returns an environment variable as a bool, or def on conversion failure.
func GetBool(key string, def bool) bool {
	return GetBoolWithOptions(key, def)
}

// GetBoolWithOptions returns an environment variable as a bool, or def on conversion failure, using custom providers.
func GetBoolWithOptions(key string, def bool, opts ...EnvOption) bool {
	cfg := applyEnvOptions(opts)
	v, ok := cfg.lookup(key)
	if !ok {
		return def
	}
	b, err := cfg.parseBool(v)
	if err != nil {
		return def
	}
	return b
}

// DumpSystemInfo writes system information to stdout.
func DumpSystemInfo() {
	DumpSystemInfoTo(os.Stdout)
}

// DumpSystemInfoTo writes system information to the specified writer.
func DumpSystemInfoTo(w io.Writer) {
	const sep = "--------------\n"
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetGoInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetOsInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetUserInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetHostInfo())
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetRuntimeInfo())
	_, _ = fmt.Fprint(w, sep)
}

// DumpSystemInfoWithOptions writes uncached system information to w using per-call providers.
func DumpSystemInfoWithOptions(w io.Writer, opts ...DumpOption) {
	if w == nil {
		w = io.Discard
	}
	cfg := applyDumpOptions(opts)
	const sep = "--------------\n"
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetGoInfoWithOptions(cfg.goOpts...))
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetOsInfoWithOptions(cfg.osOpts...))
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetUserInfoWithOptions(cfg.userOpts...))
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetHostInfoWithOptions(cfg.hostOpts...))
	_, _ = fmt.Fprint(w, sep)
	_, _ = fmt.Fprint(w, GetRuntimeInfoWithOptions(cfg.runtimeOpts...))
	_, _ = fmt.Fprint(w, sep)
}

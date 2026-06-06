package log

import "sync"

// LogFactory 对应 the utility toolkit LogFactory，提供根据名称获取 Log 实例的能力。
//
// 通过 SetFactory 可全局替换实现，默认返回 ConsoleLog。
type LogFactory interface {
	// CreateLog 根据名称创建 Log 实例。
	CreateLog(name string) Log
}

// LogFactoryFunc 适配函数为 LogFactory。
type LogFactoryFunc func(name string) Log

// CreateLog 调用底层函数。
func (f LogFactoryFunc) CreateLog(name string) Log { return f(name) }

// LoggerOption customizes logger lookup/creation for one call.
type LoggerOption func(*loggerConfig)

type loggerConfig struct {
	factory LogFactory
	cache   bool
}

// WithLoggerFactory sets the logger factory used by GetWithOptions or NewIsolatedLogger.
func WithLoggerFactory(factory LogFactory) LoggerOption {
	return func(cfg *loggerConfig) {
		if factory != nil {
			cfg.factory = factory
			cfg.cache = false
		}
	}
}

// WithLoggerConsoleOptions builds loggers with console options for one lookup/creation call.
func WithLoggerConsoleOptions(opts ...ConsoleLogOption) LoggerOption {
	return WithLoggerFactory(LogFactoryFunc(func(name string) Log {
		return NewConsoleLogWithOptions(name, opts...)
	}))
}

// WithLoggerCache controls whether GetWithOptions may use the package-level logger cache.
func WithLoggerCache(enabled bool) LoggerOption {
	return func(cfg *loggerConfig) {
		cfg.cache = enabled
	}
}

func defaultLogFactory() LogFactory {
	return LogFactoryFunc(func(name string) Log { return NewConsoleLog(name) })
}

func applyLoggerOptions(base loggerConfig, opts ...LoggerOption) loggerConfig {
	if base.factory == nil {
		base.factory = defaultLogFactory()
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&base)
		}
	}
	if base.factory == nil {
		base.factory = defaultLogFactory()
	}
	return base
}

var (
	factoryMu      sync.RWMutex
	currentFactory LogFactory = defaultLogFactory()

	logCache   = make(map[string]Log)
	logCacheMu sync.RWMutex
)

// SetFactory 设置全局日志工厂。设置后会清空已缓存的 Log 实例。
func SetFactory(factory LogFactory) {
	if factory == nil {
		return
	}
	factoryMu.Lock()
	currentFactory = factory
	factoryMu.Unlock()

	logCacheMu.Lock()
	logCache = make(map[string]Log)
	logCacheMu.Unlock()
}

// GetFactory 返回当前的日志工厂。
func GetFactory() LogFactory {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	return currentFactory
}

// Get 根据名称获取一个 Log 实例（带缓存）。
func Get(name string) Log {
	logCacheMu.RLock()
	if l, ok := logCache[name]; ok {
		logCacheMu.RUnlock()
		return l
	}
	logCacheMu.RUnlock()

	logCacheMu.Lock()
	defer logCacheMu.Unlock()
	// 双重检查
	if l, ok := logCache[name]; ok {
		return l
	}
	l := GetFactory().CreateLog(name)
	logCache[name] = l
	return l
}

// GetWithOptions returns a logger by name with per-call factory/cache options.
func GetWithOptions(name string, opts ...LoggerOption) Log {
	if len(opts) == 0 {
		return Get(name)
	}
	cfg := applyLoggerOptions(loggerConfig{factory: GetFactory(), cache: true}, opts...)
	if cfg.cache {
		return Get(name)
	}
	return cfg.factory.CreateLog(name)
}

// NewIsolatedLogger creates a logger without reading package-level factory/cache state.
func NewIsolatedLogger(name string, opts ...LoggerOption) Log {
	cfg := applyLoggerOptions(loggerConfig{factory: defaultLogFactory(), cache: false}, opts...)
	return cfg.factory.CreateLog(name)
}

// GetDefault 返回名称为 "default" 的 Log 实例。
func GetDefault() Log {
	return Get("default")
}

// GetDefaultWithOptions returns the default logger with per-call factory/cache options.
func GetDefaultWithOptions(opts ...LoggerOption) Log {
	return GetWithOptions("default", opts...)
}

package log

import "sync"

// LogFactory 对应 hutool LogFactory，提供根据名称获取 Log 实例的能力。
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

var (
	factoryMu      sync.RWMutex
	currentFactory LogFactory = LogFactoryFunc(func(name string) Log { return NewConsoleLog(name) })

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

// GetDefault 返回名称为 "default" 的 Log 实例。
func GetDefault() Log {
	return Get("default")
}

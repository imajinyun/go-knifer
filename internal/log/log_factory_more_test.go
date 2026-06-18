package log

import (
	"testing"
)

func TestWithLoggerCache(t *testing.T) {
	// WithLoggerCache(false) should bypass the cache and create a new instance.
	factory := LogFactoryFunc(func(name string) Log {
		return NewConsoleLog(name)
	})
	cfg := applyLoggerOptions(loggerConfig{factory: factory, cache: true}, WithLoggerCache(false))
	if cfg.cache {
		t.Error("WithLoggerCache(false) should set cache to false")
	}

	// WithLoggerCache(true) should keep the cache enabled.
	cfg2 := applyLoggerOptions(loggerConfig{factory: factory, cache: false}, WithLoggerCache(true))
	if !cfg2.cache {
		t.Error("WithLoggerCache(true) should set cache to true")
	}
}

func TestGetDefault(t *testing.T) {
	l := GetDefault()
	if l == nil {
		t.Fatal("GetDefault should return a non-nil Log")
	}
	if l.GetName() != "default" {
		t.Errorf("GetDefault name = %q, want %q", l.GetName(), "default")
	}
}

func TestGetDefaultWithOptions(t *testing.T) {
	factory := LogFactoryFunc(func(name string) Log {
		return NewConsoleLog("custom:" + name)
	})
	l := GetDefaultWithOptions(WithLoggerFactory(factory))
	if l == nil {
		t.Fatal("GetDefaultWithOptions should return a non-nil Log")
	}
	if l.GetName() != "custom:default" {
		t.Errorf("GetDefaultWithOptions name = %q, want %q", l.GetName(), "custom:default")
	}
}

package vlog_test

import (
	"testing"

	"github.com/imajinyun/go-knifer/vlog"
)

func TestFacadeLogger(t *testing.T) {
	log := vlog.NewConsoleLog("test")
	if log == nil {
		t.Fatal("expected non-nil logger")
	}

	// smoke test: log at each level should not panic
	log.Trace("trace")
	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
}

func TestFacadeDefaultLogger(t *testing.T) {
	log := vlog.DefaultLogger()
	if log == nil {
		t.Fatal("expected non-nil default logger")
	}
	log.Info("default logger works")
}

func TestFacadeLoggerByName(t *testing.T) {
	log1 := vlog.Logger("foo")
	log2 := vlog.Logger("foo")
	if log1 == nil || log2 == nil {
		t.Fatal("expected non-nil loggers")
	}
}

func TestFacadeLogLevel(t *testing.T) {
	old := vlog.GetLogLevel()
	vlog.SetLogLevel(vlog.LogLevelDebug)
	if vlog.GetLogLevel() != vlog.LogLevelDebug {
		t.Fatal("expected log level to be set to Debug")
	}
	vlog.SetLogLevel(old)
}

func TestFacadeStaticLog(t *testing.T) {
	// smoke test: static log functions should not panic
	vlog.Trace("static trace")
	vlog.Debug("static debug")
	vlog.Info("static info")
	vlog.Warn("static warn")
	vlog.ErrorLog("static error")
	vlog.Tracef("formatted %s", "trace")
	vlog.Debugf("formatted %s", "debug")
	vlog.Infof("formatted %s", "info")
	vlog.Warnf("formatted %s", "warn")
	vlog.Errorf("formatted %s", "error")
}

package vlog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

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

func TestFacadeLoggerWithOptions(t *testing.T) {
	log := vlog.LoggerWithOptions("facade.logger.option", vlog.WithLoggerFactory(vlog.LogFactoryFunc(func(name string) vlog.Log {
		return vlog.NewConsoleLog("facade:" + name)
	})))
	if log.GetName() != "facade:facade.logger.option" {
		t.Fatalf("LoggerWithOptions name = %q", log.GetName())
	}

	vlog.SetLogFactory(vlog.LogFactoryFunc(func(name string) vlog.Log { return vlog.NewConsoleLog("global:" + name) }))
	defer vlog.SetLogFactory(vlog.LogFactoryFunc(func(name string) vlog.Log { return vlog.NewConsoleLog(name) }))
	isolated := vlog.NewIsolatedLogger("facade.isolated")
	if isolated.GetName() != "facade.isolated" {
		t.Fatalf("NewIsolatedLogger leaked global factory: %q", isolated.GetName())
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

func TestFacadeConsoleLogOptions(t *testing.T) {
	old := vlog.GetLogLevel()
	vlog.SetLogLevel(vlog.LogLevelDebug)
	defer vlog.SetLogLevel(old)

	out := &bytes.Buffer{}
	fixed := time.Date(2024, 4, 5, 6, 7, 8, 0, time.UTC)
	log := vlog.NewConsoleLogWithOptions("facade.options",
		vlog.WithLogClock(func() time.Time { return fixed }),
		vlog.WithLogTimeLayout(time.RFC3339),
		vlog.WithLogOutput(out, &bytes.Buffer{}),
	)
	log.Info("hello")
	if !strings.Contains(out.String(), "2024-04-05T06:07:08Z") || !strings.Contains(out.String(), "hello") {
		t.Fatalf("console log options not applied: %q", out.String())
	}

	colorOut := &bytes.Buffer{}
	colorLog := vlog.NewConsoleColorLogWithOptions("facade.color",
		vlog.WithLogClock(func() time.Time { return fixed }),
		vlog.WithLogTimeLayout("15:04"),
		vlog.WithLogOutput(colorOut, &bytes.Buffer{}),
	)
	colorLog.Info("color")
	if !strings.Contains(colorOut.String(), "06:07") || !strings.Contains(colorOut.String(), "color") {
		t.Fatalf("color log options not applied: %q", colorOut.String())
	}

	customColorOut := &bytes.Buffer{}
	customColorLog := vlog.NewConsoleColorLogWithOptions("facade.color.custom",
		vlog.WithLogOutput(customColorOut, &bytes.Buffer{}),
		vlog.WithLogColorFactory(func(vlog.Level) string { return "\033[36m" }),
	)
	customColorLog.Info("custom-color")
	if !strings.Contains(customColorOut.String(), "\033[36m") || !strings.Contains(customColorOut.String(), "custom-color") {
		t.Fatalf("color factory option not applied: %q", customColorOut.String())
	}
}

func TestFacadeStaticLogWithOptions(t *testing.T) {
	old := vlog.GetLogLevel()
	vlog.SetLogLevel(vlog.LogLevelInfo)
	defer vlog.SetLogLevel(old)

	out := &bytes.Buffer{}
	fixed := time.Date(2024, 7, 8, 9, 10, 11, 0, time.UTC)
	vlog.InfoWithOptions([]vlog.LoggerOption{vlog.WithLoggerConsoleOptions(
		vlog.WithLogClock(func() time.Time { return fixed }),
		vlog.WithLogTimeLayout(time.RFC3339),
		vlog.WithLogOutput(out, &bytes.Buffer{}),
	)}, "facade-static")
	if !strings.Contains(out.String(), "2024-07-08T09:10:11Z") || !strings.Contains(out.String(), "facade-static") {
		t.Fatalf("static log options not applied: %q", out.String())
	}
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

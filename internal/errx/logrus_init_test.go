package errx

import (
	"io"
	"testing"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

func TestInitWithOptionsUsesInjectedProviders(t *testing.T) {
	var reportCaller bool
	var output io.Writer
	var formatter logrus.Formatter
	var getenvKey string
	var setDSN string
	var hookClient *raven.Client
	var hookLevels []logrus.Level
	var hookAdded bool

	InitWithOptions(
		WithSentryEnvKey("CUSTOM_DSN"),
		WithEnvLookupFunc(func(key string) string {
			getenvKey = key
			return "https://public@example.com/1"
		}),
		WithLogrusConfigurer(
			func(v bool) { reportCaller = v },
			func(w io.Writer) { output = w },
			func(f logrus.Formatter) { formatter = f },
		),
		WithRavenSetDSNFunc(func(dsn string) error {
			setDSN = dsn
			return nil
		}),
		WithSentryHookFactory(func(client *raven.Client, levels []logrus.Level) (*logrus_sentry.SentryHook, error) {
			hookClient = client
			hookLevels = append([]logrus.Level(nil), levels...)
			return &logrus_sentry.SentryHook{}, nil
		}),
		WithLogHookAdder(func(logrus.Hook) { hookAdded = true }),
		WithSentryLevels(logrus.ErrorLevel),
	)

	if !reportCaller || output != io.Discard || formatter != EmptyFormatter {
		t.Fatalf("logrus config = reportCaller %v output %T formatter %T", reportCaller, output, formatter)
	}
	if getenvKey != "CUSTOM_DSN" || setDSN != "https://public@example.com/1" {
		t.Fatalf("dsn providers key=%q dsn=%q", getenvKey, setDSN)
	}
	if hookClient == nil || len(hookLevels) != 1 || hookLevels[0] != logrus.ErrorLevel || !hookAdded {
		t.Fatalf("hook providers client=%v levels=%v added=%v", hookClient, hookLevels, hookAdded)
	}
}

func TestNewIsolatedLogrusWithOptionsDoesNotUseGlobalConfigurers(t *testing.T) {
	var globalConfigured bool
	logger := NewIsolatedLogrusWithOptions(
		WithEnvLookupFunc(func(string) string { return "" }),
		WithLogOutput(io.Discard),
		WithReportCaller(false),
		WithLogrusConfigurer(
			func(bool) { globalConfigured = true },
			func(io.Writer) { globalConfigured = true },
			func(logrus.Formatter) { globalConfigured = true },
		),
	)
	if logger == nil {
		t.Fatal("NewIsolatedLogrusWithOptions returned nil")
	}
	if globalConfigured {
		t.Fatal("isolated logger should not call global logrus configurers")
	}
	if logger.Out != io.Discard || logger.ReportCaller {
		t.Fatalf("isolated logger config out=%T reportCaller=%v", logger.Out, logger.ReportCaller)
	}
}

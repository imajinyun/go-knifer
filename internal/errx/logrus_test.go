package errx

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

func TestWrapperExecReturnsFunctionError(t *testing.T) {
	silenceLogrus(t)

	want := errors.New("wrapped failure")
	got := Wrap(func() error { return want }).WithErrorf("failed").Exec(context.Background())
	if !ErrorIs(got, want) {
		t.Fatalf("Exec() error = %v, want %v", got, want)
	}
}

func TestWrapperExecConvertsPanic(t *testing.T) {
	silenceLogrus(t)

	got := Wrap(func() error {
		panic("panic from wrapper")
	}).WithWarnf("panic").Exec(context.TODO())
	if got == nil || !strings.Contains(got.Error(), "panic from wrapper") {
		t.Fatalf("Exec() panic error = %v, want panic value", got)
	}
}

func TestWrapperExecNilFunction(t *testing.T) {
	if err := Wrap(nil).Exec(context.Background()); err != nil {
		t.Fatalf("Exec(nil) = %v, want nil", err)
	}
}

func TestRecoverHelpers(t *testing.T) {
	silenceLogrus(t)

	want := errors.New("recover failure")
	if got := Recover(func() error { return want }, "recover"); !ErrorIs(got, want) {
		t.Fatalf("Recover() = %v, want %v", got, want)
	}
	got := RecoverWithoutError(func() { panic("recover without error") }, "recover without error")
	if got == nil || !strings.Contains(got.Error(), "recover without error") {
		t.Fatalf("RecoverWithoutError() = %v, want panic value", got)
	}
	if got := RecoverWithoutError(nil, "nil function"); got != nil {
		t.Fatalf("RecoverWithoutError(nil) = %v, want nil", got)
	}
}

func TestEmptyFormatterSuppressesOutput(t *testing.T) {
	data, err := EmptyFormatter.Format(logrus.NewEntry(logrus.New()))
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Fatalf("EmptyFormatter output length = %d, want 0", len(data))
	}
}

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

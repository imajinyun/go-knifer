package errx

import (
	"io"
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

const (
	// SentryDSN is the environment variable used to override the configured DSN.
	SentryDSN = "SENTRY_DSN"
)

type initConfig struct {
	dsn          string
	envKey       string
	output       io.Writer
	formatter    logrus.Formatter
	reportCaller bool
	levels       []logrus.Level
}

// InitOption customizes logrus/Sentry initialization.
type InitOption func(*initConfig)

// WithSentryDSN sets the Sentry DSN.
func WithSentryDSN(dsn string) InitOption { return func(c *initConfig) { c.dsn = dsn } }

// WithSentryEnvKey sets the environment variable used to override the Sentry DSN.
func WithSentryEnvKey(key string) InitOption { return func(c *initConfig) { c.envKey = key } }

// WithLogOutput sets the logrus output writer.
func WithLogOutput(output io.Writer) InitOption { return func(c *initConfig) { c.output = output } }

// WithLogFormatter sets the logrus formatter.
func WithLogFormatter(formatter logrus.Formatter) InitOption {
	return func(c *initConfig) { c.formatter = formatter }
}

// WithReportCaller controls whether logrus records caller information.
func WithReportCaller(reportCaller bool) InitOption {
	return func(c *initConfig) { c.reportCaller = reportCaller }
}

// WithSentryLevels sets the log levels forwarded to Sentry.
func WithSentryLevels(levels ...logrus.Level) InitOption {
	return func(c *initConfig) { c.levels = append([]logrus.Level(nil), levels...) }
}

func applyInitOptions(dsn string, opts []InitOption) initConfig {
	cfg := initConfig{
		dsn:          dsn,
		envKey:       SentryDSN,
		output:       io.Discard,
		formatter:    EmptyFormatter,
		reportCaller: true,
		levels:       []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.output == nil {
		cfg.output = io.Discard
	}
	if cfg.formatter == nil {
		cfg.formatter = EmptyFormatter
	}
	if len(cfg.levels) == 0 {
		cfg.levels = []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}
	}
	return cfg
}

// Init configures logrus to forward logs to the internal logs hook and,
// when a DSN is provided, to Sentry as well.
func Init(sentryDSN string) {
	InitWithOptions(WithSentryDSN(sentryDSN))
}

// InitWithOptions configures logrus output and optional Sentry forwarding with custom options.
func InitWithOptions(opts ...InitOption) {
	cfg := applyInitOptions("", opts)
	logrus.SetReportCaller(cfg.reportCaller)
	logrus.SetOutput(cfg.output)
	logrus.SetFormatter(cfg.formatter)

	if dsn := os.Getenv(cfg.envKey); dsn != "" {
		cfg.dsn = dsn
	}
	if cfg.dsn == "" {
		return
	}
	if err := raven.SetDSN(cfg.dsn); err != nil {
		logrus.WithError(err).Error("raven init failed")
		return
	}

	sentry, err := logrus_sentry.NewAsyncWithClientSentryHook(
		raven.DefaultClient,
		cfg.levels,
	)
	if err != nil {
		logrus.WithError(err).Error("sentry hook init failed")
		return
	}
	sentry.StacktraceConfiguration.Enable = true
	sentry.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	logrus.AddHook(sentry)
}

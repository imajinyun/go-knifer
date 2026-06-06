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
	dsn             string
	envKey          string
	output          io.Writer
	formatter       logrus.Formatter
	reportCaller    bool
	levels          []logrus.Level
	getenv          func(string) string
	setDSN          func(string) error
	sentryClient    *raven.Client
	newSentryHook   func(*raven.Client, []logrus.Level) (*logrus_sentry.SentryHook, error)
	addHook         func(logrus.Hook)
	setReportCaller func(bool)
	setOutput       func(io.Writer)
	setFormatter    func(logrus.Formatter)
	logError        func(error, string)
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

// WithEnvLookupFunc sets the environment lookup used to override the Sentry DSN.
func WithEnvLookupFunc(getenv func(string) string) InitOption {
	return func(c *initConfig) {
		if getenv != nil {
			c.getenv = getenv
		}
	}
}

// WithRavenSetDSNFunc sets the function used to configure raven's global DSN.
func WithRavenSetDSNFunc(setDSN func(string) error) InitOption {
	return func(c *initConfig) {
		if setDSN != nil {
			c.setDSN = setDSN
		}
	}
}

// WithSentryClient sets the raven client passed to the Sentry hook factory.
func WithSentryClient(client *raven.Client) InitOption {
	return func(c *initConfig) {
		if client != nil {
			c.sentryClient = client
		}
	}
}

// WithSentryHookFactory sets the factory used to create the Sentry logrus hook.
func WithSentryHookFactory(factory func(*raven.Client, []logrus.Level) (*logrus_sentry.SentryHook, error)) InitOption {
	return func(c *initConfig) {
		if factory != nil {
			c.newSentryHook = factory
		}
	}
}

// WithLogHookAdder sets the function used to register the Sentry hook.
func WithLogHookAdder(addHook func(logrus.Hook)) InitOption {
	return func(c *initConfig) {
		if addHook != nil {
			c.addHook = addHook
		}
	}
}

// WithLogrusConfigurer sets the logrus global configuration functions used during initialization.
func WithLogrusConfigurer(setReportCaller func(bool), setOutput func(io.Writer), setFormatter func(logrus.Formatter)) InitOption {
	return func(c *initConfig) {
		if setReportCaller != nil {
			c.setReportCaller = setReportCaller
		}
		if setOutput != nil {
			c.setOutput = setOutput
		}
		if setFormatter != nil {
			c.setFormatter = setFormatter
		}
	}
}

// WithInitErrorLogger sets the logger used for initialization failures.
func WithInitErrorLogger(logError func(error, string)) InitOption {
	return func(c *initConfig) {
		if logError != nil {
			c.logError = logError
		}
	}
}

func defaultInitErrorLogger(err error, msg string) { logrus.WithError(err).Error(msg) }

func applyInitOptions(dsn string, opts []InitOption) initConfig {
	cfg := initConfig{
		dsn:             dsn,
		envKey:          SentryDSN,
		output:          io.Discard,
		formatter:       EmptyFormatter,
		reportCaller:    true,
		levels:          []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel},
		getenv:          os.Getenv,
		setDSN:          raven.SetDSN,
		sentryClient:    raven.DefaultClient,
		newSentryHook:   logrus_sentry.NewAsyncWithClientSentryHook,
		addHook:         logrus.AddHook,
		setReportCaller: logrus.SetReportCaller,
		setOutput:       logrus.SetOutput,
		setFormatter:    logrus.SetFormatter,
		logError:        defaultInitErrorLogger,
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
	if cfg.getenv == nil {
		cfg.getenv = os.Getenv
	}
	if cfg.setDSN == nil {
		cfg.setDSN = raven.SetDSN
	}
	if cfg.sentryClient == nil {
		cfg.sentryClient = raven.DefaultClient
	}
	if cfg.newSentryHook == nil {
		cfg.newSentryHook = logrus_sentry.NewAsyncWithClientSentryHook
	}
	if cfg.addHook == nil {
		cfg.addHook = logrus.AddHook
	}
	if cfg.setReportCaller == nil {
		cfg.setReportCaller = logrus.SetReportCaller
	}
	if cfg.setOutput == nil {
		cfg.setOutput = logrus.SetOutput
	}
	if cfg.setFormatter == nil {
		cfg.setFormatter = logrus.SetFormatter
	}
	if cfg.logError == nil {
		cfg.logError = defaultInitErrorLogger
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
	cfg.setReportCaller(cfg.reportCaller)
	cfg.setOutput(cfg.output)
	cfg.setFormatter(cfg.formatter)

	if dsn := cfg.getenv(cfg.envKey); dsn != "" {
		cfg.dsn = dsn
	}
	if cfg.dsn == "" {
		return
	}
	if err := cfg.setDSN(cfg.dsn); err != nil {
		cfg.logError(err, "raven init failed")
		return
	}

	sentry, err := cfg.newSentryHook(
		cfg.sentryClient,
		cfg.levels,
	)
	if err != nil {
		cfg.logError(err, "sentry hook init failed")
		return
	}
	sentry.StacktraceConfiguration.Enable = true
	sentry.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	cfg.addHook(sentry)
}

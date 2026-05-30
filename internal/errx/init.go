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

// Init configures logrus to forward logs to the internal logs hook and,
// when a DSN is provided, to Sentry as well.
func Init(sentryDSN string) {
	logrus.SetReportCaller(true)
	logrus.SetOutput(io.Discard)
	logrus.SetFormatter(EmptyFormatter)

	if dsn := os.Getenv(SentryDSN); dsn != "" {
		sentryDSN = dsn
	}
	if sentryDSN == "" {
		return
	}
	if err := raven.SetDSN(sentryDSN); err != nil {
		logrus.WithError(err).Error("raven init failed")
		return
	}

	sentry, err := logrus_sentry.NewAsyncWithClientSentryHook(
		raven.DefaultClient,
		[]logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel},
	)
	if err != nil {
		logrus.WithError(err).Error("sentry hook init failed")
		return
	}
	sentry.StacktraceConfiguration.Enable = true
	sentry.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	logrus.AddHook(sentry)
}

package errx

import (
	"context"

	"github.com/sirupsen/logrus"
)

// MustExit logs err and panics when err is non-nil.
func MustExit(ctx context.Context, err error) {
	if err == nil {
		return
	}
	logrus.WithContext(ctx).WithError(err).Error("exit with error")
	panic(err)
}

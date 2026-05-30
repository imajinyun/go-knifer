package errx

import (
	"context"
	"errors"
	"testing"
)

func TestMustExitNoopOnNilError(t *testing.T) {
	MustExit(context.Background(), nil)
}

func TestMustExitPanicsOnError(t *testing.T) {
	silenceLogrus(t)

	want := errors.New("fatal")
	defer func() {
		got := recover()
		if got != want {
			t.Fatalf("panic = %v, want original error", got)
		}
	}()
	MustExit(context.Background(), want)
}

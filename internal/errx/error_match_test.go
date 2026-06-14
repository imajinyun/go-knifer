package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"
)

func TestErrorIsHandlesNestedMultierror(t *testing.T) {
	target := errors.New("target")
	nested := multierror.Append(nil, errors.New("other"), fmt.Errorf("wrapped: %w", target))
	err := multierror.Append(nil, errors.New("top"), nested)

	if !ErrorIs(err, target) {
		t.Fatalf("ErrorIs() = false, want true for nested multierror")
	}
	if ErrorIs(err, errors.New("missing")) {
		t.Fatal("ErrorIs() = true for an unrelated error")
	}
	if !ErrorIs(nil, nil) {
		t.Fatal("ErrorIs(nil, nil) should be true")
	}
	if ErrorIs(err, nil) {
		t.Fatal("ErrorIs(non-nil, nil) should be false")
	}
}

package bean

import (
	"errors"
	"strconv"
	"testing"
)

func TestBeanErrorContract(t *testing.T) {
	_, err := ToMap(nil)
	assertBeanInvalidInput(t, err)

	err = FillMap(sourceProfile{}, nil)
	assertBeanInvalidInput(t, err)

	var dst targetProfile
	err = CopyProperties(map[string]any{"age": "not-a-number"}, &dst)
	assertBeanInvalidInput(t, err)
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fatalf("CopyProperties should preserve strconv.NumError cause: %v", err)
	}

	err = CopyProperties(map[string]any{"age": "42"}, &dst, WithWeaklyTyped(false))
	assertBeanInvalidInput(t, err)
}

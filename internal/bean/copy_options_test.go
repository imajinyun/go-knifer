package bean

import (
	"strconv"
	"testing"
)

func TestCopyPropertiesWithParserOptions(t *testing.T) {
	src := map[string]any{
		"age":   "custom-int",
		"admin": "custom-bool",
		"score": "custom-float",
		"quota": "custom-uint",
	}
	type target struct {
		Age   int
		Admin bool
		Score float64
		Quota uint
	}
	var dst target
	var intCalled, boolCalled, floatCalled, uintCalled int
	err := CopyProperties(src, &dst,
		WithIntParser(func(text string, base, bits int) (int64, error) {
			intCalled++
			if text == "custom-int" {
				return 42, nil
			}
			return strconv.ParseInt(text, base, bits)
		}),
		WithBoolParser(func(text string) (bool, error) {
			boolCalled++
			return text == "custom-bool", nil
		}),
		WithFloatParser(func(text string, bits int) (float64, error) {
			floatCalled++
			if text == "custom-float" {
				return 9.5, nil
			}
			return strconv.ParseFloat(text, bits)
		}),
		WithUintParser(func(text string, base, bits int) (uint64, error) {
			uintCalled++
			if text == "custom-uint" {
				return 7, nil
			}
			return strconv.ParseUint(text, base, bits)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if dst != (target{Age: 42, Admin: true, Score: 9.5, Quota: 7}) {
		t.Fatalf("CopyProperties dst = %+v", dst)
	}
	if intCalled != 1 || boolCalled != 1 || floatCalled != 1 || uintCalled != 1 {
		t.Fatalf("parser calls int=%d bool=%d float=%d uint=%d", intCalled, boolCalled, floatCalled, uintCalled)
	}
}

func TestWeaklyTypedDisabled(t *testing.T) {
	var dst targetProfile
	err := CopyProperties(map[string]any{"age": "42"}, &dst, WithWeaklyTyped(false))
	if err == nil {
		t.Fatal("expected strict assignment error")
	}
}

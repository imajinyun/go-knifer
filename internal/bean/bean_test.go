package bean

import (
	"errors"
	"strconv"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

type embeddedProfile struct {
	Trace string `bean:"trace_id"`
}

type sourceProfile struct {
	embeddedProfile
	Name  string `bean:"name,alias=full_name|displayName"`
	Age   string `bean:"age"`
	Admin string `bean:"admin"`
	Skip  string `bean:"-"`
	Empty string `bean:"empty"`
}

type targetProfile struct {
	Name  string `bean:"name,alias=full_name|displayName" json:"full_name"`
	Age   int    `json:"age"`
	Admin bool   `json:"admin"`
	Trace string `json:"trace_id"`
	Empty string `json:"empty"`
}

func TestCopyPropertiesStructToStructWithAliasAndWeakConversion(t *testing.T) {
	src := sourceProfile{
		embeddedProfile: embeddedProfile{Trace: "t-1"},
		Name:            "alice",
		Age:             "42",
		Admin:           "yes",
		Skip:            "ignored",
	}
	var dst targetProfile
	if err := CopyProperties(src, &dst, WithIgnoreEmpty(true)); err != nil {
		t.Fatalf("CopyProperties() error = %v", err)
	}
	if dst.Name != "alice" || dst.Age != 42 || !dst.Admin || dst.Trace != "t-1" || dst.Empty != "" {
		t.Fatalf("dst = %+v", dst)
	}
}

func TestCopyPropertiesMapToStruct(t *testing.T) {
	src := map[string]any{
		"displayName": "bob",
		"age":         7.9,
		"admin":       1,
		"trace_id":    "t-2",
	}
	var dst targetProfile
	if err := Copy(src, &dst); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}
	if dst.Name != "bob" || dst.Age != 7 || !dst.Admin || dst.Trace != "t-2" {
		t.Fatalf("dst = %+v", dst)
	}
}

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

func TestToMapUsesPrimaryTagAndOmit(t *testing.T) {
	got, err := ToMap(sourceProfile{Name: "alice", Age: "18", Skip: "hidden"})
	if err != nil {
		t.Fatalf("ToMap() error = %v", err)
	}
	if got["name"] != "alice" || got["age"] != "18" {
		t.Fatalf("map = %#v", got)
	}
	if _, ok := got["Skip"]; ok {
		t.Fatalf("omit field leaked: %#v", got)
	}
}

func TestWeaklyTypedDisabled(t *testing.T) {
	var dst targetProfile
	err := CopyProperties(map[string]any{"age": "42"}, &dst, WithWeaklyTyped(false))
	if err == nil {
		t.Fatal("expected strict assignment error")
	}
}

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

func assertBeanInvalidInput(t *testing.T, err error) {
	t.Helper()
	const code = knifer.ErrCodeInvalidInput
	if err == nil {
		t.Fatalf("err = nil, want %s", code)
	}
	if !errors.Is(err, code) {
		t.Fatalf("errors.Is(%v, %s) = false", err, code)
	}
	got, ok := knifer.CodeOf(err)
	if !ok || got != code {
		t.Fatalf("CodeOf(%v) = %q, %v; want %q, true", err, got, ok, code)
	}
	var beanErr *BeanError
	if !errors.As(err, &beanErr) {
		t.Fatalf("errors.As(err, *BeanError) = false: %v", err)
	}
}

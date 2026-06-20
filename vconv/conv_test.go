package vconv

import (
	"errors"
	"math"
	"testing"

	knifer "github.com/imajinyun/go-knifer"
)

type facadeNamedString string

func TestConvFacade(t *testing.T) {
	if ToString(12) != "12" || ToStringDefault(nil, "x") != "x" {
		t.Fatal("string conversion failed")
	}
	if ToInt("12") != 12 || ToIntDefault("bad", 7) != 7 {
		t.Fatal("int conversion failed")
	}
	if ToInt64(true) != 1 || ToFloat64("3.5") != 3.5 {
		t.Fatal("number conversion failed")
	}
	if !ToBool("yes") || ToBoolDefault("bad", true) != true {
		t.Fatal("bool conversion failed")
	}
	if string(ToBytes("go")) != "go" {
		t.Fatal("bytes conversion failed")
	}
}

func TestConvFacadeDefaultValues(t *testing.T) {
	if got := ToInt64Default("bad", 42); got != 42 {
		t.Fatalf("ToInt64Default = %d, want 42", got)
	}
	if got := ToInt64Default("123", 42); got != 123 {
		t.Fatalf("ToInt64Default = %d, want 123", got)
	}
	if got := ToFloat64Default("bad", 3.14); got != 3.14 {
		t.Fatalf("ToFloat64Default = %v, want 3.14", got)
	}
	if got := ToFloat64Default("2.5", 1.0); got != 2.5 {
		t.Fatalf("ToFloat64Default = %v, want 2.5", got)
	}
}

func TestConvFacadeWithOptions(t *testing.T) {
	if ToStringWithOptions(true, WithFormatBoolFunc(func(bool) string { return "BOOL" })) != "BOOL" {
		t.Fatal("ToStringWithOptions bool formatter failed")
	}
	if ToStringDefaultWithOptions(nil, "fallback", WithFormatBoolFunc(func(bool) string { return "BOOL" })) != "fallback" {
		t.Fatal("ToStringDefaultWithOptions fallback failed")
	}
	parseInt := WithParseIntFunc(func(string, int, int) (int64, error) { return 41, nil })
	if ToIntWithOptions("ignored", parseInt) != 41 || ToIntDefaultWithOptions("ignored", 7, parseInt) != 41 {
		t.Fatal("int option conversion failed")
	}
	if ToInt64WithOptions("ignored", parseInt) != 41 || ToInt64DefaultWithOptions("ignored", 7, parseInt) != 41 {
		t.Fatal("int64 option conversion failed")
	}
	parseFloat := WithParseFloatFunc(func(string, int) (float64, error) { return 6.25, nil })
	if ToFloat64WithOptions("ignored", parseFloat) != 6.25 || ToFloat64DefaultWithOptions("ignored", 7, parseFloat) != 6.25 {
		t.Fatal("float option conversion failed")
	}
	parseBool := WithBoolParser(func(string) (bool, error) { return true, nil })
	if !ToBoolWithOptions("ignored", parseBool) || !ToBoolDefaultWithOptions("ignored", false, parseBool) {
		t.Fatal("bool option conversion failed")
	}
	if got := string(ToBytesWithOptions(3.5, WithFormatFloatFunc(func(float64, byte, int, int) string { return "FLOAT" }))); got != "FLOAT" {
		t.Fatalf("ToBytesWithOptions = %q", got)
	}
}

func TestConvFacadeErrorReturningHelpers(t *testing.T) {
	if got, err := ToIntE("42"); err != nil || got != 42 {
		t.Fatalf("ToIntE = %d, %v", got, err)
	}
	if got, err := ToInt64E(facadeNamedString("77")); err != nil || got != 77 {
		t.Fatalf("ToInt64E named string = %d, %v", got, err)
	}
	if got, err := ToIntEWithOptions("ignored", WithParseIntFunc(func(string, int, int) (int64, error) { return 17, nil })); err != nil || got != 17 {
		t.Fatalf("ToIntEWithOptions = %d, %v", got, err)
	}
	if got, err := ToInt64E("42.9"); err != nil || got != 42 {
		t.Fatalf("ToInt64E = %d, %v", got, err)
	}
	if got, err := ToFloat64E(facadeNamedString("1.25")); err != nil || got != 1.25 {
		t.Fatalf("ToFloat64E named string = %v, %v", got, err)
	}
	if got, err := ToFloat64E("3.5"); err != nil || got != 3.5 {
		t.Fatalf("ToFloat64E = %v, %v", got, err)
	}
	if got, err := ToBoolE("yes"); err != nil || !got {
		t.Fatalf("ToBoolE = %v, %v", got, err)
	}

	for _, err := range []error{
		func() error { _, err := ToIntE("bad"); return err }(),
		func() error { _, err := ToIntE(uint64(math.MaxInt64) + 1); return err }(),
		func() error { _, err := ToInt64E(nil); return err }(),
		func() error { _, err := ToInt64E(uint64(math.MaxInt64) + 1); return err }(),
		func() error { _, err := ToInt64E(math.Inf(1)); return err }(),
		func() error { _, err := ToFloat64E("bad"); return err }(),
		func() error { _, err := ToBoolE("maybe"); return err }(),
	} {
		if !errors.Is(err, ErrInvalidConversion) {
			t.Fatalf("error = %v, want ErrInvalidConversion", err)
		}
		if !errors.Is(err, knifer.ErrCodeInvalidInput) {
			t.Fatalf("error = %v, want invalid input code", err)
		}
	}
}

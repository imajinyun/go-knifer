package vsys_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/imajinyun/go-knifer/vsys"
)

func TestFacadeEnv(t *testing.T) {
	_ = os.Setenv("GO_KNIFER_TEST_KEY", "test_value")
	defer os.Unsetenv("GO_KNIFER_TEST_KEY")

	if got := vsys.Env("GO_KNIFER_TEST_KEY"); got != "test_value" {
		t.Fatalf("expected 'test_value', got %q", got)
	}
	if got := vsys.EnvOrDefault("GO_KNIFER_TEST_MISSING", "default"); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
	if got := vsys.EnvInt("GO_KNIFER_TEST_KEY", 0); got != 0 {
		t.Fatalf("expected 0 for non-int env, got %d", got)
	}
	if got := vsys.EnvBool("GO_KNIFER_TEST_KEY", false); got != false {
		t.Fatalf("expected false for non-bool env, got %v", got)
	}

	lookup := vsys.WithEnvLookupFunc(func(key string) (string, bool) {
		switch key {
		case "A":
			return "value", true
		case "N":
			return "13", true
		case "B":
			return "true", true
		default:
			return "", false
		}
	})
	var warning bytes.Buffer
	if got := vsys.EnvWithOptions("A", lookup); got != "value" {
		t.Fatalf("EnvWithOptions = %q", got)
	}
	if got := vsys.GetWithOptions("MISSING", false, lookup, vsys.WithEnvWarningWriter(&warning)); got != "" || warning.Len() == 0 {
		t.Fatalf("GetWithOptions missing = %q warning=%q", got, warning.String())
	}
	if got := vsys.EnvOrDefaultWithOptions("MISSING", "def", lookup); got != "def" {
		t.Fatalf("EnvOrDefaultWithOptions = %q", got)
	}
	if got := vsys.EnvIntWithOptions("N", 0, lookup); got != 13 {
		t.Fatalf("EnvIntWithOptions = %d", got)
	}
	if got := vsys.EnvBoolWithOptions("B", false, lookup); !got {
		t.Fatalf("EnvBoolWithOptions = %v", got)
	}
}

func TestFacadeDumpSystemInfo(t *testing.T) {
	var buf bytes.Buffer
	vsys.DumpSystemInfoTo(&buf)
	if buf.Len() == 0 {
		t.Fatal("expected non-empty system info dump")
	}
}

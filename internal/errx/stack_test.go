package errx

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetStackTraceAndFormatting(t *testing.T) {
	stack := GetStackTrace(0)
	if len(stack) == 0 {
		t.Fatal("GetStackTrace() returned an empty stack")
	}
	if neg := GetStackTrace(-1); len(neg) == 0 {
		t.Fatal("GetStackTrace(-1) returned an empty stack")
	}

	short := fmt.Sprintf("%v", stack)
	if !strings.HasPrefix(short, "[") || !strings.HasSuffix(short, "]") {
		t.Fatalf("short stack format = %q, want bracketed slice", short)
	}
	verbose := fmt.Sprintf("%+v", stack)
	if !strings.Contains(verbose, "TestGetStackTraceAndFormatting") {
		t.Fatalf("verbose stack format should include test function, got %q", verbose)
	}
	goSyntax := fmt.Sprintf("%#v", stack)
	if !strings.Contains(goSyntax, "errx.Frame") {
		t.Fatalf("go-syntax stack format = %q, want frame type", goSyntax)
	}
}

func TestFrameFormatting(t *testing.T) {
	stack := GetStackTrace(0)
	frame := stack[0]

	if got := fmt.Sprintf("%s", frame); got == "" || got == "unknown" {
		t.Fatalf("frame %%s = %q", got)
	}
	if got := fmt.Sprintf("%d", frame); got == "0" || got == "" {
		t.Fatalf("frame %%d = %q", got)
	}
	if got := fmt.Sprintf("%n", frame); got == "" || got == "unknown" {
		t.Fatalf("frame %%n = %q", got)
	}
	if got := fmt.Sprintf("%+s", frame); !strings.Contains(got, "\n\t") {
		t.Fatalf("frame %%+s = %q, want function and file", got)
	}
}

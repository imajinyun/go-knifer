package vcli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestFacadeRunWithInjectedRunner(t *testing.T) {
	runner := RunnerFunc(func(ctx context.Context, req ExecRequest) (ExecResult, error) {
		if req.Name != "tool" || len(req.Args) != 1 || req.Args[0] != "arg" {
			t.Fatalf("request = %+v", req)
		}
		return ExecResult{Stdout: "ok"}, nil
	})
	got, err := Output(context.Background(), "tool", []string{"arg"}, WithRunner(runner))
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("Output = %q", got)
	}
}

func TestFacadeCommandAndHelp(t *testing.T) {
	var stdout bytes.Buffer
	root := &Command{Name: "app", Usage: "app <command>", Summary: "demo"}
	root.Add(&Command{
		Name:    "echo",
		Summary: "print args",
		Run: func(ctx context.Context, inv *Invocation) error {
			_, _ = fmt.Fprint(inv.Stdout, strings.Join(inv.Args, ","))
			return nil
		},
	})
	if err := root.Execute(context.Background(), []string{"echo", "a", "b"}, WithStdout(&stdout)); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stdout.String() != "a,b" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if help := RenderHelp(root, WithColorMode(ColorNever)); !strings.Contains(help, "app <command>") {
		t.Fatalf("help = %q", help)
	}
}

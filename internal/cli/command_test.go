package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestCommandExecuteRoutesToSubcommand(t *testing.T) {
	var stdout bytes.Buffer
	root := &Command{Name: "app", Summary: "demo app"}
	root.Add(&Command{
		Name:    "hello",
		Summary: "print greeting",
		Run: func(ctx context.Context, inv *Invocation) error {
			_, _ = fmt.Fprintf(inv.Stdout, "hello %s", inv.Args[0])
			return nil
		},
	})
	err := root.Execute(context.Background(), []string{"hello", "gopher"}, WithStdout(&stdout))
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if stdout.String() != "hello gopher" {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestCommandExecuteUnknownCommandReturnsUsageError(t *testing.T) {
	root := &Command{Name: "app"}
	err := root.Execute(context.Background(), []string{"missing"})
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("Execute unknown command error = %v, want ErrUsage", err)
	}
}

func TestCommandExecuteParsesLocalFlags(t *testing.T) {
	var gotName string
	cmd := &Command{Name: "serve"}
	cmd.Flags = func(flags *FlagParser) {
		name := flags.String("name", "world", "name to greet")
		cmd.Run = func(ctx context.Context, inv *Invocation) error {
			gotName = *name
			return nil
		}
	}
	err := cmd.Execute(context.Background(), []string{"--name", "gopher"})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if gotName != "gopher" {
		t.Fatalf("flag value = %q", gotName)
	}
}

func TestCommandExecutePropagatesContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd := &Command{
		Name: "root",
		Run: func(ctx context.Context, inv *Invocation) error {
			return ctx.Err()
		},
	}
	err := cmd.Execute(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Execute canceled error = %v, want context.Canceled", err)
	}
}

func TestCommandExecuteNilCommandReturnsUsageError(t *testing.T) {
	var cmd *Command
	err := cmd.Execute(context.Background(), nil)
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("nil command error = %v, want ErrUsage", err)
	}
}

func TestCommandExecuteCommandWithoutNameReturnsUsageError(t *testing.T) {
	cmd := &Command{}
	err := cmd.Execute(context.Background(), nil)
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("unnamed command error = %v, want ErrUsage", err)
	}
}

func TestCommandExecuteRendersHelpWhenNoRun(t *testing.T) {
	var stdout bytes.Buffer
	cmd := &Command{Name: "app", Summary: "demo app"}
	err := cmd.Execute(context.Background(), nil, WithStdout(&stdout))
	if err != nil {
		t.Fatalf("Execute help returned error: %v", err)
	}
	if got := stdout.String(); got != "Usage: app\n\ndemo app\n" {
		t.Fatalf("help = %q", got)
	}
}

func TestCommandExecuteRejectsUnexpectedArgsWhenNoRun(t *testing.T) {
	cmd := &Command{Name: "app"}
	err := cmd.Execute(context.Background(), []string{"extra"})
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("unexpected args error = %v, want ErrUsage", err)
	}
}

func TestCommandExecuteWritesFlagErrorsToStderr(t *testing.T) {
	var stderr bytes.Buffer
	cmd := &Command{Name: "serve"}
	cmd.Flags = func(flags *FlagParser) {
		flags.Bool("debug", false, "enable debug")
	}
	err := cmd.Execute(context.Background(), []string{"--missing"}, WithStderr(&stderr))
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("flag error = %v, want ErrUsage", err)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("flag provided but not defined")) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

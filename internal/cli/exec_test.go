package cli

import (
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fakeRunner struct {
	requests []ExecRequest
	result   ExecResult
	err      error
}

func (r *fakeRunner) Run(ctx context.Context, req ExecRequest) (ExecResult, error) {
	r.requests = append(r.requests, req)
	select {
	case <-ctx.Done():
		return ExecResult{}, ctx.Err()
	default:
	}
	return r.result, r.err
}

func TestRunUsesInjectedRunnerAndClonesMutableInputs(t *testing.T) {
	runner := &fakeRunner{result: ExecResult{Stdout: "ok", ExitCode: 0}}
	env := []string{"A=1"}
	result, err := Run(
		context.Background(),
		"git",
		[]string{"status"},
		WithRunner(runner),
		WithDir("/tmp/work"),
		WithEnv(env),
		WithStdin(strings.NewReader("input")),
	)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.Stdout != "ok" || result.ExitCode != 0 {
		t.Fatalf("Run result = %+v", result)
	}
	env[0] = "A=changed"
	if len(runner.requests) != 1 {
		t.Fatalf("runner requests = %d", len(runner.requests))
	}
	got := runner.requests[0]
	if got.Name != "git" || !reflect.DeepEqual(got.Args, []string{"status"}) || got.Dir != "/tmp/work" {
		t.Fatalf("request = %+v", got)
	}
	if !reflect.DeepEqual(got.Env, []string{"A=1"}) {
		t.Fatalf("request env was not cloned: %#v", got.Env)
	}
	if got.Stdin == nil {
		t.Fatalf("request stdin is nil")
	}
}

func TestRunRejectsEmptyCommandName(t *testing.T) {
	_, err := Run(context.Background(), "", nil)
	if !errors.Is(err, ErrEmptyCommand) {
		t.Fatalf("Run empty command error = %v, want ErrEmptyCommand", err)
	}
}

func TestRunPropagatesContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Run(ctx, "tool", nil, WithRunner(&fakeRunner{}))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run canceled error = %v, want context.Canceled", err)
	}
}

func TestRunWrapsRunnerError(t *testing.T) {
	runnerErr := errors.New("boom")
	_, err := Run(context.Background(), "tool", nil, WithRunner(&fakeRunner{err: runnerErr}))
	if !errors.Is(err, runnerErr) {
		t.Fatalf("Run error = %v, want wrapped runner error", err)
	}
}

func TestRunEnforcesOutputLimitOnInjectedRunner(t *testing.T) {
	runner := &fakeRunner{result: ExecResult{Stdout: "abcdef", Stderr: "xyz"}}
	result, err := Run(context.Background(), "tool", nil, WithRunner(runner), WithMaxOutputBytes(5))
	if !errors.Is(err, ErrOutputLimitExceeded) {
		t.Fatalf("Run limit error = %v, want ErrOutputLimitExceeded", err)
	}
	if result.Stdout != "abcde" || result.Stderr != "" {
		t.Fatalf("limited result = %+v", result)
	}
}

func TestOutputReturnsStdout(t *testing.T) {
	runner := &fakeRunner{result: ExecResult{Stdout: "hello\n", ExitCode: 0}}
	got, err := Output(context.Background(), "printf", []string{"hello"}, WithRunner(runner))
	if err != nil {
		t.Fatalf("Output returned error: %v", err)
	}
	if got != "hello\n" {
		t.Fatalf("Output = %q", got)
	}
}

func TestWithTimeoutAppliesDeadline(t *testing.T) {
	runner := RunnerFunc(func(ctx context.Context, req ExecRequest) (ExecResult, error) {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatalf("ctx has no deadline")
		}
		if time.Until(deadline) > time.Second {
			t.Fatalf("deadline too far away: %v", deadline)
		}
		return ExecResult{Stdout: "ok"}, nil
	})
	_, err := Run(context.Background(), "tool", nil, WithRunner(runner), WithTimeout(time.Second))
	if err != nil {
		t.Fatalf("Run with timeout returned error: %v", err)
	}
}

func TestExecRequestStdinAcceptsNilAndReader(t *testing.T) {
	var _ io.Reader = strings.NewReader("x")
	request := ExecRequest{Name: "tool", Stdin: strings.NewReader("x")}
	if request.Stdin == nil {
		t.Fatalf("stdin reader is nil")
	}
}

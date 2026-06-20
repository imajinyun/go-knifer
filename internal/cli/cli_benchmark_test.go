package cli

import (
	"context"
	"io"
	"testing"
	"time"
)

func BenchmarkFlagParserParse(b *testing.B) {
	for b.Loop() {
		parser := NewFlagParser("serve")
		parser.String("host", "127.0.0.1", "host to bind")
		parser.Int("port", 8080, "port to bind")
		parser.Bool("debug", false, "enable debug")
		parser.Duration("timeout", time.Second, "request timeout")
		_, _ = parser.Parse([]string{
			"--host", "0.0.0.0",
			"--port", "9090",
			"--debug",
			"--timeout", "2s",
			"api",
		})
	}
}

func BenchmarkRenderHelp(b *testing.B) {
	root := &Command{Name: "app", Usage: "app <command>", Summary: "demo app"}
	root.Add(&Command{Name: "serve", Summary: "start server"})
	root.Add(&Command{Name: "version", Summary: "print version"})
	for b.Loop() {
		_ = RenderHelp(root, WithColorMode(ColorNever))
	}
}

func BenchmarkRunInjectedRunner(b *testing.B) {
	runner := RunnerFunc(func(ctx context.Context, req ExecRequest) (ExecResult, error) {
		return ExecResult{Stdout: "ok"}, nil
	})
	for b.Loop() {
		_, _ = Run(
			context.Background(),
			"tool",
			[]string{"arg"},
			WithRunner(runner),
			WithStdin(io.Reader(nil)),
		)
	}
}

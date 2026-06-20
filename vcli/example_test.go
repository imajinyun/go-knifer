package vcli_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imajinyun/go-knifer/vcli"
)

func ExampleRun() {
	runner := vcli.RunnerFunc(func(ctx context.Context, req vcli.ExecRequest) (vcli.ExecResult, error) {
		return vcli.ExecResult{Stdout: req.Name + " " + strings.Join(req.Args, " ")}, nil
	})
	result, err := vcli.Run(context.Background(), "echo", []string{"hello"}, vcli.WithRunner(runner))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Stdout)
	// Output: echo hello
}

func ExampleNewFlagParser() {
	parser := vcli.NewFlagParser("serve")
	port := parser.Int("port", 8080, "port to bind")
	debug := parser.Bool("debug", false, "enable debug")
	result, err := parser.Parse([]string{"--port", "9090", "--debug", "api"})
	if err != nil {
		panic(err)
	}
	fmt.Println(*port, *debug, result.Args[0])
	// Output: 9090 true api
}

func ExampleCommand_Execute() {
	cmd := &vcli.Command{
		Name: "hello",
		Run: func(ctx context.Context, inv *vcli.Invocation) error {
			_, _ = fmt.Fprintf(inv.Stdout, "hello %s", inv.Args[0])
			return nil
		},
	}
	var out strings.Builder
	if err := cmd.Execute(context.Background(), []string{"gopher"}, vcli.WithStdout(&out)); err != nil {
		panic(err)
	}
	fmt.Println(out.String())
	// Output: hello gopher
}

func ExampleRenderHelp() {
	root := &vcli.Command{Name: "app", Usage: "app <command>", Summary: "demo app"}
	root.Add(&vcli.Command{Name: "serve", Summary: "start server"})
	fmt.Print(vcli.RenderHelp(root, vcli.WithColorMode(vcli.ColorNever)))
	// Output:
	// Usage: app <command>
	//
	// demo app
	//
	// Commands:
	//   serve	start server
}

func ExampleColorize() {
	fmt.Println(vcli.Colorize("ok", vcli.Green, vcli.WithColorMode(vcli.ColorNever)))
	fmt.Println(strings.Contains(vcli.Colorize("ok", vcli.Green, vcli.WithColorMode(vcli.ColorAlways)), "\x1b["))
	// Output:
	// ok
	// true
}

func ExampleWithTimeout() {
	runner := vcli.RunnerFunc(func(ctx context.Context, req vcli.ExecRequest) (vcli.ExecResult, error) {
		_, hasDeadline := ctx.Deadline()
		return vcli.ExecResult{Stdout: fmt.Sprint(hasDeadline)}, nil
	})
	stdout, err := vcli.Output(
		context.Background(),
		"tool",
		nil,
		vcli.WithRunner(runner),
		vcli.WithTimeout(time.Second),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(stdout)
	// Output: true
}

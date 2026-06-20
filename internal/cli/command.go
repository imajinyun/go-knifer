package cli

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Handler runs a command invocation.
type Handler func(context.Context, *Invocation) error

// Command describes a lightweight command or subcommand.
type Command struct {
	Name     string
	Summary  string
	Usage    string
	Flags    func(*FlagParser)
	Run      Handler
	Children []*Command
}

// Invocation contains parsed arguments and command I/O streams.
type Invocation struct {
	Command *Command
	Args    []string
	Stdout  io.Writer
	Stderr  io.Writer
}

type executeConfig struct {
	stdout io.Writer
	stderr io.Writer
}

// ExecuteOption customizes command execution.
type ExecuteOption func(*executeConfig)

// WithStdout sets command stdout.
func WithStdout(w io.Writer) ExecuteOption {
	return func(c *executeConfig) {
		if w != nil {
			c.stdout = w
		}
	}
}

// WithStderr sets command stderr.
func WithStderr(w io.Writer) ExecuteOption {
	return func(c *executeConfig) {
		if w != nil {
			c.stderr = w
		}
	}
}

func applyExecuteOptions(opts []ExecuteOption) executeConfig {
	cfg := executeConfig{stdout: os.Stdout, stderr: os.Stderr}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// Add appends subcommands to cmd.
func (cmd *Command) Add(children ...*Command) {
	for _, child := range children {
		if child != nil {
			cmd.Children = append(cmd.Children, child)
		}
	}
}

// Execute routes args to cmd or one of its subcommands.
func (cmd *Command) Execute(ctx context.Context, args []string, opts ...ExecuteOption) error {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg := applyExecuteOptions(opts)
	return cmd.execute(ctx, append([]string(nil), args...), cfg)
}

func (cmd *Command) execute(ctx context.Context, args []string, cfg executeConfig) error {
	if cmd == nil || cmd.Name == "" {
		return fmt.Errorf("execute command: %w", ErrUsage)
	}
	if len(args) > 0 {
		for _, child := range cmd.Children {
			if child != nil && child.Name == args[0] {
				return child.execute(ctx, args[1:], cfg)
			}
		}
		if len(cmd.Children) > 0 && cmd.Run == nil {
			return fmt.Errorf("unknown command %q: %w", args[0], ErrUsage)
		}
	}
	parser := NewFlagParser(cmd.Name, WithFlagOutput(cfg.stderr))
	if cmd.Flags != nil {
		cmd.Flags(parser)
	}
	parsed, err := parser.Parse(args)
	if err != nil {
		return err
	}
	if cmd.Run == nil && len(parsed.Args) > 0 {
		return fmt.Errorf("unexpected arguments for %q: %w", cmd.Name, ErrUsage)
	}
	if cmd.Run == nil {
		_, _ = io.WriteString(cfg.stdout, RenderHelp(cmd))
		return nil
	}
	return cmd.Run(ctx, &Invocation{
		Command: cmd,
		Args:    parsed.Args,
		Stdout:  cfg.stdout,
		Stderr:  cfg.stderr,
	})
}

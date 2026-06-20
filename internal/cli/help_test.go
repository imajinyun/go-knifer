package cli

import (
	"strings"
	"testing"
)

func TestRenderHelpIncludesUsageSummaryAndCommands(t *testing.T) {
	root := &Command{Name: "app", Summary: "demo app", Usage: "app <command>"}
	root.Add(&Command{Name: "serve", Summary: "start server"})
	help := RenderHelp(root, WithColorMode(ColorNever))
	for _, want := range []string{"Usage: app <command>", "demo app", "serve", "start server"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help %q does not contain %q", help, want)
		}
	}
}

func TestColorizeHonorsModes(t *testing.T) {
	plain := Colorize("hello", Green, WithColorMode(ColorNever))
	if plain != "hello" {
		t.Fatalf("ColorNever = %q", plain)
	}
	colored := Colorize("hello", Green, WithColorMode(ColorAlways))
	if !strings.Contains(colored, "\x1b[") || !strings.Contains(colored, "hello") {
		t.Fatalf("ColorAlways = %q", colored)
	}
}

func TestColorizeReturnsPlainTextForEmptyColor(t *testing.T) {
	got := Colorize("hello", "", WithColorMode(ColorAlways))
	if got != "hello" {
		t.Fatalf("Colorize empty color = %q", got)
	}
}

func TestRenderHelpNilCommandReturnsEmptyString(t *testing.T) {
	if got := RenderHelp(nil); got != "" {
		t.Fatalf("RenderHelp(nil) = %q", got)
	}
}

func TestRenderHelpUsesDefaultCommandUsage(t *testing.T) {
	cmd := &Command{Name: "app"}
	if got := RenderHelp(cmd); got != "Usage: app\n" {
		t.Fatalf("RenderHelp default usage = %q", got)
	}
}

func TestRenderHelpSkipsNilChildrenAndUsesCommandPlaceholder(t *testing.T) {
	cmd := &Command{Name: "app"}
	cmd.Add(nil, &Command{Name: "serve"})
	help := RenderHelp(cmd)
	for _, want := range []string{"Usage: app <command>", "Commands:", "serve"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help %q does not contain %q", help, want)
		}
	}
}

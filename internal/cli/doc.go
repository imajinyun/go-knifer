// Package cli implements dependency-free helpers for lightweight command-line tooling.
//
// Command execution helpers accept a command name plus an argument slice and never
// execute through a shell by default. Tests should inject a Runner instead of
// starting real processes.
package cli

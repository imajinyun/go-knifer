// Package log provides a unified logging interface and default console implementations matching the utility toolkit-log.
//
// Main types:
//   - Level：log levels: Trace, Debug, Info, Warn, Error, Fatal, All, and Off.
//   - Log：unified logging interface with methods such as Trace, Debug, Info, Warn, and Error.
//   - LogFactory：constructs Log instances by name.
//   - StaticLog：static convenience methods for printing logs directly.
//   - ConsoleLog / ConsoleColorLog：default console implementations, with optional colored output.
package log

package vregex

import regeximpl "github.com/imajinyun/go-knifer/internal/regex"

func Match(pattern, s string) bool       { return regeximpl.ReMatch(pattern, s) }
func Find(pattern, s string) string      { return regeximpl.ReFind(pattern, s) }
func FindAll(pattern, s string) []string { return regeximpl.ReFindAll(pattern, s) }
func Replace(pattern, s, replacement string) string {
	return regeximpl.ReReplace(pattern, s, replacement)
}

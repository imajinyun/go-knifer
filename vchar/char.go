package vchar

import charimpl "github.com/imajinyun/go-knifer/internal/char"

func IsBlankChar(r rune) bool     { return charimpl.IsBlankChar(r) }
func IsLetter(r rune) bool        { return charimpl.IsLetter(r) }
func IsDigit(r rune) bool         { return charimpl.IsDigit(r) }
func IsAscii(r rune) bool         { return charimpl.IsAscii(r) }
func IsLetterOrDigit(r rune) bool { return charimpl.IsLetterOrDigit(r) }

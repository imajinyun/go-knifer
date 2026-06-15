package vimg_test

import "io"

type fixedGenerator struct{ code string }

func (g fixedGenerator) Gen() string { return g.code }

func (g fixedGenerator) Verify(code, userInput string) bool { return code == userInput }

type nopWriteCloser struct{ io.Writer }

func (w nopWriteCloser) Close() error { return nil }

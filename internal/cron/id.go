package cron

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func newIDGeneratorWithReader(r io.Reader) func() string {
	if r == nil {
		r = rand.Reader
	}
	return func() string { return generateIDWithReader(r) }
}

func generateIDWithReader(r io.Reader) string {
	if r == nil {
		r = rand.Reader
	}
	var b [8]byte
	_, _ = io.ReadFull(r, b[:])
	return hex.EncodeToString(b[:])
}

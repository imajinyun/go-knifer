package cron

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// generateID creates a random hexadecimal id with 16 characters.
func generateID() string {
	return generateIDWithReader(rand.Reader)
}

func generateIDWithReader(r io.Reader) string {
	if r == nil {
		r = rand.Reader
	}
	var b [8]byte
	_, _ = io.ReadFull(r, b[:])
	return hex.EncodeToString(b[:])
}

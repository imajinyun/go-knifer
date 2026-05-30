package cron

import (
	"crypto/rand"
	"encoding/hex"
)

// generateID creates a random hexadecimal id with 16 characters.
func generateID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

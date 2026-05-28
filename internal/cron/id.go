package cron

import (
	"crypto/rand"
	"encoding/hex"
)

// generateID 生成一个随机十六进制 id（16 个字符）。
func generateID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

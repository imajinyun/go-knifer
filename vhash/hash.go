package vhash

import hashimpl "github.com/imajinyun/go-knifer/internal/hash"

func AdditiveHash(s string, prime int) int { return hashimpl.AdditiveHash(s, prime) }
func FnvHash(s string) uint32              { return hashimpl.FnvHash(s) }
func MD5Hex(s string) string               { return hashimpl.MD5Hex(s) }
func SHA1Hex(s string) string              { return hashimpl.SHA1Hex(s) }
func SHA256Hex(s string) string            { return hashimpl.SHA256Hex(s) }

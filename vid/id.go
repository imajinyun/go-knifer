package vid

import idimpl "github.com/imajinyun/go-knifer/internal/id"

func SimpleUUID() string   { return idimpl.SimpleUUID() }
func FastUUID() string     { return idimpl.FastUUID() }
func UUID() string         { return idimpl.SimpleUUID() }
func ObjectId() string     { return idimpl.ObjectId() }
func NanoId() string       { return idimpl.NanoId() }
func NanoIdN(n int) string { return idimpl.NanoIdN(n) }

package conf

import (
	"io/fs"
	"time"
)

type watchTestTicker struct {
	stopped chan struct{}
}

func (t *watchTestTicker) Stop() { close(t.stopped) }

type fakeFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return f.size }
func (f fakeFileInfo) Mode() fs.FileMode  { return 0o644 }
func (f fakeFileInfo) ModTime() time.Time { return f.modTime }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }

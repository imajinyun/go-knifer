package system

import (
	"bytes"
	"strings"
	"testing"
)

func TestDumpSystemInfo(t *testing.T) {
	var buf bytes.Buffer
	DumpSystemInfoTo(&buf)
	out := buf.String()
	for _, kw := range []string{"Go Version:", "OS Name:", "User Name:", "Host Name:", "Goroutine Count:"} {
		if !strings.Contains(out, kw) {
			t.Errorf("Dump 输出缺少 %q：\n%s", kw, out)
		}
	}
}

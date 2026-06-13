package system

import (
	"os"
	"testing"
)

func TestReadableSize(t *testing.T) {
	cases := []struct {
		in   uint64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
	}
	for _, c := range cases {
		got := readableSize(c.in)
		if got != c.want {
			t.Errorf("readableSize(%d): 期望 %q 实际 %q", c.in, c.want, got)
		}
	}
}

func TestFixPath(t *testing.T) {
	if fixPath("") != "" {
		t.Errorf("空字符串应保持空")
	}
	sep := string(os.PathSeparator)
	if fixPath("/tmp"+sep) != "/tmp"+sep {
		t.Errorf("已带后缀不应再追加")
	}
	if fixPath("/tmp") != "/tmp"+sep {
		t.Errorf("应追加分隔符")
	}
}

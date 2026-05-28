package system

import "os"

// osHostname 是 os.Hostname 的别名，便于测试覆盖。
func osHostname() (string, error) {
	return os.Hostname()
}

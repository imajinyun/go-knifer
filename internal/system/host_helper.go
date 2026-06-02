package system

import "os"

// osHostname aliases os.Hostname so tests can override it.
func osHostname() (string, error) {
	return os.Hostname()
}

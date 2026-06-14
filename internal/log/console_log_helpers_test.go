package log

import "bytes"

func newTestConsoleLog(name string) (*ConsoleLog, *bytes.Buffer, *bytes.Buffer) {
	c := NewConsoleLog(name)
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	c.SetOutput(out, errOut)
	return c, out, errOut
}

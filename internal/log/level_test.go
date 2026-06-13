package log

import "testing"

func TestLevelString(t *testing.T) {
	cases := map[Level]string{
		LevelAll:   "ALL",
		LevelTrace: "TRACE",
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
		LevelOff:   "OFF",
	}
	for l, want := range cases {
		if got := l.String(); got != want {
			t.Errorf("Level(%d).String()=%q, want %q", l, got, want)
		}
	}
}

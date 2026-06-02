package cron

import "time"

// Config is aligned with the utility toolkit CronConfig and configures the scheduler.
type Config struct {
	// Location is the scheduler time zone.
	Location *time.Location
	// MatchSecond reports whether expressions match seconds; when false, tasks fire by minute.
	MatchSecond bool
}

// NewConfig creates the default config using the local time zone and minute-level matching.
func NewConfig() *Config {
	return &Config{Location: time.Local, MatchSecond: false}
}

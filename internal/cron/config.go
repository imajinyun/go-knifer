package cron

import "time"

// Config is aligned with the utility toolkit CronConfig and configures the scheduler.
type Config struct {
	// Location is the scheduler time zone.
	Location *time.Location
	// MatchSecond reports whether expressions match seconds; when false, tasks fire by minute.
	MatchSecond bool
}

// ConfigOption customizes Config construction.
type ConfigOption func(*Config)

// WithConfigLocation sets the scheduler time zone on Config.
func WithConfigLocation(loc *time.Location) ConfigOption {
	return func(c *Config) {
		if loc != nil {
			c.Location = loc
		}
	}
}

// WithConfigMatchSecond sets whether expressions match seconds on Config.
func WithConfigMatchSecond(matchSecond bool) ConfigOption {
	return func(c *Config) { c.MatchSecond = matchSecond }
}

// NewConfig creates the default config using the local time zone and minute-level matching.
func NewConfig() *Config {
	return NewConfigWithOptions()
}

// NewConfigWithOptions creates a config customized by options.
func NewConfigWithOptions(opts ...ConfigOption) *Config {
	cfg := &Config{Location: time.Local, MatchSecond: false}
	for _, opt := range opts {
		if opt != nil {
			opt(cfg)
		}
	}
	if cfg.Location == nil {
		cfg.Location = time.Local
	}
	return cfg
}

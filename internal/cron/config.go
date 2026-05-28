package cron

import "time"

// Config 对应 hutool 的 CronConfig，调度器配置。
type Config struct {
	// Location 时区。
	Location *time.Location
	// MatchSecond 是否匹配到秒；为 false 时按分钟触发。
	MatchSecond bool
}

// NewConfig 创建默认配置（本地时区，按分钟触发）。
func NewConfig() *Config {
	return &Config{Location: time.Local, MatchSecond: false}
}

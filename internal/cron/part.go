package cron

// Part 表示 Cron 表达式的字段位置，对应 hutool 的 cn.hutool.cron.pattern.Part。
type Part int

const (
	// PartSecond 秒。
	PartSecond Part = iota
	// PartMinute 分。
	PartMinute
	// PartHour 时。
	PartHour
	// PartDayOfMonth 日（1-31，32 为 L 哨兵）。
	PartDayOfMonth
	// PartMonth 月（1-12）。
	PartMonth
	// PartDayOfWeek 周（0=周日 ~ 6=周六）。
	PartDayOfWeek
	// PartYear 年。
	PartYear
)

// partInfo 描述每个 Part 的取值范围。
type partInfo struct {
	min int
	max int
}

var partInfos = [...]partInfo{
	PartSecond:     {0, 59},
	PartMinute:     {0, 59},
	PartHour:       {0, 23},
	PartDayOfMonth: {1, 32}, // 32 为 "L" 的哨兵值
	PartMonth:      {1, 12},
	PartDayOfWeek:  {0, 6},
	PartYear:       {1970, 2099},
}

// Min 返回字段最小值。
func (p Part) Min() int { return partInfos[p].min }

// Max 返回字段最大值。
func (p Part) Max() int { return partInfos[p].max }

// CheckValue 校验值是否合法，越界返回错误。
func (p Part) CheckValue(v int) error {
	if v < p.Min() || v > p.Max() {
		return NewCronError("value %d out of range [%d, %d] for part %d", v, p.Min(), p.Max(), int(p))
	}
	return nil
}

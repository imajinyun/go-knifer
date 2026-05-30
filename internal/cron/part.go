package cron

// Part represents a cron expression field position aligned with hutool cn.hutool.cron.pattern.Part.
type Part int

const (
	// PartSecond is the seconds field.
	PartSecond Part = iota
	// PartMinute is the minutes field.
	PartMinute
	// PartHour is the hours field.
	PartHour
	// PartDayOfMonth is the day-of-month field, where 32 is the L sentinel.
	PartDayOfMonth
	// PartMonth is the month field, from 1 to 12.
	PartMonth
	// PartDayOfWeek is the day-of-week field, from 0 for Sunday to 6 for Saturday.
	PartDayOfWeek
	// PartYear is the year field.
	PartYear
)

// partInfo describes the value range of each Part.
type partInfo struct {
	min int
	max int
}

var partInfos = [...]partInfo{
	PartSecond:     {0, 59},
	PartMinute:     {0, 59},
	PartHour:       {0, 23},
	PartDayOfMonth: {1, 32}, // 32 is the sentinel value for "L".
	PartMonth:      {1, 12},
	PartDayOfWeek:  {0, 6},
	PartYear:       {1970, 2099},
}

// Min returns the minimum field value.
func (p Part) Min() int { return partInfos[p].min }

// Max returns the maximum field value.
func (p Part) Max() int { return partInfos[p].max }

// CheckValue validates the value and returns an error when it is out of range.
func (p Part) CheckValue(v int) error {
	if v < p.Min() || v > p.Max() {
		return NewCronError("value %d out of range [%d, %d] for part %d", v, p.Min(), p.Max(), int(p))
	}
	return nil
}

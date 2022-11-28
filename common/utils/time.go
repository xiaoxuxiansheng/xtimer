package utils

import (
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/consts"
)

func GetStartMinute(timeStr string) (time.Time, error) {
	return time.ParseInLocation(consts.MinuteFormat, timeStr, time.Local)
}

func GetDayStr(t time.Time) string {
	return t.Format(consts.DayFormat)
}

func GetStartHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local)
}

func GetMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

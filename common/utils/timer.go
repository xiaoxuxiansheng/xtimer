package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/consts"
)

func UnionTimerIDUnix(timeID uint, unix int64) string {
	return fmt.Sprintf("%d_%d", timeID, unix)
}

func SplitTimerIDUnix(str string) (uint, int64, error) {
	timerIDUnix := strings.Split(str, "_")
	if len(timerIDUnix) != 2 {
		return 0, 0, fmt.Errorf("invalid timerID unix str: %s", str)
	}

	timerID, _ := strconv.ParseInt(timerIDUnix[0], 10, 64)
	unix, _ := strconv.ParseInt(timerIDUnix[1], 10, 64)
	return uint(timerID), unix, nil
}

func GetTaskBloomFilterKeyByDay(dayStr string) string {
	return "task_bloom_" + dayStr
}

func GetTimeBucketLockKey(t time.Time, bucketID int) string {
	return fmt.Sprintf("time_bucket_lock_%s_%d", t.Format(consts.MinuteFormat), bucketID)
}

func GetMigratorLockKey(t time.Time) string {
	return fmt.Sprintf("migrator_lock_%s", t.Format(consts.HourFormat))
}

func GetSliceMsgKey(t time.Time, bucketID int) string {
	return fmt.Sprintf("%s_%d", t.Format(consts.MinuteFormat), bucketID)
}

func GetForwardTwoMigrateStepEnd(cur time.Time, diff time.Duration) time.Time {
	end := cur.Add(diff)
	return time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, time.Local)
}

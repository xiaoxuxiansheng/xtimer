package cron

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type CronParser struct {
}

func NewCronParser() *CronParser {
	return &CronParser{}
}

func (c *CronParser) NextFromNow(cron string) (time.Time, error) {
	expr, err := cronexpr.Parse(cron)
	if err != nil {
		return time.Time{}, err
	}

	nextTime := expr.Next(time.Now())
	if nextTime.UnixNano() < 0 {
		return time.Time{}, fmt.Errorf("fail to parse time from cron: %s", cron)
	}

	return nextTime, nil
}

func (c *CronParser) NextsBefore(cron string, end time.Time) ([]time.Time, error) {
	return c.NextsBetween(cron, time.Now(), end)
}

func (c *CronParser) NextsBetween(cron string, start, end time.Time) ([]time.Time, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end can not earlier than start, start: %v, end: %v", start, end)
	}

	expr, err := cronexpr.Parse(cron)
	if err != nil {
		return nil, err
	}

	var nexts []time.Time
	for start.Before(end) {
		next := expr.Next(start)
		if next.UnixNano() < 0 {
			return nil, fmt.Errorf("fail to parse time from cron: %s", cron)
		}
		nexts = append(nexts, next)
		start = next
	}

	return nexts, nil
}

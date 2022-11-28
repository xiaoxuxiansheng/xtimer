package task

import (
	"time"

	"gorm.io/gorm"
)

type Option func(*gorm.DB) *gorm.DB

func WithTaskID(id uint) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("id = ?", id)
	}
}

func WithTimerID(timerID uint) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("timer_id = ?", timerID)
	}
}

func WithRunTimer(runTimer time.Time) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("run_timer = ?", runTimer)
	}
}

func WithStartTime(start time.Time) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("run_timer >= ?", start)
	}
}

func WithEndTime(end time.Time) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("run_timer < ?", end)
	}
}

func WithStatus(status int32) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("status = ?", status)
	}
}

func WithAsc() Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Order("created_at ASC")
	}
}

func WithDesc() Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Order("run_timer DESC")
	}
}

func WithPageLimit(offset, limit int) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Offset(offset).Limit(limit)
	}
}

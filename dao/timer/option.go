package timer

import (
	"gorm.io/gorm"
)

type Option func(*gorm.DB) *gorm.DB

func WithID(id uint) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("id = ?", id)
	}
}

func WithIDs(ids []uint) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("id IN ?", ids)
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
		return d.Order("created_at DESC")
	}
}

func WithApp(app string) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("app = ?", app)
	}
}

func WithPageLimit(offset, limit int) Option {
	return func(d *gorm.DB) *gorm.DB {
		return d.Offset(offset).Limit(limit)
	}
}

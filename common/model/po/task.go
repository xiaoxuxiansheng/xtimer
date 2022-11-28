package po

import (
	"time"

	"gorm.io/gorm"
)

// Task 运行流水记录
type Task struct {
	gorm.Model
	App      string    `gorm:"column:app;NOT NULL"`           // 定义ID
	TimerID  uint      `gorm:"column:timer_id;NOT NULL"`      // 定义ID
	Output   string    `gorm:"column:output;default:null"`    // 执行结果
	RunTimer time.Time `gorm:"column:run_timer;default:null"` // 执行时间
	CostTime int       `gorm:"column:cost_time"`              // 执行耗时
	Status   int       `gorm:"column:status;NOT NULL"`        // 当前状态
}

func (t *Task) TableName() string {
	return "task"
}

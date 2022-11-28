package po

import (
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/consts"

	"gorm.io/gorm"
)

// Timer 定时器定义
type Timer struct {
	gorm.Model
	App             string `gorm:"column:app;NOT NULL" json:"app,omitempty"`                             // 定时器定义名称
	Name            string `gorm:"column:name;NOT NULL" json:"name,omitempty"`                           // 定时器定义名称
	Creator         string `gorm:"column:creator;NOT NULL" json:"creator,omitempty"`                     // 创建人
	Status          int    `gorm:"column:status;NOT NULL" json:"status,omitempty"`                       // 定时器定义状态，1:未激活, 2:已激活
	Cron            string `gorm:"column:cron;NOT NULL" json:"cron,omitempty"`                           // 定时器定时配置
	NotifyHTTPParam string `gorm:"column:notify_http_param;NOT NULL" json:"notify_http_param,omitempty"` // Http 回调参数
}

func (t *Timer) TableName() string {
	return "timer"
}

func (t *Timer) BatchTasksFromTimer(executeTimes []time.Time) []*Task {
	tasks := make([]*Task, 0, len(executeTimes))
	for _, executeTime := range executeTimes {
		tasks = append(tasks, &Task{
			App:      t.App,
			TimerID:  t.Model.ID,
			Status:   consts.NotRunned.ToInt(),
			RunTimer: executeTime,
		})
	}
	return tasks
}

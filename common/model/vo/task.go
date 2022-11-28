package vo

import (
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
)

type GetTasksReq struct {
	PageLimiter
	TimerID uint `form:"timerID" binding:"required"`
}

type GetTasksResp struct {
	CodeMsg
	Total int64   `json:"total"`
	Data  []*Task `json:"data"`
}

func NewGetTasksResp(tasks []*Task, total int64, codeMsg CodeMsg) *GetTasksResp {
	return &GetTasksResp{
		CodeMsg: codeMsg,
		Total:   total,
		Data:    tasks,
	}
}

// Task 运行流水记录
type Task struct {
	ID       uint      `json:"id"`       // 任务 ID
	App      string    `json:"app"`      // 定义ID
	TimerID  uint      `json:"timerID"`  // 定义ID
	Output   string    `json:"output"`   // 执行结果
	RunTimer time.Time `json:"runTimer"` // 执行时间
	CostTime int       `json:"costTime"` // 执行耗时
	Status   int       `json:"status"`   // 当前状态
}

func NewTask(task *po.Task) *Task {
	return &Task{
		ID:       task.ID,
		App:      task.App,
		TimerID:  task.TimerID,
		Output:   task.Output,
		RunTimer: task.RunTimer,
		CostTime: task.CostTime,
		Status:   task.Status,
	}
}

func NewTasks(tasks []*po.Task) []*Task {
	vTasks := make([]*Task, 0, len(tasks))
	for _, task := range tasks {
		vTasks = append(vTasks, NewTask(task))
	}
	return vTasks
}

func (t *Task) ToPO() *po.Task {
	return &po.Task{
		App:      t.App,
		TimerID:  t.TimerID,
		Output:   t.Output,
		RunTimer: t.RunTimer,
		CostTime: t.CostTime,
		Status:   t.Status,
	}
}

package webserver

import (
	"context"

	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	dao "github.com/xiaoxuxiansheng/xtimer/dao/task"
)

type TaskService struct {
	dao *dao.TaskDAO
}

func NewTaskService(dao *dao.TaskDAO) *TaskService {
	return &TaskService{dao: dao}
}

func (t *TaskService) GetTask(ctx context.Context, id uint) (*vo.Task, error) {
	task, err := t.dao.GetTask(ctx, dao.WithTaskID(id))
	if err != nil {
		return nil, err
	}
	return vo.NewTask(task), nil
}

func (t *TaskService) GetTasks(ctx context.Context, req *vo.GetTasksReq) ([]*vo.Task, int64, error) {
	total, err := t.dao.Count(ctx, dao.WithTimerID(req.TimerID))
	if err != nil {
		return nil, -1, err
	}

	offset, limit := req.Get()
	if total <= int64(offset) {
		return []*vo.Task{}, total, nil
	}
	tasks, err := t.dao.GetTasks(ctx, dao.WithTimerID(req.TimerID), dao.WithPageLimit(offset, limit), dao.WithDesc())
	if err != nil {
		return nil, -1, err
	}

	return vo.NewTasks(tasks), total, nil
}

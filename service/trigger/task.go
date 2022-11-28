package trigger

import (
	"context"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	dao "github.com/xiaoxuxiansheng/xtimer/dao/task"
)

type TaskService struct {
	cache *dao.TaskCache
	dao   taskDAO
}

func NewTaskService(dao *dao.TaskDAO, cache *dao.TaskCache) *TaskService {
	return &TaskService{
		dao:   dao,
		cache: cache,
	}
}

func (t *TaskService) GetTasksByTime(ctx context.Context, key string, start, end time.Time) ([]*vo.Task, error) {
	// 先走缓存
	if tasks, err := t.cache.GetTasksByTime(ctx, key, start.Unix(), end.Unix()); err == nil {
		return vo.NewTasks(tasks), nil
	}

	// 倘若缓存 miss 再走 db
	tasks, err := t.dao.GetTasks(ctx, dao.WithStartTime(start), dao.WithEndTime(end), dao.WithStatus(int32(consts.NotRunned.ToInt())))
	if err != nil {
		return nil, err
	}
	return vo.NewTasks(tasks), nil
}

type taskDAO interface {
	GetTasks(ctx context.Context, opts ...dao.Option) ([]*po.Task, error)
}

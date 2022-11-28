package task

import (
	"context"

	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/pkg/mysql"
)

type TaskDAO struct {
	client *mysql.Client
}

func NewTaskDAO(client *mysql.Client) *TaskDAO {
	return &TaskDAO{
		client: client,
	}
}

func (t *TaskDAO) GetTask(ctx context.Context, opts ...Option) (*po.Task, error) {
	db := t.client.DB.WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}

	var task po.Task
	return &task, db.First(&task).Error
}

func (t *TaskDAO) GetTasks(ctx context.Context, opts ...Option) ([]*po.Task, error) {
	db := t.client.DB.WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}

	var tasks []*po.Task
	return tasks, db.Model(&po.Task{}).Scan(&tasks).Error
}

func (t *TaskDAO) UpdateTask(ctx context.Context, task *po.Task) error {
	return t.client.DB.WithContext(ctx).Updates(task).Error
}

func (t *TaskDAO) Count(ctx context.Context, opts ...Option) (int64, error) {
	db := t.client.DB.WithContext(ctx).Model(&po.Task{})
	for _, opt := range opts {
		db = opt(db)
	}
	var cnt int64
	return cnt, db.Count(&cnt).Error
}

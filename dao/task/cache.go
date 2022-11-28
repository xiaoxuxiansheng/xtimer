package task

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
)

type TaskCache struct {
	client       cacheClient
	confProvider *conf.SliceConfProvider
}

func NewTaskCache(client *redis.Client, confProvider *conf.SliceConfProvider) *TaskCache {
	return &TaskCache{client: client, confProvider: confProvider}
}

func (t *TaskCache) BatchCreateTasks(ctx context.Context, tasks []*po.Task) error {
	// TODO(@weixuxu): 动态分桶
	commands := make([]*redis.Command, 0, len(tasks))
	for _, task := range tasks {
		unix := task.RunTimer.Unix()
		commands = append(commands, redis.NewZAddCommand(t.GetTableName(task), unix, utils.UnionTimerIDUnix(task.TimerID, unix)))
	}

	_, err := t.client.Transaction(ctx, commands...)
	return err
}

func (t *TaskCache) GetTasksByTime(ctx context.Context, table string, start, end int64) ([]*po.Task, error) {
	timerIDUnixs, err := t.client.ZrangeByScore(ctx, table, start, end)
	if err != nil {
		return nil, err
	}

	tasks := make([]*po.Task, 0, len(timerIDUnixs))
	for _, timerIDUnix := range timerIDUnixs {
		timerID, unix, _ := utils.SplitTimerIDUnix(timerIDUnix)
		tasks = append(tasks, &po.Task{
			TimerID:  timerID,
			RunTimer: time.Unix(unix, 0),
		})
	}

	return tasks, nil
}

func (t *TaskCache) GetTableName(task *po.Task) string {
	bucket := task.TimerID % uint(t.confProvider.Get().BucketsNum)
	return fmt.Sprintf("%s_%d", task.RunTimer.Format(consts.MinuteFormat), bucket)
}

type cacheClient interface {
	Transaction(ctx context.Context, commands ...*redis.Command) ([]interface{}, error)
	ZrangeByScore(ctx context.Context, table string, score1, score2 int64) ([]string, error)
}

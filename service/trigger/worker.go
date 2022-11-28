package trigger

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/pool"
	"github.com/xiaoxuxiansheng/xtimer/service/executor"
)

type Worker struct {
	sync.Once
	task         taskService
	confProvider confProvider
	pool         pool.WorkerPool
	executor     *executor.Worker
}

func NewWorker(executor *executor.Worker, task *TaskService, pool *pool.GoWorkerPool, confProvider *conf.TriggerAppConfProvider) *Worker {
	return &Worker{
		executor:     executor,
		task:         task,
		pool:         pool,
		confProvider: confProvider,
	}
}

func (w *Worker) Work(ctx context.Context, minuteBucketKey string) error {
	// 进行为时一分钟的 zrange 处理
	conf := w.confProvider.Get()
	startTime, err := getStartMinute(minuteBucketKey)
	if err != nil {
		return err
	}

	for move := startTime; move.Before(startTime.Add(time.Minute)); move = move.Add(time.Duration(conf.ZRangeGapSeconds) * time.Second) {
		if err := w.handleBatch(ctx, minuteBucketKey, move, move.Add(time.Duration(conf.ZRangeGapSeconds)*time.Second)); err != nil {
			return err
		}
	}

	log.InfoContextf(ctx, "handle all tasks of key: %s", minuteBucketKey)
	// 任务全部执行完成，此时执行 ack
	return nil
}

func (w *Worker) handleBatch(ctx context.Context, key string, start, end time.Time) error {
	tasks, err := w.task.GetTasksByTime(ctx, key, start, end)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := w.pool.Submit(func() {
			if err := w.executor.Work(ctx, utils.UnionTimerIDUnix(task.TimerID, task.RunTimer.Unix())); err != nil {
				log.ErrorContextf(ctx, "executor work failed, err: %v", err)
			}
		}); err != nil {
			log.ErrorContextf(ctx, "handle task failed, err: %v", err)
		}
	}
	return nil
}

func getStartMinute(slice string) (time.Time, error) {
	timeBucket := strings.Split(slice, "_")
	if len(timeBucket) != 2 {
		return time.Time{}, fmt.Errorf("invalid format of msg key: %s", slice)
	}

	return utils.GetStartMinute(timeBucket[0])
}

type taskService interface {
	GetTasksByTime(ctx context.Context, key string, start, end time.Time) ([]*vo.Task, error)
}

type confProvider interface {
	Get() *conf.TriggerAppConf
}

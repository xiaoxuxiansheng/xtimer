package scheduler

import (
	"context"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
	"github.com/xiaoxuxiansheng/xtimer/service/trigger"
)

type Worker struct {
	appConfProvider  appConfProvider
	lockConfProvider *conf.LockConfProvider
	trigger          *trigger.Worker
	lockService      lockService
	stop             func()
}

func NewWorker(trigger *trigger.Worker, lockService *redis.Client, appConfProvider *conf.SchedulerAppConfProvider, lockConfProvider *conf.LockConfProvider) *Worker {
	return &Worker{
		trigger:          trigger,
		lockService:      lockService,
		lockConfProvider: lockConfProvider,
		appConfProvider:  appConfProvider,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	workerID, _ := ctx.Value(consts.WorkerIDContextKey).(int)
	lockConf := w.lockConfProvider.Get()
	ticker := time.NewTicker(time.Duration(lockConf.TryLockGapSeconds) * time.Second)
	w.stop = ticker.Stop

	for range ticker.C {
		select {
		case <-ctx.Done():
			log.WarnContextf(ctx, "worker: %d is stopped", workerID)
			return nil
		default:
		}

		w.handleSlices(ctx)
	}
	return nil
}

func (w *Worker) Stop() {
	if w.stop != nil {
		w.stop()
	}
}

func (w *Worker) handleSlices(ctx context.Context) {
	conf := w.appConfProvider.Get()
	for i := 0; i < conf.BucketsNum; i++ {
		w.handleSlice(ctx, i)
	}
}

func (w *Worker) handleSlice(ctx context.Context, bucketID int) {
	now := time.Now()
	go w.asyncHandleSlice(ctx, now.Add(-1*time.Minute), bucketID)
	go w.asyncHandleSlice(ctx, now, bucketID)
}

func (w *Worker) asyncHandleSlice(ctx context.Context, t time.Time, bucketID int) {
	lockConf := w.lockConfProvider.Get()
	locker := w.lockService.GetDistributionLock(utils.GetTimeBucketLockKey(t, bucketID))
	if err := locker.Lock(ctx, int64(lockConf.TryLockSeconds)); err != nil {
		return
	}

	if err := w.trigger.Work(ctx, utils.GetSliceMsgKey(t, bucketID)); err != nil {
		log.ErrorContextf(ctx, "trigger work failed, err: %v", err)
	}
}

type appConfProvider interface {
	Get() *conf.SchedulerAppConf
}

type lockService interface {
	GetDistributionLock(key string) redis.DistributeLocker
}

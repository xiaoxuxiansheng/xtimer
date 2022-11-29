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
	appConfProvider appConfProvider
	trigger         *trigger.Worker
	lockService     lockService
}

func NewWorker(trigger *trigger.Worker, lockService *redis.Client, appConfProvider *conf.SchedulerAppConfProvider) *Worker {
	return &Worker{
		trigger:         trigger,
		lockService:     lockService,
		appConfProvider: appConfProvider,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	workerID, _ := ctx.Value(consts.WorkerIDContextKey).(int)
	ticker := time.NewTicker(time.Duration(w.appConfProvider.Get().TryLockGapSeconds) * time.Second)
	defer ticker.Stop()

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

func (w *Worker) handleSlices(ctx context.Context) {
	for i := 0; i < w.appConfProvider.Get().BucketsNum; i++ {
		w.handleSlice(ctx, i)
	}
}

func (w *Worker) handleSlice(ctx context.Context, bucketID int) {
	now := time.Now()
	go w.asyncHandleSlice(ctx, now.Add(-time.Minute), bucketID)
	go w.asyncHandleSlice(ctx, now, bucketID)
}

func (w *Worker) asyncHandleSlice(ctx context.Context, t time.Time, bucketID int) {
	locker := w.lockService.GetDistributionLock(utils.GetTimeBucketLockKey(t, bucketID))
	if err := locker.Lock(ctx, int64(w.appConfProvider.Get().TryLockSeconds)); err != nil {
		return
	}

	workerID, _ := ctx.Value(consts.WorkerIDContextKey).(int)
	log.InfoContextf(ctx, "get scheduler lock success, key: %s, worker: %d", utils.GetTimeBucketLockKey(t, bucketID), workerID)

	ack := func() {
		if err := locker.ExpireLock(ctx, int64(w.appConfProvider.Get().SuccessExpireSeconds)); err != nil {
			log.ErrorContextf(ctx, "expire lock failed, lock key: %s, err: %v", utils.GetTimeBucketLockKey(t, bucketID), err)
		}
	}

	if err := w.trigger.Work(ctx, utils.GetSliceMsgKey(t, bucketID), ack); err != nil {
		log.ErrorContextf(ctx, "trigger work failed, err: %v", err)
	}
}

type appConfProvider interface {
	Get() *conf.SchedulerAppConf
}

type lockService interface {
	GetDistributionLock(key string) redis.DistributeLocker
}

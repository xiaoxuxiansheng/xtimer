package scheduler

import (
	"context"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/pool"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
	"github.com/xiaoxuxiansheng/xtimer/service/trigger"
)

type Worker struct {
	pool            pool.WorkerPool
	appConfProvider appConfProvider
	trigger         *trigger.Worker
	lockService     lockService
}

func NewWorker(trigger *trigger.Worker, lockService *redis.Client, appConfProvider *conf.SchedulerAppConfProvider) *Worker {
	return &Worker{
		pool:            pool.NewGoWorkerPool(appConfProvider.Get().WorkersNum),
		trigger:         trigger,
		lockService:     lockService,
		appConfProvider: appConfProvider,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(w.appConfProvider.Get().TryLockGapSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.InfoContextf(ctx, "ticking: %v", time.Now())
		select {
		case <-ctx.Done():
			log.WarnContext(ctx, "stopped")
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
	if err := w.pool.Submit(func() {
		w.asyncHandleSlice(ctx, now.Add(-time.Minute), bucketID)
	}); err != nil {
		log.ErrorContextf(ctx, "[handle slice] submit task failed, err: %v", err)
	}
	if err := w.pool.Submit(func() {
		w.asyncHandleSlice(ctx, now, bucketID)
	}); err != nil {
		log.ErrorContextf(ctx, "[handle slice] submit task failed, err: %v", err)
	}
}

func (w *Worker) asyncHandleSlice(ctx context.Context, t time.Time, bucketID int) {
	locker := w.lockService.GetDistributionLock(utils.GetTimeBucketLockKey(t, bucketID))
	if err := locker.Lock(ctx, int64(w.appConfProvider.Get().TryLockSeconds)); err != nil {
		log.WarnContextf(ctx, "get lock failed, err: %v, key: %s", err, utils.GetTimeBucketLockKey(t, bucketID))
		return
	}

	log.InfoContextf(ctx, "get scheduler lock success, key: %s", utils.GetTimeBucketLockKey(t, bucketID))

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

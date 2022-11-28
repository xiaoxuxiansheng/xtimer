package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/pool"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
	"github.com/xiaoxuxiansheng/xtimer/service/trigger"
)

type Worker struct {
	sync.Once
	appConfProvider  appConfProvider
	lockConfProvider *conf.LockConfProvider
	trigger          *trigger.Worker
	lockService      lockService
	pool             pool.WorkerPool
	stop             func()
}

func NewWorker(trigger *trigger.Worker, lockService *redis.Client, appConfProvider *conf.WorkerAppConfProvider, lockConfProvider *conf.LockConfProvider, pool *pool.GoWorkerPool) *Worker {
	return &Worker{
		trigger:          trigger,
		lockService:      lockService,
		lockConfProvider: lockConfProvider,
		appConfProvider:  appConfProvider,
		pool:             pool,
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

		w.tryProduceSlices(ctx)
	}
	return nil
}

func (w *Worker) Stop() {
	if w.stop != nil {
		w.stop()
	}
}

func (w *Worker) tryProduceSlices(ctx context.Context) {
	conf := w.appConfProvider.Get()
	for i := 0; i < conf.BucketsNum; i++ {
		if err := w.tryProduceSlice(ctx, i); err != nil {
			return
		}
	}
}

func (w *Worker) tryProduceSlice(ctx context.Context, bucketID int) error {
	lockConf := w.lockConfProvider.Get()
	now := time.Now()
	locker := w.lockService.GetDistributionLock(utils.GetTimeBucketLockKey(now, bucketID))
	if err := locker.Lock(ctx, int64(lockConf.TryLockSeconds)); err != nil {
		return nil
	}

	key := utils.GetSliceMsgKey(now, bucketID)
	if err := w.produceSlice(ctx, key); err != nil {
		log.ErrorContextf(ctx, "produce slice token failed, key: %s", key)
		return err
	}
	log.InfoContextf(ctx, "produce slice token successed, key: %s", key)

	return locker.ExpireLock(ctx, int64(lockConf.SuccessExpireSeconds))
}

func (w *Worker) produceSlice(ctx context.Context, key string) error {
	return w.pool.Submit(func() {
		if err := w.trigger.Work(ctx, key); err != nil {
			log.ErrorContextf(ctx, "trigger work failed, err: %v", err)
		}
	})
}

type appConfProvider interface {
	Get() *conf.WorkerAppConf
}

type lockService interface {
	GetDistributionLock(key string) redis.DistributeLocker
}

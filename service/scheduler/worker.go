package scheduler

import (
	"context"
	"strconv"
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
	pool            pool.WorkerPool
	appConfProvider appConfProvider
	trigger         *trigger.Worker
	lockService     lockService
	bucketGetter    bucketGetter
	minuteBuckets   map[string]int
}

func NewWorker(trigger *trigger.Worker, redisClient *redis.Client, appConfProvider *conf.SchedulerAppConfProvider) *Worker {
	return &Worker{
		pool:            pool.NewGoWorkerPool(appConfProvider.Get().WorkersNum),
		trigger:         trigger,
		lockService:     redisClient,
		bucketGetter:    redisClient,
		appConfProvider: appConfProvider,
		minuteBuckets:   make(map[string]int),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.trigger.Start(ctx)

	ticker := time.NewTicker(time.Duration(w.appConfProvider.Get().TryLockGapSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
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
	for i := 0; i < w.getValidBucket(ctx); i++ {
		w.handleSlice(ctx, i)
	}
}

func (w *Worker) getValidBucket(ctx context.Context) int {
	now := time.Now()
	// 删除上一分钟的数据
	delete(w.minuteBuckets, now.Add(-time.Minute).Format(consts.MinuteFormat))

	// 复用一分钟内的数据
	bucket, ok := w.minuteBuckets[now.Format(consts.MinuteFormat)]
	if ok {
		return bucket
	}

	// 更新入map进行复用
	defer func() {
		w.minuteBuckets[now.Format(consts.MinuteFormat)] = bucket
	}()

	bucket = w.appConfProvider.Get().BucketsNum
	bucketKey := utils.GetBucketCntKey(now.Format(consts.MinuteFormat))
	res, err := w.bucketGetter.Get(ctx, bucketKey)
	if err != nil {
		log.ErrorContextf(ctx, "[scheduler] get bucket failed, key: %s, err:%v", bucketKey, err)
		return bucket
	}

	_bucket, err := strconv.Atoi(res)
	if err != nil {
		log.ErrorContextf(ctx, "[scheduler] get invalid bucket, key: %s, got:%v", bucketKey, res)
		return bucket
	}

	bucket = _bucket
	log.InfoContextf(ctx, "[scheduler] get valid bucket success, bucket: %d, cur: %v", _bucket, time.Now())

	return bucket
}

func (w *Worker) handleSlice(ctx context.Context, bucketID int) {
	// log.InfoContextf(ctx, "scheduler_1 start: %v", time.Now())
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
	// log.InfoContextf(ctx, "scheduler_1 end: %v", time.Now())
}

func (w *Worker) asyncHandleSlice(ctx context.Context, t time.Time, bucketID int) {
	// log.InfoContextf(ctx, "scheduler_2 start: %v", time.Now())
	// defer func() {
	// 	log.InfoContextf(ctx, "scheduler_2 end: %v", time.Now())
	// }()

	locker := w.lockService.GetDistributionLock(utils.GetTimeBucketLockKey(t, bucketID))
	if err := locker.Lock(ctx, int64(w.appConfProvider.Get().TryLockSeconds)); err != nil {
		// log.WarnContextf(ctx, "get lock failed, err: %v, key: %s", err, utils.GetTimeBucketLockKey(t, bucketID))
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

type bucketGetter interface {
	Get(ctx context.Context, key string) (string, error)
}

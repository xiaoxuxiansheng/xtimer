package service

import (
	"context"
	"sync"
	"time"

	mconf "github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	taskdao "github.com/xiaoxuxiansheng/xtimer/dao/task"
	timerdao "github.com/xiaoxuxiansheng/xtimer/dao/timer"
	"github.com/xiaoxuxiansheng/xtimer/pkg/cron"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
)

type Worker struct {
	timerDAO          *timerdao.TimerDAO
	taskDAO           *taskdao.TaskDAO
	taskCache         *taskdao.TaskCache
	cronParser        *cron.CronParser
	lockService       *redis.Client
	appConfigProvider *mconf.MigratorAppConfProvider
}

func NewWorker(timerDAO *timerdao.TimerDAO, taskDAO *taskdao.TaskDAO, taskCache *taskdao.TaskCache, lockService *redis.Client,
	cronParser *cron.CronParser, appConfigProvider *mconf.MigratorAppConfProvider) *Worker {
	return &Worker{
		timerDAO:          timerDAO,
		taskDAO:           taskDAO,
		taskCache:         taskCache,
		lockService:       lockService,
		cronParser:        cronParser,
		appConfigProvider: appConfigProvider,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	conf := w.appConfigProvider.Get()
	ticker := time.NewTicker(time.Duration(conf.MigrateStepMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		locker := w.lockService.GetDistributionLock(utils.GetMigratorLockKey(utils.GetStartHour(time.Now())))
		if err := locker.Lock(ctx, int64(conf.MigrateTryLockMinutes)*int64(time.Minute/time.Second)); err != nil {
			continue
		}

		if err := w.migrate(ctx); err != nil {
			continue
		}

		_ = locker.ExpireLock(ctx, int64(conf.MigrateSucessExpireMinutes)*int64(time.Minute/time.Second))
	}
	return nil
}

func (w *Worker) migrate(ctx context.Context) error {
	// 获取下一个小时的时间范围
	timers, err := w.timerDAO.GetTimers(ctx, timerdao.WithStatus(int32(consts.Enabled.ToInt())))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	now := time.Now()
	start, end := utils.GetStartHour(now.Add(time.Hour)), utils.GetStartHour(now.Add(2*time.Hour))
	for _, timer := range timers {
		wg.Add(1)
		go func(timer *po.Timer) {
			defer wg.Done()
			nexts, _ := w.cronParser.NextsBetween(timer.Cron, start, end)
			if err := w.timerDAO.BatchCreateRecords(ctx, timer.BatchTasksFromTimer(nexts)); err != nil {
				log.ErrorContextf(ctx, "batch create records for timer: %d failed, err: %v", timer.ID, err)
			}
		}(timer)
	}

	wg.Wait()
	// 迁移完成后，将所有添加的 task 取出，添加到 redis 当中
	tasks, err := w.taskDAO.GetTasks(ctx, taskdao.WithStartTime(start), taskdao.WithEndTime(end))
	if err != nil {
		return err
	}

	return w.taskCache.BatchCreateTasks(ctx, tasks)
}

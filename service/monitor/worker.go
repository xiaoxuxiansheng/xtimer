package monitor

import (
	"context"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	taskdao "github.com/xiaoxuxiansheng/xtimer/dao/task"
	timerdao "github.com/xiaoxuxiansheng/xtimer/dao/timer"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/promethus"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
)

type Worker struct {
	lockService *redis.Client
	taskDAO     *taskdao.TaskDAO
	timerDAO    *timerdao.TimerDAO
	reporter    *promethus.Reporter
}

func NewWorker(taskDAO *taskdao.TaskDAO, timerDAO *timerdao.TimerDAO, lockService *redis.Client, reporter *promethus.Reporter) *Worker {
	return &Worker{
		taskDAO:     taskDAO,
		timerDAO:    timerDAO,
		lockService: lockService,
		reporter:    reporter,
	}
}

// 每隔1分钟上报失败的定时器数量
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
		}

		now := time.Now()
		lock := w.lockService.GetDistributionLock(utils.GetMonitorLockKey(now))
		if err := lock.Lock(ctx, 2*int64(time.Minute/time.Second)); err != nil {
			continue
		}

		// 取上一分钟的定时器进行查询
		minute := utils.GetMinute(now)
		go w.reportUnexecedTasksCnt(ctx, minute)
		go w.reportEnabledTimersCnt(ctx)
	}
}

func (w *Worker) reportUnexecedTasksCnt(ctx context.Context, minute time.Time) {
	unexecedTasksCnt, err := w.taskDAO.Count(ctx, taskdao.WithStartTime(minute.Add(-time.Minute)), taskdao.WithEndTime(minute), taskdao.WithStatus(int32(consts.NotRunned)))
	if err != nil {
		log.ErrorContextf(ctx, "[monitor] get unexeced tasks cnt failed, err: %v", err)
		return
	}
	w.reporter.ReportTimerUnexecedRecord(float64(unexecedTasksCnt))
	log.InfoContextf(ctx, "[monitor] report unexeced tasks cnt success, cnt: %d", unexecedTasksCnt)
}

func (w *Worker) reportEnabledTimersCnt(ctx context.Context) {
	enabledTimerCnt, err := w.timerDAO.Count(ctx, timerdao.WithStatus(int32(consts.Enabled)))
	if err != nil {
		log.ErrorContextf(ctx, "[monitor] get enabled timer cnt failed, err: %v", err)
		return
	}
	w.reporter.ReportTimerEnabledRecord(float64(enabledTimerCnt))
	log.InfoContextf(ctx, "[monitor] report enabled timer cnt success, cnt: %d", enabledTimerCnt)
}

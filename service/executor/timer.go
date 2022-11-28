package executor

import (
	"context"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	taskdao "github.com/xiaoxuxiansheng/xtimer/dao/task"
	timerdao "github.com/xiaoxuxiansheng/xtimer/dao/timer"
)

type TimerService struct {
	confProvider *conf.MigratorAppConfProvider
	ctx          context.Context
	stop         func()
	timers       map[uint]*vo.Timer
	timerDAO     timerDAO
	taskDAO      *taskdao.TaskDAO
}

func NewTimerService(timerDAO *timerdao.TimerDAO, taskDAO *taskdao.TaskDAO, confProvider *conf.MigratorAppConfProvider) *TimerService {
	return &TimerService{
		confProvider: confProvider,
		timers:       make(map[uint]*vo.Timer),
		timerDAO:     timerDAO,
		taskDAO:      taskDAO,
	}
}

func (t *TimerService) Start(ctx context.Context) {
	t.ctx, t.stop = context.WithCancel(ctx)

	stepMinutes := t.confProvider.Get().TimerDetailCacheMinutes
	ticker := time.NewTicker(time.Duration(stepMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		start := time.Now()
		t.timers, _ = t.getTimersByTime(ctx, start, start.Add(time.Duration(stepMinutes)*time.Minute))
	}
}

func (t *TimerService) getTimersByTime(ctx context.Context, start, end time.Time) (map[uint]*vo.Timer, error) {
	tasks, err := t.taskDAO.GetTasks(ctx, taskdao.WithStartTime(start), taskdao.WithEndTime(end))
	if err != nil {
		return nil, err
	}

	timerIDs := getTimerIDs(tasks)
	if len(timerIDs) == 0 {
		return nil, nil
	}
	pTimers, err := t.timerDAO.GetTimers(ctx, timerdao.WithIDs(timerIDs), timerdao.WithStatus(int32(consts.Enabled)))
	if err != nil {
		return nil, err
	}

	return getTimersMap(pTimers)
}

func getTimerIDs(tasks []*po.Task) []uint {
	timerIDSet := make(map[uint]struct{})
	for _, task := range tasks {
		if _, ok := timerIDSet[task.TimerID]; ok {
			continue
		}
		timerIDSet[task.TimerID] = struct{}{}
	}
	timerIDs := make([]uint, 0, len(timerIDSet))
	for id := range timerIDSet {
		timerIDs = append(timerIDs, id)
	}
	return timerIDs
}

func getTimersMap(pTimers []*po.Timer) (map[uint]*vo.Timer, error) {
	vTimers, err := vo.NewTimers(pTimers)
	if err != nil {
		return nil, err
	}

	timers := make(map[uint]*vo.Timer, len(vTimers))
	for _, vTimer := range vTimers {
		timers[vTimer.ID] = vTimer
	}
	return timers, nil
}

func (t *TimerService) GetTimer(ctx context.Context, id uint) (*vo.Timer, error) {
	if vTimer, ok := t.timers[id]; ok {
		return vTimer, nil
	}

	timer, err := t.timerDAO.GetTimer(ctx, timerdao.WithID(id))
	if err != nil {
		return nil, err
	}

	return vo.NewTimer(timer)
}

func (t *TimerService) Stop() {
	t.stop()
}

type timerDAO interface {
	GetTimer(context.Context, ...timerdao.Option) (*po.Timer, error)
	GetTimers(ctx context.Context, opts ...timerdao.Option) ([]*po.Timer, error)
}

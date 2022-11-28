package webserver

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	taskdao "github.com/xiaoxuxiansheng/xtimer/dao/task"
	timerdao "github.com/xiaoxuxiansheng/xtimer/dao/timer"
	"github.com/xiaoxuxiansheng/xtimer/pkg/cron"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/mysql"
)

type TimerService struct {
	dao                 timerDAO
	confProvider        confProvider
	migrateConfProvider *conf.MigratorAppConfProvider
	cronParser          cronParser
	taskCache           taskCache
}

func NewTimerService(dao *timerdao.TimerDAO, taskCache *taskdao.TaskCache,
	confProvider *conf.WebServerAppConfProvider, migrateConfProvider *conf.MigratorAppConfProvider, parser *cron.CronParser) *TimerService {
	return &TimerService{
		dao:                 dao,
		confProvider:        confProvider,
		migrateConfProvider: migrateConfProvider,
		taskCache:           taskCache,
		cronParser:          parser,
	}
}

func (t *TimerService) CreateTimer(ctx context.Context, timer *vo.Timer) (uint, error) {
	pTimer, err := timer.ToPO()
	if err != nil {
		return 0, err
	}
	return t.dao.CreateTimer(ctx, pTimer)
}

func (t *TimerService) DeleteTimer(ctx context.Context, id uint) error {
	return t.dao.DeleteTimer(ctx, id)
}

func (t *TimerService) UpdateTimer(ctx context.Context, timer *vo.Timer) error {
	pTimer, err := timer.ToPO()
	if err != nil {
		return err
	}
	return t.dao.UpdateTimer(ctx, pTimer)
}

func (t *TimerService) GetTimer(ctx context.Context, id uint) (*vo.Timer, error) {
	pTimer, err := t.dao.GetTimer(ctx, timerdao.WithID(id))
	if err != nil {
		return nil, err
	}

	return vo.NewTimer(pTimer)
}

func (t *TimerService) EnableTimer(ctx context.Context, id uint) error {
	do := func(ctx context.Context, dao *timerdao.TimerDAO, timer *po.Timer) error {
		// 状态校验
		if timer.Status != consts.Unabled.ToInt() {
			return fmt.Errorf("not unabled status, enable failed, timer id: %d", id)
		}

		// 取得批量的执行时机
		// end 为下两个切片的右边界
		executeTimes, err := t.cronParser.NextsBefore(timer.Cron,
			utils.GetForwardTwoMigrateStepEnd(time.Now(), 3*time.Duration(t.migrateConfProvider.Get().MigrateStepMinutes)*time.Minute))
		if err != nil {
			log.ErrorContextf(ctx, "get executeTimes failed, err: %v", err)
			return err
		}

		// 执行时机加入数据库
		tasks := timer.BatchTasksFromTimer(executeTimes)
		// 基于 timer_id + run_timer 唯一键，保证任务不被重复插入
		if err := dao.BatchCreateRecords(ctx, tasks); err != nil && !mysql.IsDuplicateEntryErr(err) {
			return err
		}

		// 执行时机加入 redis 跳表
		if err := t.taskCache.BatchCreateTasks(ctx, tasks); err != nil {
			return err
		}

		// 修改 timer 状态为激活态
		timer.Status = consts.Enabled.ToInt()
		return dao.UpdateTimer(ctx, timer)
	}

	return t.dao.DoWithLock(ctx, id, do)
}

func (t *TimerService) UnableTimer(ctx context.Context, id uint) error {
	do := func(ctx context.Context, dao *timerdao.TimerDAO, timer *po.Timer) error {
		// 状态校验
		if timer.Status != consts.Enabled.ToInt() {
			return fmt.Errorf("not enabled status, unable failed, timer id: %d", id)
		}

		// 修改 timer 状态为激活态
		timer.Status = consts.Unabled.ToInt()
		return dao.UpdateTimer(ctx, timer)
	}

	return t.dao.DoWithLock(ctx, id, do)
}

func (t *TimerService) GetAppTimers(ctx context.Context, req *vo.GetAppTimersReq) ([]*vo.Timer, int64, error) {
	total, err := t.dao.Count(ctx, timerdao.WithApp(req.App))
	if err != nil {
		return nil, -1, err
	}

	offset, limit := req.Get()
	if total <= int64(offset) {
		return []*vo.Timer{}, total, nil
	}

	timers, err := t.dao.GetTimers(ctx, timerdao.WithApp(req.App), timerdao.WithPageLimit(offset, limit), timerdao.WithDesc())
	if err != nil {
		return nil, -1, err
	}

	vTimers, err := vo.NewTimers(timers)
	return vTimers, total, err
}

type timerDAO interface {
	CreateTimer(ctx context.Context, timer *po.Timer) (uint, error)
	DeleteTimer(ctx context.Context, id uint) error
	UpdateTimer(ctx context.Context, timer *po.Timer) error
	GetTimer(ctx context.Context, opts ...timerdao.Option) (*po.Timer, error)
	BatchCreateRecords(ctx context.Context, tasks []*po.Task) error
	DoWithLock(ctx context.Context, id uint, do func(ctx context.Context, dao *timerdao.TimerDAO, timer *po.Timer) error) error
	GetTimers(ctx context.Context, opts ...timerdao.Option) ([]*po.Timer, error)
	Count(ctx context.Context, opts ...timerdao.Option) (int64, error)
}

type confProvider interface {
	Get() *conf.WebServerAppConf
}

type taskCache interface {
	BatchCreateTasks(ctx context.Context, tasks []*po.Task) error
}

type cronParser interface {
	NextsBefore(cron string, end time.Time) ([]time.Time, error)
}

package scheduler

import (
	"context"
	"sync"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	service "github.com/xiaoxuxiansheng/xtimer/service/scheduler"
)

// 读取配置启动多个协程进行
type WorkerApp struct {
	sync.Once
	confProvider confProvider
	service      workerService
	ctx          context.Context
	stop         func()
}

func NewWorkerApp(service *service.Worker, confProvider *conf.SchedulerAppConfProvider) *WorkerApp {
	w := WorkerApp{
		service:      service,
		confProvider: confProvider,
	}

	w.ctx, w.stop = context.WithCancel(context.Background())
	return &w
}

func (w *WorkerApp) Start() {
	w.Do(w.start)
}

func (w *WorkerApp) start() {
	log.InfoContext(w.ctx, "worker app is starting")
	for i := 0; i < w.confProvider.Get().WorkersNum; i++ {
		i := i
		go func() {
			ctx := context.WithValue(w.ctx, consts.WorkerIDContextKey, i)
			if err := w.service.Start(ctx); err != nil {
				log.ErrorContextf(ctx, "worker start failed, err: %v", err)
			}
		}()
	}
}

func (w *WorkerApp) Stop() {
	w.stop()
	log.WarnContext(w.ctx, "worker app is stopped")
}

type workerService interface {
	Start(context.Context) error
}

type confProvider interface {
	Get() *conf.SchedulerAppConf
}

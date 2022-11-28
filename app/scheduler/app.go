package scheduler

import (
	"context"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/pool"
	service "github.com/xiaoxuxiansheng/xtimer/service/scheduler"
)

// 读取配置启动多个协程进行
type WorkerApp struct {
	pool         pool.WorkerPool
	confProvider confProvider
	service      workerService
	ctx          context.Context
	stop         func()
}

// var (
// 	app *WorkerApp
// )

// // func GetWorkerApp() (*WorkerApp, error) {
// // 	return app, container.Invoke(func(_app *WorkerApp) {
// // 		app = _app
// // 	})
// // }

func NewWorkerApp(service *service.Worker, pool *pool.GoWorkerPool, confProvider *conf.SchedulerAppConfProvider) *WorkerApp {
	w := WorkerApp{
		service:      service,
		pool:         pool,
		confProvider: confProvider,
	}

	w.ctx, w.stop = context.WithCancel(context.Background())
	return &w
}

func (w *WorkerApp) Start() {
	log.InfoContext(w.ctx, "worker app is starting")
	for i := 0; i < w.confProvider.Get().WorkersNum; i++ {
		i := i
		if err := w.pool.Submit(func() {
			ctx := context.WithValue(w.ctx, consts.WorkerIDContextKey, i)
			if err := w.service.Start(ctx); err != nil {
				log.ErrorContextf(ctx, "worker start failed, err: %v", err)
			}
		}); err != nil {
			log.ErrorContextf(w.ctx, "worker start task submit to pool failed, err: %v", err)
		}
	}
}

func (w *WorkerApp) Stop() {
	w.stop()
	w.service.Stop()
	log.WarnContext(w.ctx, "worker app is stopped")
}

type workerService interface {
	Start(context.Context) error
	Stop()
}

type confProvider interface {
	Get() *conf.SchedulerAppConf
}

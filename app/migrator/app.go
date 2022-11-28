package migrator

import (
	"context"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	service "github.com/xiaoxuxiansheng/xtimer/service/migrator"
)

// var (
// 	app *MigratorApp
// )

// func GetWorkerApp() (*MigratorApp, error) {
// 	return app, container.Invoke(func(_app *MigratorApp) {
// 		app = _app
// 	})
// }

// 定期从 timer 表中加载一系列 task 记录添加到 task 表中
// 并且将 一系列 task 添加到 redis zset 当中
type MigratorApp struct {
	ctx            context.Context
	stop           func()
	worker         *service.Worker
	configProvider *conf.MigratorAppConfProvider
}

func NewMigratorApp(worker *service.Worker, configProvider *conf.MigratorAppConfProvider) *MigratorApp {
	m := MigratorApp{
		worker:         worker,
		configProvider: configProvider,
	}

	m.ctx, m.stop = context.WithCancel(context.Background())
	return &m
}

func (m *MigratorApp) Start() {
	for i := 0; i < m.configProvider.Get().WorkersNum; i++ {
		go func() {
			if err := m.worker.Start(m.ctx); err != nil {
				log.ErrorContextf(m.ctx, "start worker failed, err: %v", err)
			}
		}()
	}
}

func (m *MigratorApp) Stop() {
	m.stop()
}

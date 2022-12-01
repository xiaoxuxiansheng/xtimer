package migrator

import (
	"context"
	"sync"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	service "github.com/xiaoxuxiansheng/xtimer/service/migrator"
)

// 定期从 timer 表中加载一系列 task 记录添加到 task 表中
type MigratorApp struct {
	sync.Once
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
	m.Do(func() {
		log.InfoContext(m.ctx, "migrator is starting")
		go func() {
			if err := m.worker.Start(m.ctx); err != nil {
				log.ErrorContextf(m.ctx, "start worker failed, err: %v", err)
			}
		}()
	})
}

func (m *MigratorApp) Stop() {
	m.stop()
}

package pool

import (
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"

	"github.com/panjf2000/ants/v2"
)

// WorkerPool 协程工作池.
type WorkerPool interface {
	Submit(func()) error
}

// GoWorkerPool golang 协程工作池.
type GoWorkerPool struct {
	pool         *ants.Pool
	confProvider confProvider
}

// Submit 提交任务.
func (g *GoWorkerPool) Submit(f func()) error {
	return g.pool.Submit(f)
}

func NewGoWorkerPool(confProvider confProvider) (*GoWorkerPool, error) {
	p := GoWorkerPool{
		confProvider: confProvider,
	}

	conf := p.confProvider.Get()

	pool, err := ants.NewPool(
		conf.Size,
		ants.WithExpiryDuration(time.Duration(conf.ExpireSeconds)*time.Second),
	)
	if err != nil {
		return nil, err
	}
	p.pool = pool
	return &p, nil
}

type confProvider interface {
	Get() *conf.WorkerPoolConf
}

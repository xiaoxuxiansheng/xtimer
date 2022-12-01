package pool

import (
	"time"

	"github.com/panjf2000/ants/v2"
)

// WorkerPool 协程工作池.
type WorkerPool interface {
	Submit(func()) error
}

// GoWorkerPool golang 协程工作池.
type GoWorkerPool struct {
	pool *ants.Pool
}

// Submit 提交任务.
func (g *GoWorkerPool) Submit(f func()) error {
	return g.pool.Submit(f)
}

func NewGoWorkerPool(size int) *GoWorkerPool {
	pool, err := ants.NewPool(
		size,
		ants.WithExpiryDuration(time.Minute),
	)
	if err != nil {
		panic(err)
	}
	return &GoWorkerPool{pool: pool}
}

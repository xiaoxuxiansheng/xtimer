package conf

type WorkerPoolConf struct {
	Size          int `yaml:"size"`
	ExpireSeconds int `yaml:"expireSeconds"`
}

type WorkerPoolConfProvider struct {
	conf *WorkerPoolConf
}

func NewWorkerPoolConfProvider(conf *WorkerPoolConf) *WorkerPoolConfProvider {
	return &WorkerPoolConfProvider{
		conf: conf,
	}
}

func (w *WorkerPoolConfProvider) Get() *WorkerPoolConf {
	return w.conf
}

var defaultWorkerPoolConfProvider *WorkerPoolConfProvider

func DefaultWorkerPoolConfProvider() *WorkerPoolConfProvider {
	return defaultWorkerPoolConfProvider
}

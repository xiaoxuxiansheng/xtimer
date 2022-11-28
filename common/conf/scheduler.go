package conf

type WorkerAppConf struct {
	WorkersNum int `yaml:"workersNum"`
	BucketsNum int `yaml:"bucketsNum"`
}

var defaultWorkerAppConfProvider *WorkerAppConfProvider

type WorkerAppConfProvider struct {
	conf *WorkerAppConf
}

func NewWorkerAppConfProvider(conf *WorkerAppConf) *WorkerAppConfProvider {
	return &WorkerAppConfProvider{conf: conf}
}

func (w *WorkerAppConfProvider) Get() *WorkerAppConf {
	return w.conf
}

func DefaultWorkerAppConfProvider() *WorkerAppConfProvider {
	return defaultWorkerAppConfProvider
}

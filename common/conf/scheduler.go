package conf

type SchedulerAppConf struct {
	SchedulersNum int `yaml:"schedulersNum"`
	WorkersNum    int `yaml:"workersNum"`
	// 在默认桶数的基础上，每多 200 个任务增加一个桶数
	BucketsNum             int `yaml:"bucketsNum"`
	TryLockSeconds         int `yaml:"tryLockSeconds"`
	TryLockGapMilliSeconds int `yaml:"tryLockGapMilliSeconds"`
	SuccessExpireSeconds   int `yaml:"successExpireSeconds"`
}

var defaultSchedulerAppConfProvider *SchedulerAppConfProvider

type SchedulerAppConfProvider struct {
	conf *SchedulerAppConf
}

func NewSchedulerAppConfProvider(conf *SchedulerAppConf) *SchedulerAppConfProvider {
	return &SchedulerAppConfProvider{conf: conf}
}

func (s *SchedulerAppConfProvider) Get() *SchedulerAppConf {
	return s.conf
}

func DefaultSchedulerAppConfProvider() *SchedulerAppConfProvider {
	return defaultSchedulerAppConfProvider
}

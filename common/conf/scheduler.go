package conf

type SchedulerAppConf struct {
	SchedulersNum        int `yaml:"schedulersNum"`
	WorkersNum           int `yaml:"workersNum"`
	BucketsNum           int `yaml:"bucketsNum"`
	TryLockSeconds       int `yaml:"tryLockSeconds"`
	TryLockGapSeconds    int `yaml:"tryLockGapSeconds"`
	SuccessExpireSeconds int `yaml:"successExpireSeconds"`
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

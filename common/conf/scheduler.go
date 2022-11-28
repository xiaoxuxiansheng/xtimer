package conf

type SchedulerAppConf struct {
	WorkersNum int `yaml:"workersNum"`
	BucketsNum int `yaml:"bucketsNum"`
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

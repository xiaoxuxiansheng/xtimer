package conf

type TriggerAppConf struct {
	ZRangeGapSeconds int `yaml:"zrangeGapSeconds"`
	WorkersNum       int `yaml:"workersNum"`
}

var defaultTriggerAppConfProvider *TriggerAppConfProvider

type TriggerAppConfProvider struct {
	conf *TriggerAppConf
}

func NewTriggerAppConfProvider(conf *TriggerAppConf) *TriggerAppConfProvider {
	return &TriggerAppConfProvider{conf: conf}
}

func (t *TriggerAppConfProvider) Get() *TriggerAppConf {
	return t.conf
}

func DefaultTriggerAppConfProvider() *TriggerAppConfProvider {
	return defaultTriggerAppConfProvider
}

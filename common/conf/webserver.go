package conf

type WebServerAppConf struct {
	Port                    int `yaml:"port"`
	MigrateTimerStepSeconds int `yaml:"migrateTimeStepSeconds"`
	BucketsNum              int `yaml:"bucketsNum"`
}

var defaultWebServerAppConfProvider *WebServerAppConfProvider

type WebServerAppConfProvider struct {
	conf *WebServerAppConf
}

func NewWebServerAppConfProvider(conf *WebServerAppConf) *WebServerAppConfProvider {
	return &WebServerAppConfProvider{conf: conf}
}

func (w *WebServerAppConfProvider) Get() *WebServerAppConf {
	return w.conf
}

func DefaultWebServerAppConfProvider() *WebServerAppConfProvider {
	return defaultWebServerAppConfProvider
}

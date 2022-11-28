package conf

// LockConfig 锁配置
type LockConfig struct {
	TryLockSeconds       int `yaml:"tryLockSeconds"`
	TryLockGapSeconds    int `yaml:"tryLockGapSeconds"`
	SuccessExpireSeconds int `yaml:"successExpireSeconds"`
}

type LockConfProvider struct {
	conf *LockConfig
}

func NewLockConfProvider(conf *LockConfig) *LockConfProvider {
	return &LockConfProvider{
		conf: conf,
	}
}

func (l *LockConfProvider) Get() *LockConfig {
	return l.conf
}

var defaultLockConfProvider *LockConfProvider

func DefaultLockConfProvider() *LockConfProvider {
	return defaultLockConfProvider
}

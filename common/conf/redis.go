package conf

// RedisConfig 缓存配置
type RedisConfig struct {
	Network            string `yaml:"network"`
	Address            string `yaml:"address"`
	Password           string `yaml:"password"`
	MaxIdle            int    `yaml:"maxIdle"`
	IdleTimeoutSeconds int    `yaml:"idleTimeout"`
	// 连接池最大存活的连接数.
	MaxActive int `yaml:"maxActive"`
	// 当连接数达到上限时，新的请求是等待还是立即报错.
	Wait bool `yaml:"wait"`
}

type RedisConfigProvider struct {
	conf *RedisConfig
}

func NewRedisConfigProvider(conf *RedisConfig) *RedisConfigProvider {
	return &RedisConfigProvider{
		conf: conf,
	}
}

func (r *RedisConfigProvider) Get() *RedisConfig {
	return r.conf
}

var defaultRedisConfProvider *RedisConfigProvider

func DefaultRedisConfigProvider() *RedisConfigProvider {
	return defaultRedisConfProvider
}

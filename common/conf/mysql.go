package conf

// MySQLConfig 数据库配置
type MySQLConfig struct {
	DSN string `yaml:"dsn"`
	// 最大连接数
	MaxOpenConns int `yaml:"maxOpenConns"`
	// 最大空闲连接数
	MaxIdleConns int `yaml:"maxIdleConns"`
}

type MysqlConfProvider struct {
	conf *MySQLConfig
}

func NewMysqlConfProvider(conf *MySQLConfig) *MysqlConfProvider {
	return &MysqlConfProvider{
		conf: conf,
	}
}

func (m *MysqlConfProvider) Get() *MySQLConfig {
	return m.conf
}

var defaultMysqlConfProvider *MysqlConfProvider

func DefaultMysqlConfProvider() *MysqlConfProvider {
	return defaultMysqlConfProvider
}

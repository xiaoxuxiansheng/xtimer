package conf

// MySQLConfig 数据库配置
type MySQLConfig struct {
	DSN string `yaml:"dsn"`
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

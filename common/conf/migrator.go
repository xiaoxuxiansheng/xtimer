package conf

type MigratorAppConf struct {
	WorkersNum                 int `yaml:"workersNum"`
	MigrateStepMinutes         int `yaml:"migrateStepMinutes"`
	MigrateSucessExpireMinutes int `yaml:"migrateSuccessExpireMinutes"`
	MigrateTryLockMinutes      int `yaml:"migrateTryLockMinutes"`
	TimerDetailCacheMinutes    int `yaml:"timerDetailCacheMinutes"`
}

var defaultMigratorAppConfProvider *MigratorAppConfProvider

type MigratorAppConfProvider struct {
	conf *MigratorAppConf
}

func NewMigratorAppConfProvider(conf *MigratorAppConf) *MigratorAppConfProvider {
	return &MigratorAppConfProvider{
		conf: conf,
	}
}

func (m *MigratorAppConfProvider) Get() *MigratorAppConf {
	return m.conf
}

func DefaultMigratorAppConfProvider() *MigratorAppConfProvider {
	return defaultMigratorAppConfProvider
}

package conf

import (
	"os"

	"github.com/spf13/viper"
)

func init() {
	// 获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := viper.New()
	config.AddConfigPath(path)   //设置读取的文件路径
	config.SetConfigName("conf") //设置读取的文件名
	config.SetConfigType("yaml") //设置文件的类型
	// 尝试进行配置读取
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := config.Unmarshal(&gConf); err != nil {
		panic(err)
	}

	defaultMigratorAppConfProvider = NewMigratorAppConfProvider(gConf.Migrator)
	defaultLockConfProvider = NewLockConfProvider(gConf.Lock)
	defaultWorkerPoolConfProvider = NewWorkerPoolConfProvider(gConf.Pool)
	defaultMysqlConfProvider = NewMysqlConfProvider(gConf.Mysql)
	defaultRedisConfProvider = NewRedisConfigProvider(gConf.Redis)
	defaultWorkerPoolConfProvider = NewWorkerPoolConfProvider(gConf.Pool)
	defaultSliceConfigProvider = NewSliceConfProvider(gConf.Slice)
	defaultTriggerAppConfProvider = NewTriggerAppConfProvider(gConf.Trigger)
	defaultSchedulerAppConfProvider = NewSchedulerAppConfProvider(gConf.Scheduler)
	defaultWebServerAppConfProvider = NewWebServerAppConfProvider(gConf.WebServer)
}

var gConf GloablConf

type GloablConf struct {
	Migrator  *MigratorAppConf  `yaml:"migrator"`
	Lock      *LockConfig       `yaml:"lock"`
	Mysql     *MySQLConfig      `yaml:"mysql"`
	Redis     *RedisConfig      `yaml:"redis"`
	Pool      *WorkerPoolConf   `yaml:"pool"`
	Slice     *SliceConf        `yaml:"slice"`
	Trigger   *TriggerAppConf   `yaml:"trigger"`
	Scheduler *SchedulerAppConf `yaml:"scheduler"`
	WebServer *WebServerAppConf `yaml:"webServer"`
}

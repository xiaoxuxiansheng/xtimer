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
	defaultMysqlConfProvider = NewMysqlConfProvider(gConf.Mysql)
	defaultRedisConfProvider = NewRedisConfigProvider(gConf.Redis)
	defaultTriggerAppConfProvider = NewTriggerAppConfProvider(gConf.Trigger)
	defaultSchedulerAppConfProvider = NewSchedulerAppConfProvider(gConf.Scheduler)
	defaultWebServerAppConfProvider = NewWebServerAppConfProvider(gConf.WebServer)
}

// 兜底配置
var gConf GloablConf = GloablConf{
	Migrator: &MigratorAppConf{
		// 单节点并行协程数
		WorkersNum: 1000,
		// 每次迁移数据的时间间隔，单位：min
		MigrateStepMinutes: 60,
		// 迁移成功更新的锁过期时间，单位：min
		MigrateSucessExpireMinutes: 120,
		// 迁移器获取锁时，初设的过期时间，单位：min
		MigrateTryLockMinutes: 20,
		// 迁移器提前将定时器数据缓存到内存中的保存时间，单位：min
		TimerDetailCacheMinutes: 2,
	},

	Scheduler: &SchedulerAppConf{
		// 单节点并行协程数
		WorkersNum: 100,
		// 分桶数量
		BucketsNum: 10,
		// 调度器获取分布式锁时初设的过期时间，单位：s
		TryLockSeconds: 70,
		// 调度器每次尝试获取分布式锁的时间间隔，单位：s
		TryLockGapSeconds: 1,
		// 时间片执行成功后，更新的分布式锁时间，单位：s
		SuccessExpireSeconds: 130,
	},

	Trigger: &TriggerAppConf{
		// 触发器轮询定时任务 zset 的时间间隔，单位：s
		ZRangeGapSeconds: 1,
		// 并发协程数
		WorkersNum: 10000,
	},

	WebServer: &WebServerAppConf{
		Port: 8092,
	},
	Redis: &RedisConfig{
		Network: "tcp",
		// 最大空闲连接数
		MaxIdle: 500,
		// 空闲连接超时时间，单位：s
		IdleTimeoutSeconds: 30,
		// 连接池最大存活的连接数
		MaxActive: 5000,
		// 当连接数达到上限时，新的请求是等待还是立即报错
		Wait: true,
	},
}

type GloablConf struct {
	Migrator  *MigratorAppConf  `yaml:"migrator"`
	Mysql     *MySQLConfig      `yaml:"mysql"`
	Redis     *RedisConfig      `yaml:"redis"`
	Trigger   *TriggerAppConf   `yaml:"trigger"`
	Scheduler *SchedulerAppConf `yaml:"scheduler"`
	WebServer *WebServerAppConf `yaml:"webServer"`
}

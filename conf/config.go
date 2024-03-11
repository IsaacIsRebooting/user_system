package conf

import (
	rlog "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

var (
	config GlobalConfig // 全局业务配置文件
	once   sync.Once
)

type LogConf struct {
	LogPattern string `yaml:"log_pattern" mapstructure:"log_pattern"` // 日志输出标准， 终端输出/文件输出
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`       // 日志路径
	SaveDays   uint   `yaml:"save_days" mapstructure:"save_days"`     // 日志保存天数
	Level      string `yaml:"level" mapstructure:"level"`             // 日志级别
}

type DbConf struct {
	Host        string `yaml:"host" mapstructure:"host"`                   // db主机地址
	Port        string `yaml:"port" mapstructure:"port"`                   // db端口
	User        string `yaml:"user" mapstructure:"user"`                   // 用户名
	Password    string `yaml:"password" mapstructure:"password"`           // 密码
	Dbname      string `yaml:"dbname" mapstructure:"dbname"`               // db名
	MaxIdleConn int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"` // 最大空闲连接数
	MaxOpenConn int    `yaml:"max_open_conn" mapstructure:"max_open_conn"` // 最大打开的连接数
	MaxIdleTime int64  `yaml:"max_idle_time" mapstructure:"max_idle_time"` // 连接最大空闲时间
}
type RedisConf struct {
	Host     string `yaml:"rhost" mapstructure:"rhost"` // db主机地址
	Port     int    `yaml:"rport" mapstructure:"rport"` // db端口
	DB       int    `yaml:"rdb" mapstructure:"rdb"`
	PassWord string `yaml:"passwd" mapstructure:"passwd"`
	PoolSize int    `yaml:"poolsize" mapstructure:"poolsize"`
}

type Cache struct {
	SessionExpired int `yaml:"session_expired" mapstructure:"session_expired"`
	UserExpired    int `yaml:"user_expired" mapstructure:"user_expired"`
}

type Appconf struct {
	AppName string `yaml:"app_name" mapstructure:"app_name"` // 业务名
	Version string `yaml:"version" mapstructure:"version"`   // 版本
	Port    int    `yaml:"port" mapstructure:"port"`         // 端口
	RunMode string `yaml:"run_mode" mapstructure:"run_mode"` // 运行模式
}

type GlobalConfig struct {
	CorsOrigin  []string  `yaml:"cors_origin" mapstructure:"cors_origin"`
	LogConfig   LogConf   `yaml:"log" mapstructure:"log"`
	AppConfig   Appconf   `yaml:"app" mapstructure:"app"`
	DbConfig    DbConf    `yaml:"db" mapstructure:"db"`
	RedisConfig RedisConf `yaml:"redis" mapstructure:"redis"`
	Cache       Cache     `yaml:"cache" mapstructure:"cache"`
}

func GetGlobalConfig() *GlobalConfig {
	once.Do(readConf)
	return &config
}

// readConf 函数用于读取配置文件并加载配置信息
func readConf() {
	// 设置配置文件名和类型
	viper.SetConfigName("app.yml")
	viper.SetConfigType("yml")

	// 添加配置文件搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		// 读取配置文件失败，程序终止并打印错误信息
		panic("read config file err：" + err.Error())
	}
	err = viper.Unmarshal(&config) //将 Viper 实例中的配置数据解析为指定的结构体 config,别忘写了
	// 解析配置文件内容到全局变量 config 中
	if err != nil {
		// 配置文件解析失败，程序终止并打印错误信息
		panic("config file unmarshal err:" + err.Error())
	}

	// 打印配置信息
	log.Infof("config=== %+v1", config)
}

func InitConfig() {
	globalConf := GetGlobalConfig()
	level, err := log.ParseLevel(globalConf.LogConfig.Level)
	if err != nil {
		panic("log level parse err:" + err.Error())
	}
	log.SetFormatter(&logFormatter{
		log.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}})
	log.SetReportCaller(true)
	log.SetLevel(level)
	switch globalConf.LogConfig.LogPattern {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "file":
		logger, err := rlog.New(
			globalConf.LogConfig.LogPath+".%Y%m%d",
			rlog.WithRotationCount(globalConf.LogConfig.SaveDays),
			rlog.WithRotationTime(time.Hour*24),
		)
		if err != nil {
			panic("log conf err: " + err.Error())
		}
		log.SetOutput(logger)
	default:
		panic("log conf err, check log_pattern is logsvr.yaml")
	}

}

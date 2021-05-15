package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//使用结构体变量保存结构体信息
// Conf ... 全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

// AppConfig ...
type AppConfig struct {
	Name    string `mapstructure:"app_name"`
	Mode    string `mapstructure:"mod"`
	Version string `mapstructure:"version"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

// LogConfig ...
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// MySQLConfig ...
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBname       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig ...
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Init(filepath string) (err error) {
	// viper.SetConfigName("config")
	// viper.SetConfigFile("./config.yaml") //配置文件的位置 相对于可执行文件的位置(即相当于main.go的位置)
	viper.SetConfigFile(filepath)
	err = viper.ReadInConfig() // 读取配置文件信息
	if err != nil {            // 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() error:%v", err)
		return
	}
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
	}
	//把读取到的配置信息读取到Conf变量中
	// 监控配置文件变化
	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件已被修改...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		}
	})
	return
}

package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigFile("./config.yaml")
	err = viper.ReadInConfig() // 读取配置文件信息
	if err != nil {            // 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() error:%v", err)
		return
	}

	// 监控配置文件变化
	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件已被修改...")
	})
	return
}

package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.paddword"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
	})

	_, err = rdb.Ping().Result()
	return
}

func Close() {
	_ = rdb.Close()
}

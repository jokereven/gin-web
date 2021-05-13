package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdbpipeline *redis.Client

// 初始化连接
func initClientpipeline() (err error) {
	rdbpipeline = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	_, err = rdbpipeline.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

/*
	Pipeline
Pipeline 主要是一种网络优化。它本质上意味着客户端缓冲一堆命令并一次性将它们发送到服务器。这些命令不能保证在事务中执行。这样做的好处是节省了每个命令的网络往返时间（RTT）。
*/

func pipelineDome() {
	pipe := rdbpipeline.Pipeline()

	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)

	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}

/*
	事务
Redis是单线程的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行，例如在它们之间交替执行。但是，Multi/exec能够确保在multi/exec两个语句之间的命令之间没有其他客户端正在执行命令。

在这种场景我们需要使用TxPipeline。TxPipeline总体上类似于上面的Pipeline，但是它内部会使用MULTI/EXEC包裹排队的命令。
*/

func execpipeline() {
	pipe := rdbpipeline.TxPipeline()

	incr := pipe.Incr("tx_pipeline_counter")
	pipe.Expire("tx_pipeline_counter", time.Hour)

	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}

/*
	Watch
在某些场景下，我们除了要使用MULTI/EXEC命令外，还需要配合使用WATCH命令。在用户使用WATCH命令监视某个键之后，直到该用户执行EXEC命令的这段时间里，如果有其他用户抢先对被监视的键进行了替换、更新、删除等操作，那么当用户尝试执行EXEC的时候，事务将失败并返回一个错误，用户可以根据这个错误选择重试事务或者放弃事务。

Watch(fn func(*Tx) error, keys ...string) error
*/

func watchpipe() {
	// 监视watch_count的值，并在值不变的前提下将其值+1
	key := "watch_count"
	err := rdbpipeline.Watch(func(tx *redis.Tx) error {
		n, err := tx.Get(key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Set(key, n+1, 0)
			return nil
		})
		return err
	}, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("sucess...")
}
func main_pipwline_watch() {
	err := initClientpipeline() //连接redis数据库
	if err != nil {
		fmt.Printf("go connect redis err:%v\n", err)
		return
	}
	fmt.Println("go connect redis sucess...")

	pipelineDome() //网络优化

	execpipeline() //事务

	watchpipe()         //watchpipe
	rdbpipeline.Close() //程序退出关闭服务
}

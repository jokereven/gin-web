package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitee.com/jokereven/bluebell/src/dao/mysql"
	"gitee.com/jokereven/bluebell/src/dao/redis"
	"gitee.com/jokereven/bluebell/src/logger"
	"gitee.com/jokereven/bluebell/src/routes"
	"gitee.com/jokereven/bluebell/src/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//Go web 比较通用的模板脚手架

func main() {
	if len(os.Args) < 2 {
		fmt.Println("输入无效...请输入 可执行文件+filepath")
		return
	}

	var path string
	flag.StringVar(&path, "path", "config.yaml", "the_path")

	//1. 加载配置文件
	if err := settings.Init(path); err != nil {
		fmt.Printf("init settings file ,err:%v", err)
		return
	}
	//2. 初始化配置文件
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger file ,err:%v", err)
		return
	}
	defer zap.L().Sync() //延迟注册
	zap.L().Debug("logger init sucess...")
	//3. 初始化MySQL连接
	if err := mysql.Init(); err != nil { //通过全局变量获取
		fmt.Printf("init mysql file ,err:%v", err)
		return
	}
	defer mysql.Close()
	//4. 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis file ,err:%v", err)
		return
	}
	defer redis.Close()
	//5. 注册路由
	r := routes.Setup()
	//6. 启动服务(幽雅关机)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}

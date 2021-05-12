package main

import (
	"net/http"

	"go.uber.org/zap"
)

/*
	在Go语言项目中使用Zap日志库
介绍
在许多Go语言项目中，我们需要一个好的日志记录器能够提供下面这些功能：

能够将事件记录到文件中，而不是应用程序控制台。
日志切割-能够根据文件大小、时间或间隔等来切割日志文件。
支持不同的日志级别。例如INFO，DEBUG，ERROR等。
能够打印基本信息，如调用文件/函数名和行号，日志时间等。

配置Zap Logger
Zap提供了两种类型的日志记录器—Sugared Logger和Logger。

在性能很好但不是很关键的上下文中，使用SugaredLogger。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。

在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。

Logger
通过调用zap.NewProduction()/zap.NewDevelopment()或者zap.Example()创建一个Logger。
上面的每一个函数都将创建一个logger。唯一的区别在于它将记录的信息不同。例如production logger默认记录调用函数信息、日期和时间等。
通过Logger调用Info/Error等。
默认情况下日志都会打印到应用程序的console界面。
*/

var logger *zap.Logger

var sugarLogger *zap.SugaredLogger

func main_zap() {
	InitLogger_zap()
	defer logger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}

func InitLogger_zap() {
	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

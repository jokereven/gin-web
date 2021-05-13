package main

import (
	"net/http"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger_custom *zap.Logger

var sugarLogger_custom *zap.SugaredLogger

func InitLogger_custom() {
	writeSyncer := getLogWriter_splice()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger_custom = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) //json格式日志
	// return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// func getLogWriter() zapcore.WriteSyncer {
// 	// file, _ := os.Create("./test.log") //覆盖
// 	file, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744) //追加
// 	return zapcore.AddSync(file)
// }

func simpleHttpGet_custom(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger_custom.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger_custom.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

/*
	更改时间编码并添加调用者详细信息
鉴于我们对配置所做的更改，有下面两个问题：

时间是以非人类可读的方式展示，例如1.572161051846623e+09
调用方函数的详细信息没有显示在日志中
我们要做的第一件事是覆盖默认的ProductionConfig()，并进行以下更改:

修改时间编码器
在日志文件中使用大写字母记录日志级别

接下来，我们将修改zap logger代码，添加将调用函数信息记录到日志中的功能。为此，我们将在zap.New(..)函数中添加一个Option。

logger := zap.New(core, zap.AddCaller())
*/

func getLogWriter_splice() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     //大小单位是m
		MaxBackups: 5,     //最大备份数
		MaxAge:     30,    //最大备天数
		Compress:   false, //是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func main_zap_logger() {
	InitLogger_custom()
	defer logger_custom.Sync()

	for i := 0; i < 100000; i++ {
		sugarLogger_custom.Info("test logger ...")
	}
	simpleHttpGet_custom("www.baidu.com")
	simpleHttpGet_custom("http://www.baidu.com")
}

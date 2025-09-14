package middleware

import (
	"io"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true, // 包含完整的时间戳
	}
	logger.SetReportCaller(true)

	// 配置lumberjack.Logger
	logFile := &lumberjack.Logger{
		Filename:   "../../../.logs/app.log", // 日志文件的位置和名称
		MaxSize:    100,                      // 每个日志文件保存的最大尺寸 单位：MB
		MaxBackups: 7,                        // 日志文件最多保存多少个备份
		MaxAge:     30,                       // 文件最多保存多少天
		Compress:   true,                     // 是否压缩/归档过期的日志文件
	}

	writer := logFile
	writerCloser := io.WriteCloser(writer)
	logger.SetOutput(writerCloser) // 设置输出到lumberjack.Logger实例

	return logger
}

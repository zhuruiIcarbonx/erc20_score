package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {

	Log = logrus.New()

	Log.Level = logrus.DebugLevel

	// 添加默认格式化器
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置标准输出
	Log.SetOutput(os.Stdout)

	Log.Info("日志系统初始化完成")

	// Log.Formatter = &logrus.TextFormatter{
	// 	FullTimestamp: true, // 包含完整的时间戳
	// }
	// Log.SetReportCaller(true)
	// Log.SetLevel(logrus.DebugLevel)

	// // 配置lumberjack.Logger
	// logFile := &lumberjack.Logger{
	// 	Filename:   "../../.logs/app.log", // 日志文件的位置和名称
	// 	MaxSize:    100,                   // 每个日志文件保存的最大尺寸 单位：MB
	// 	MaxBackups: 7,                     // 日志文件最多保存多少个备份
	// 	MaxAge:     30,                    // 文件最多保存多少天
	// 	Compress:   true,                  // 是否压缩/归档过期的日志文件
	// }

	// writer := logFile
	// writerCloser := io.WriteCloser(writer)
	// Log.SetOutput(writerCloser) // 设置输出到lumberjack.Logger实例

}

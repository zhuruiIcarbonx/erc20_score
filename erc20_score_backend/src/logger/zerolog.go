package logger

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func InitLog() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	/***
	    panic (zerolog.PanicLevel, 5)
		fatal (zerolog.FatalLevel, 4)
		error (zerolog.ErrorLevel, 3)
		warn (zerolog.WarnLevel, 2)
		info (zerolog.InfoLevel, 1)
		debug (zerolog.DebugLevel, 0)
		trace (zerolog.TraceLevel, -1)

	***/
	zerolog.SetGlobalLevel(zerolog.DebugLevel) //0

	proj, _ := filepath.Abs(filepath.Dir("."))
	log.Info().Str("项目路径", proj).Msg("")
	log.Info().Str("app.log路径", proj+"/.logs/app.log").Msg("")
	file, _ := os.OpenFile(proj+"/.logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	Log = zerolog.New(multi).With().Timestamp().Logger()

	//要记录人性化的彩色输出，请使用zerolog.ConsoleWriter
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal().Msgf("无法打开日志文件: %v", err)
	// }
	// defer file.Close()
	// // 使用文件作为输出创建logger
	// logger := log.Output(file).With().Timestamp().Logger()

	//demo
	// Log.Print("hello world")
	// Log.Info().Msg("hello world")
	// Log.Debug().Msg("This message appears only when log level set to Debug")

	// err := errors.New("seems we have an error here")
	// Log.Error().Err(err).Msg("")

	// err1 := errors.New("A repo man spends his life getting into tense situations")

	// Log.Log().
	// 	Err(err1).
	// 	Str("foo", "bar").
	// 	Msg("")
}

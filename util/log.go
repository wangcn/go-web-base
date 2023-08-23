package util

import (
	"context"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

// Log 获取全局日志
func Log() *zap.Logger {
	return log
}

// SetLog 初始化时设置全局日志
func SetLog(logger *zap.Logger) {
	log = logger
}

func LogWithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return log
	}
	return log.With(zap.String("request_id", ctx.Value("request_id").(string)))
}

func InitLog() {
	var logCoreArr []zapcore.Core
	zcfg := zap.NewProductionEncoderConfig()
	zcfg.EncodeTime = TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(zcfg)
	logLevel := zap.NewAtomicLevelAt(zapcore.Level(viper.GetInt("log.level")))
	if viper.GetBool("log.stdout") {
		stdoutW := zapcore.Lock(os.Stdout)
		logCoreArr = append(logCoreArr, zapcore.NewCore(jsonEncoder, stdoutW, logLevel))
	}
	outputFile := viper.GetString("log.file")
	if outputFile != "" {
		fileW := zapcore.AddSync(&lumberjack.Logger{
			Filename:   outputFile,
			MaxSize:    viper.GetInt("log.file_size"), // megabytes
			MaxBackups: viper.GetInt("log.max_backup"),
			MaxAge:     viper.GetInt("log.max_age"), // days
		})
		logCoreArr = append(logCoreArr, zapcore.NewCore(jsonEncoder, fileW, logLevel))
	}
	if len(logCoreArr) == 0 {
		panic("must set one log path or use stdout")
	}
	core := zapcore.NewTee(logCoreArr...)
	log := zap.New(core)
	SetLog(log)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}

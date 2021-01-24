package log

import (
	"os"

	"github.com/tinysrc/z9go/pkg/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func toZapLevel(str string) (level zapcore.Level) {
	switch str {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}
	return
}

func newLogger() *zap.Logger {
	var cores []zapcore.Core
	// json encoder
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)
	// log to stdout
	level := toZapLevel(conf.Global.GetString("log.level"))
	if conf.Global.GetBool("log.console") {
		stdoutWriter := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(jsonEncoder, stdoutWriter, level))
	}
	// log to file
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.Global.GetString("log.filename"),
		MaxSize:    conf.Global.GetInt("log.maxSize"),
		MaxAge:     conf.Global.GetInt("log.maxAge"),
		MaxBackups: conf.Global.GetInt("log.maxBackups"),
		Compress:   conf.Global.GetBool("log.compress"),
	})
	cores = append(cores, zapcore.NewCore(jsonEncoder, fileWriter, level))
	// merge
	merge := zapcore.NewTee(cores...)
	return zap.New(merge)
}

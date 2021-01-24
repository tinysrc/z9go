package log

import (
	"github.com/tinysrc/z9go/pkg/conf"
	"go.uber.org/zap"
)

var logger *zap.Logger

func initConfig() {
	conf.Global.SetDefault("log.console", true)
	conf.Global.SetDefault("log.level", "debug")
	conf.Global.SetDefault("log.filename", "./service.log")
	conf.Global.SetDefault("log.maxSize", 100)
	conf.Global.SetDefault("log.maxAge", 0)
	conf.Global.SetDefault("log.maxBackups", 0)
	conf.Global.SetDefault("log.compress", false)
}

func init() {
	initConfig()
	logger = newLogger()
	if logger == nil {
		panic("new logger failed")
	}
}

// Close make sure the log is finished
func Close() {
	logger.Sync()
}

// Debug wrap
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info wrap
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn wrap
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error wrap
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// DPanic wrap
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

// Panic wrap
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal wrap
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

package log

import (
	"github.com/tinysrc/z9go/pkg/conf"
	"go.uber.org/zap"
)

// Logger instance
var Logger *zap.Logger

func init() {
	conf.Global.SetDefault("log.console", true)
	conf.Global.SetDefault("log.level", "debug")
	conf.Global.SetDefault("log.filename", "./service.log")
	conf.Global.SetDefault("log.maxSize", 100)
	conf.Global.SetDefault("log.maxAge", 0)
	conf.Global.SetDefault("log.maxBackups", 0)
	conf.Global.SetDefault("log.compress", false)
	Logger = newLogger()
	if Logger == nil {
		panic("new logger failed")
	}
}

// Close make sure the log is finished
func Close() {
	Logger.Sync()
}

// Debug wrap
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info wrap
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn wrap
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error wrap
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// DPanic wrap
func DPanic(msg string, fields ...zap.Field) {
	Logger.DPanic(msg, fields...)
}

// Panic wrap
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

// Fatal wrap
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

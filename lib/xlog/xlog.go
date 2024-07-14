package xlog

import (
	"go.uber.org/zap"
	"sync"
)

var (
	globalLogger     *zap.Logger
	globalLoggerLock sync.RWMutex
)

func SetGlobalLogger(logger *zap.Logger) {
	globalLoggerLock.Lock()
	defer globalLoggerLock.Unlock()

	globalLogger = logger

}

func GetGlobalLogger() *zap.Logger {
	globalLoggerLock.RLock()
	defer globalLoggerLock.RUnlock()

	return globalLogger
}

func Debug(msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	log.Warn(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	log.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	log.Error(msg, fields...)
}

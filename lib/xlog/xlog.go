package xlog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

const (
	DEBUG  = "DEBUG"
	INFO   = "INFO"
	WARN   = "WARN"
	ERROR  = "ERROR"
	DPANIC = "DPANIC"
	PANIC  = "PANIC"
	FATAL  = "FATAL"
)

var (
	logger      *zap.Logger
	loggerMutex sync.RWMutex
)

func SetDefaultLogger(logLevel string, development bool) error {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(getZapLogLevel(logLevel)),
		Development:       development,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	logger = log
	return nil
}

func GetGlobalLogger() *zap.Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return logger
}

func getZapLogLevel(logLevel string) zapcore.Level {
	switch strings.ToUpper(logLevel) {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	case DPANIC:
		return zapcore.DPanicLevel
	case PANIC:
		return zapcore.PanicLevel
	case FATAL:
		return zapcore.FatalLevel
	default:
		panic(fmt.Sprintf("invalid log level: %s", logLevel))
	}
}

// TODO: context fields

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.Error(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.Warn(msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.Debug(msg, fields...)
}

func DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.DPanic(msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.Panic(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	log := GetGlobalLogger()
	ctxFields := GetContextFields(ctx)
	fields = mergeFields(ctxFields, fields)
	log.Fatal(msg, fields...)
}

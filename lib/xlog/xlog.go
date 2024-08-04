package xlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
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

	logger, err := config.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
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

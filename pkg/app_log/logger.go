package app_log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(level string, format string) error {
	config := zap.NewProductionConfig()
	
	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Set encoding format
	if format == "console" {
		config.Encoding = "console"
	} else {
		config.Encoding = "json"
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Debug(msg string, fields ...zapcore.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	Logger.Fatal(msg, fields...)
}

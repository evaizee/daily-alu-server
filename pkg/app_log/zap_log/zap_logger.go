package zap_log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(level, format, logPath string) error {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if format == "console" {
		encoder = zapcore.NewConsoleEncoder(config)
	} else {
		encoder = zapcore.NewJSONEncoder(config)
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	writeSyncer := zapcore.AddSync(logFile)

	var levelEnabler zapcore.Level
	switch level {
	case "debug":
		levelEnabler = zap.DebugLevel
	case "info":
		levelEnabler = zap.InfoLevel
	case "warn":
		levelEnabler = zap.WarnLevel
	case "error":
		levelEnabler = zap.ErrorLevel
	default:
		levelEnabler = zap.InfoLevel
	}

	core := zapcore.NewCore(encoder, writeSyncer, levelEnabler)
	Logger = zap.New(core)

	return nil
}

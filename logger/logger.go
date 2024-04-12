package logger

import (
	"os"

	"go.uber.org/zap"
)

var logger *zap.Logger

func GetLogger() *zap.Logger {
	if logger == nil {
		lvl, ok := os.LookupEnv("LOG_LEVEL")
		// LOG_LEVEL not set, let's default to debug
		if !ok {
			lvl = "error"
		}
		ll, err := zap.ParseAtomicLevel(lvl)
		if err != nil {
			ll = zap.NewAtomicLevelAt(zap.ErrorLevel)
		}
		cfg := zap.NewProductionConfig()
		cfg.Level = ll
		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = logger.Sync()
		}()
	}
	return logger
}

func Error(msg string, fields ...zap.Field) {
	logger = GetLogger()
	logger.Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger = GetLogger()
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger = GetLogger()
	logger.Debug(msg, fields...)
}

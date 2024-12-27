package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger")
	}
	defer Log.Sync()
}

func Info(message string, fields ...zap.Field) {
	Log.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	Log.Error(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	Log.Debug(message, fields...)
}

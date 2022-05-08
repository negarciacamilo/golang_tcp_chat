package logger

import "go.uber.org/zap"

var log *zap.Logger

func New() {
	// I don't actually have any scope, so I will use the production one
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	log = l
}

func Panic(msg string, fields ...zap.Field) {
	log.Panic(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

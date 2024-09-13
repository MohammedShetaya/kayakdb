package log

import (
	"go.uber.org/zap"
	"log"
)

var (
	Logger *zap.Logger
)

func InitLogger() *zap.Logger {
	Logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
	return Logger
}

package utils

import (
	"go.uber.org/zap"
	"log"
)

var (
	Log *zap.Logger
)

func InitLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
}

func Sync() {
	if err := Log.Sync(); err != nil {
		log.Fatalf("failed to flush logs: %v", err)
	}
}

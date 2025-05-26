package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Logger *zap.Logger
)

func InitLogger(level string) *zap.Logger {

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		// If the level is invalid, default to InfoLevel
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel) // Set log level to Debug

	Logger, err := config.Build()
	if err != nil {
		_ = fmt.Errorf("cannot initialize logger: %w", err)
		os.Exit(1)
	}
	return Logger
}

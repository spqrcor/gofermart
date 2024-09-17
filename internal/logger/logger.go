package logger

import (
	"github.com/spqrcor/gofermart/internal/config"
	"go.uber.org/zap"
	"log"
)

var Log *zap.Logger = zap.NewNop()

func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(config.Cfg.LogLevel)
	zl, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return zl
}

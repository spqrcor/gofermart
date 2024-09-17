package config

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type Config struct {
	RunAddr              string        `env:"RUN_ADDRESS"`
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             zapcore.Level `env:"LOG_LEVEL"`
	DatabaseURI          string        `env:"DATABASE_URI"`
	QueryTimeOut         time.Duration `env:"QUERY_TIME_OUT"`
	WorkerCount          int           `env:"WORKER_COUNT"`
}

var Cfg = Config{
	RunAddr:              "localhost:8080",
	LogLevel:             zap.InfoLevel,
	AccrualSystemAddress: "",
	DatabaseURI:          "postgres://postgres:Sp123456@localhost:5432/gofermart?sslmode=disable",
	QueryTimeOut:         3,
	WorkerCount:          3,
}

func Init() {
	flag.StringVar(&Cfg.RunAddr, "a", Cfg.RunAddr, "address and port to run server")
	flag.StringVar(&Cfg.AccrualSystemAddress, "r", Cfg.AccrualSystemAddress, "accrual system address")
	flag.StringVar(&Cfg.DatabaseURI, "d", Cfg.DatabaseURI, "database uri")
	flag.Parse()

	serverAddressEnv, findAddress := os.LookupEnv("RUN_ADDRESS")
	serverDatabaseURI, findDatabaseURI := os.LookupEnv("DATABASE_URI")
	serverAccrualSystemAddress, findAccrualSystemAddress := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")

	if findAddress {
		Cfg.RunAddr = serverAddressEnv
	}
	if findDatabaseURI {
		Cfg.DatabaseURI = serverDatabaseURI
	}
	if findAccrualSystemAddress {
		Cfg.AccrualSystemAddress = serverAccrualSystemAddress
	}
}

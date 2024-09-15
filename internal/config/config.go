package config

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	RunAddr              string        `env:"RUN_ADDRESS"`
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             zapcore.Level `env:"LOG_LEVEL"`
	DatabaseURI          string        `env:"DATABASE_URI"`
}

var Cfg = Config{
	RunAddr:              "localhost:8080",
	LogLevel:             zap.InfoLevel,
	AccrualSystemAddress: "",
	DatabaseURI:          "",
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

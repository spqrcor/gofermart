package config

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

type Config struct {
	RunAddr               string        `env:"RUN_ADDRESS"`
	AccrualSystemAddress  string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel              zapcore.Level `env:"LOG_LEVEL"`
	DatabaseURI           string        `env:"DATABASE_URI"`
	QueryTimeOut          time.Duration `env:"QUERY_TIME_OUT"`
	WorkerCount           int           `env:"WORKER_COUNT"`
	RetryStartWorkerCount int           `env:"RETRY_START_WORKER_COUNT"`
	SecretKey             string        `env:"SECRET_KEY"`
	TokenExp              time.Duration `env:"TOKEN_EXP"`
}

var cfg = Config{
	RunAddr:               "localhost:8080",
	LogLevel:              zap.InfoLevel,
	AccrualSystemAddress:  "",
	DatabaseURI:           "",
	QueryTimeOut:          3,
	WorkerCount:           3,
	RetryStartWorkerCount: 10,
	SecretKey:             "KLJ-fo3Fksd3fl!=",
	TokenExp:              time.Hour * 3,
}
var once sync.Once

func NewConfig() Config {
	once.Do(func() {
		flag.StringVar(&cfg.RunAddr, "a", cfg.RunAddr, "address and port to run server")
		flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "accrual system address")
		flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database uri")
		flag.Parse()

		serverAddressEnv, findAddress := os.LookupEnv("RUN_ADDRESS")
		serverDatabaseURI, findDatabaseURI := os.LookupEnv("DATABASE_URI")
		serverAccrualSystemAddress, findAccrualSystemAddress := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")

		if findAddress {
			cfg.RunAddr = serverAddressEnv
		}
		if findDatabaseURI {
			cfg.DatabaseURI = serverDatabaseURI
		}
		if findAccrualSystemAddress {
			cfg.AccrualSystemAddress = serverAccrualSystemAddress
		}
	})
	return cfg
}

package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spqrcor/gofermart/internal/config"
	"go.uber.org/zap"
)

func connect(logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.Cfg.DatabaseURI)
	if err != nil {
		logger.Fatal(err.Error())
		return nil, err
	}
	if err := db.Ping(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return db, nil
}

func NewDB(logger *zap.Logger) *sql.DB {
	res, err := connect(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal(err.Error())
	}
	if err := goose.Up(res, "internal/migrations"); err != nil {
		logger.Fatal(err.Error())
	}
	return res
}

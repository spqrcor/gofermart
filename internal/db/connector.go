package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"time"
)

func connect(logger *zap.Logger, databaseURI string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		logger.Fatal(err.Error())
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return db, nil
}

func NewDB(logger *zap.Logger, databaseURI string) (*sql.DB, error) {
	res, err := connect(logger, databaseURI)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Migrate(res *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(res, "internal/migrations"); err != nil {
		return err
	}
	return nil
}

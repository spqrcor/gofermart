package db

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"github.com/spqrcor/gofermart/internal/logger"
)

func Migrate(db *sql.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Log.Fatal(err.Error())
	}

	if err := goose.Up(db, "internal/migrations"); err != nil {
		logger.Log.Fatal(err.Error())
	}
}

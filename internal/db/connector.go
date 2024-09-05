package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/logger"
)

var Source *sql.DB

func connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", config.Cfg.DatabaseURI)
	if err != nil {
		logger.Log.Fatal(err.Error())
		return nil, err
	}
	if err := db.Ping(); err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return db, nil
}

func Init() {
	res, err := connect()
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	Migrate(res)
	Source = res
}

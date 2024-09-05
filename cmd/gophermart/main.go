package main

import (
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/server"
)

func main() {
	config.Init()
	logger.Init()
	db.Init()

	server.Start()
}

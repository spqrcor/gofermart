package main

import (
	"context"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/server"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/workers"
)

func main() {
	mainCtx := context.Background()
	config.Init()
	loggerRes := logger.NewLogger()
	dbRes := db.NewDB(loggerRes)

	orderQueue := make(chan string)
	defer close(orderQueue)

	userService := services.NewUserService(dbRes)
	orderService := services.NewOrderService(orderQueue, dbRes, loggerRes)
	withdrawalService := services.NewWithdrawalService(dbRes, loggerRes)

	orderWorker := workers.NewOrderWorker(mainCtx, orderService, orderQueue, loggerRes)
	orderWorker.Run()

	appServer := server.NewServer(userService, orderService, withdrawalService, loggerRes)
	if err := appServer.Start(); err != nil {
		loggerRes.Fatal(err.Error())
	}
}

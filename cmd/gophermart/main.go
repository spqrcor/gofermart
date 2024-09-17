package main

import (
	"context"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/server"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/workers"
	"log"
)

func main() {
	mainCtx := context.Background()
	conf := config.NewConfig()

	loggerRes, err := logger.NewLogger(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	dbRes, err := db.NewDB(loggerRes, conf.DatabaseURI)
	if err != nil {
		loggerRes.Fatal(err.Error())
	}

	orderQueue := make(chan string)
	defer close(orderQueue)

	userService := services.NewUserService(dbRes, conf.QueryTimeOut)
	orderService := services.NewOrderService(orderQueue, dbRes, loggerRes, conf.QueryTimeOut)
	withdrawalService := services.NewWithdrawalService(dbRes, loggerRes, conf.QueryTimeOut)
	orderClient := client.NewOrderClient(loggerRes, conf.AccrualSystemAddress)

	orderWorker := workers.NewOrderWorker(mainCtx, orderService, orderQueue, loggerRes, conf.WorkerCount, orderClient)
	orderWorker.Run()

	appServer := server.NewServer(userService, orderService, withdrawalService, loggerRes, conf.RunAddr)
	if err := appServer.Start(); err != nil {
		loggerRes.Error(err.Error())
	}
}

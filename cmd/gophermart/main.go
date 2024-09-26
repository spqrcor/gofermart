package main

import (
	"context"
	"errors"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/client"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/server"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/workers"
	"log"
	"net/http"
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
	if err := db.Migrate(dbRes); err != nil {
		loggerRes.Fatal(err.Error())
	}

	orderQueue := make(chan string)
	defer close(orderQueue)

	authService := authenticate.NewAuthenticateService(
		authenticate.WithLogger(loggerRes),
		authenticate.WithSecretKey(conf.SecretKey),
		authenticate.WithTokenExp(conf.TokenExp),
	)
	userService := services.NewUserService(dbRes, conf.QueryTimeOut)
	orderService := services.NewOrderService(orderQueue, dbRes, loggerRes, conf.QueryTimeOut)
	withdrawalService := services.NewWithdrawalService(dbRes, loggerRes, conf.QueryTimeOut)
	orderClient := client.NewOrderClient(
		client.WithLogger(loggerRes),
		client.WithAccrualSystemAddress(conf.AccrualSystemAddress),
	)

	orderWorker := workers.NewOrderWorker(
		workers.WithCtx(mainCtx),
		workers.WithOrderService(orderService),
		workers.WithOrderQueue(orderQueue),
		workers.WithLogger(loggerRes),
		workers.WithConfig(conf),
		workers.WithOrderClient(orderClient),
	)
	orderWorker.Run()

	appServer := server.NewServer(
		server.WithUserService(userService),
		server.WithOrderService(orderService),
		server.WithWithdrawalService(withdrawalService),
		server.WithLogger(loggerRes),
		server.WithRunAddress(conf.RunAddr),
		server.WithAuthService(authService),
	)
	err = appServer.Start()
	if errors.Is(err, http.ErrServerClosed) {
		loggerRes.Error("Server stop")
	}
	if err != nil {
		loggerRes.Error(err.Error())
	}
}

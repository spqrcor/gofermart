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
	logger.Init()
	db.Init()

	userService := services.NewUserService()
	orderService := services.NewOrderService()
	withdrawalService := services.NewWithdrawalService()

	orderWorker := workers.NewOrderWorker(mainCtx, orderService)
	orderWorker.Run()

	server.Start(userService, orderService, withdrawalService)
}

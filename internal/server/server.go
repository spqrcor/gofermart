package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/handlers"
	"github.com/spqrcor/gofermart/internal/services"
	"log"
	"net/http"
)

var publicRoutes = []string{"/api/user/register", "/api/user/login"}

func Start() {
	r := chi.NewRouter()
	r.Use(authenticateMiddleware)
	r.Use(loggerMiddleware)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(getBodyMiddleware)

	userService := services.NewUserService()
	orderService := services.NewOrderService()
	withdrawalService := services.NewWithdrawalService()

	r.Post("/api/user/register", handlers.RegisterHandler(userService))
	r.Post("/api/user/login", handlers.LoginHandler(userService))
	r.Post("/api/user/orders", handlers.AddOrdersHandler(orderService))
	r.Get("/api/user/orders", handlers.GetOrdersHandler(orderService))
	r.Get("/api/user/balance", handlers.GetBalanceHandler(withdrawalService))
	r.Post("/api/user/balance/withdraw", handlers.AddWithdrawalHandler(withdrawalService))
	r.Get("/api/user/withdrawals", handlers.GetWithdrawalsHandler(withdrawalService))

	r.HandleFunc(`/*`, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})

	log.Fatal(http.ListenAndServe(config.Cfg.RunAddr, r))
}

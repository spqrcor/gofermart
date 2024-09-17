package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/handlers"
	"github.com/spqrcor/gofermart/internal/services"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	userService       services.UserRepository
	orderService      services.OrderRepository
	withdrawalService services.WithdrawalRepository
	logger            *zap.Logger
}

func NewServer(userService services.UserRepository, orderService services.OrderRepository, withdrawalService services.WithdrawalRepository, logger *zap.Logger) *HTTPServer {
	return &HTTPServer{
		userService: userService, orderService: orderService, withdrawalService: withdrawalService, logger: logger,
	}
}

func (s *HTTPServer) Start() error {
	r := chi.NewRouter()
	r.Use(loggerMiddleware(s.logger))
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(getBodyMiddleware(s.logger))

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handlers.RegisterHandler(s.userService))
		r.Post("/api/user/login", handlers.LoginHandler(s.userService))
	})

	r.Group(func(r chi.Router) {
		r.Use(authenticateMiddleware(s.logger))
		r.Post("/api/user/orders", handlers.AddOrdersHandler(s.orderService))
		r.Get("/api/user/orders", handlers.GetOrdersHandler(s.orderService))
		r.Get("/api/user/balance", handlers.GetBalanceHandler(s.withdrawalService))
		r.Post("/api/user/balance/withdraw", handlers.AddWithdrawalHandler(s.withdrawalService))
		r.Get("/api/user/withdrawals", handlers.GetWithdrawalsHandler(s.withdrawalService))
	})

	r.HandleFunc(`/*`, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})
	return http.ListenAndServe(config.Cfg.RunAddr, r)
}

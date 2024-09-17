package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/handlers"
	"github.com/spqrcor/gofermart/internal/services"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	userService       *services.UserService
	orderService      *services.OrderService
	withdrawalService *services.WithdrawalService
	logger            *zap.Logger
	runAddress        string
	authService       *authenticate.Authenticate
}

func NewServer(userService *services.UserService, orderService *services.OrderService, withdrawalService *services.WithdrawalService, logger *zap.Logger, runAddress string, authService *authenticate.Authenticate) *HTTPServer {
	return &HTTPServer{
		userService: userService, orderService: orderService, withdrawalService: withdrawalService, logger: logger, runAddress: runAddress, authService: authService,
	}
}

func (s *HTTPServer) Start() error {
	r := chi.NewRouter()
	r.Use(loggerMiddleware(s.logger))
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(getBodyMiddleware(s.logger))

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handlers.RegisterHandler(s.userService, s.authService))
		r.Post("/api/user/login", handlers.LoginHandler(s.userService, s.authService))
	})

	r.Group(func(r chi.Router) {
		r.Use(authenticateMiddleware(s.logger, s.authService))
		r.Post("/api/user/orders", handlers.AddOrdersHandler(s.orderService))
		r.Get("/api/user/orders", handlers.GetOrdersHandler(s.orderService))
		r.Get("/api/user/balance", handlers.GetBalanceHandler(s.withdrawalService))
		r.Post("/api/user/balance/withdraw", handlers.AddWithdrawalHandler(s.withdrawalService))
		r.Get("/api/user/withdrawals", handlers.GetWithdrawalsHandler(s.withdrawalService))
	})

	r.HandleFunc(`/*`, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	})
	return http.ListenAndServe(s.runAddress, r)
}

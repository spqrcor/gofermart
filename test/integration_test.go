package tests

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/services"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	login := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 10)

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
	orderService := services.NewOrderService(orderQueue, dbRes, loggerRes, conf.QueryTimeOut)

	u := url.URL{
		Scheme: "http",
		Host:   conf.RunAddr,
	}
	e := httpexpect.Default(t, u.String())

	t.Run("register not post", func(t *testing.T) {
		e.GET("/api/user/register").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("register short password", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(services.InputDataUser{
				Login:    login,
				Password: "333",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("register short login", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(services.InputDataUser{
				Login:    "l",
				Password: password,
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("success register", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(services.InputDataUser{
				Login:    login,
				Password: password,
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("user exists", func(t *testing.T) {
		e.POST("/api/user/register").
			WithJSON(services.InputDataUser{
				Login:    login,
				Password: password,
			}).
			Expect().
			Status(http.StatusConflict)
	})

	t.Run("bad login or password", func(t *testing.T) {
		e.POST("/api/user/login").
			WithJSON(services.InputDataUser{
				Login:    login,
				Password: "333",
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("success login", func(t *testing.T) {
		e.POST("/api/user/login").
			WithJSON(services.InputDataUser{
				Login:    login,
				Password: password,
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("get orders no content", func(t *testing.T) {
		e.GET("/api/user/orders").
			Expect().
			Status(http.StatusNoContent)
	})

	t.Run("add order format error", func(t *testing.T) {
		e.POST("/api/user/orders").
			WithText("343").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("add order success", func(t *testing.T) {
		e.POST("/api/user/orders").
			WithText("9399142970086005").
			Expect().
			Status(http.StatusAccepted)
	})

	t.Run("add order exists", func(t *testing.T) {
		e.POST("/api/user/orders").
			WithText("9399142970086005").
			Expect().
			Status(http.StatusOK)
	})

	t.Run("get orders have content", func(t *testing.T) {
		e.GET("/api/user/orders").
			Expect().
			Status(http.StatusOK)
	})

	t.Run("get balance zero", func(t *testing.T) {
		e.GET("/api/user/balance").
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsKey("current").HasValue("current", 0).
			ContainsKey("current").HasValue("current", 0)
	})

	t.Run("get withdrawals no content", func(t *testing.T) {
		e.GET("/api/user/withdrawals").
			Expect().
			Status(http.StatusNoContent)
	})

	t.Run("add withdraw balance error", func(t *testing.T) {
		e.POST("/api/user/balance/withdraw").
			WithJSON(services.InputWithdrawal{
				OrderNum: "9399142970086005",
				Sum:      100,
			}).
			Expect().
			Status(http.StatusPaymentRequired)
	})

	t.Run("add withdraw order format error", func(t *testing.T) {
		e.POST("/api/user/balance/withdraw").
			WithJSON(services.InputWithdrawal{
				OrderNum: "333",
				Sum:      100,
			}).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	_ = orderService.ChangeStatus(context.Background(), services.OrderFromAccrual{Status: "PROCESSED", Order: "9399142970086005", Accrual: 500})

	t.Run("add withdraw success", func(t *testing.T) {
		e.POST("/api/user/balance/withdraw").
			WithJSON(services.InputWithdrawal{
				OrderNum: "9399142970086005",
				Sum:      100,
			}).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("get withdrawals content exists", func(t *testing.T) {
		e.GET("/api/user/withdrawals").
			Expect().
			Status(http.StatusOK)
	})

	t.Run("get balance not zero", func(t *testing.T) {
		e.GET("/api/user/balance").
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			ContainsKey("current").HasValue("current", 400).
			ContainsKey("withdrawn").HasValue("withdrawn", 100)
	})

	childCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_ = dbRes.QueryRowContext(childCtx, "DELETE FROM users WHERE login = $1", login)
}

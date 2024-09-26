package client

import (
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/services"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type OrderClient interface {
	CheckOrder(OrderNum string) (services.OrderFromAccrual, int, error)
}

type OrderClientService struct {
	logger               *zap.Logger
	accrualSystemAddress string
}

func NewOrderClient(opts ...func(*OrderClientService)) *OrderClientService {
	orderClient := &OrderClientService{}
	for _, opt := range opts {
		opt(orderClient)
	}
	return orderClient
}

func WithLogger(logger *zap.Logger) func(*OrderClientService) {
	return func(o *OrderClientService) {
		o.logger = logger
	}
}

func WithAccrualSystemAddress(accrualSystemAddress string) func(*OrderClientService) {
	return func(o *OrderClientService) {
		o.accrualSystemAddress = accrualSystemAddress
	}
}

func (c OrderClientService) CheckOrder(OrderNum string) (services.OrderFromAccrual, int, error) {
	data := services.OrderFromAccrual{}
	resp, err := http.Get(c.accrualSystemAddress + "/api/orders/" + OrderNum)
	if err != nil {
		return data, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		seconds, err := strconv.Atoi(retryAfter)
		if err != nil {
			seconds = 0
		}
		return data, seconds, errors.New("Error " + resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return data, 0, errors.New("Error " + resp.Status)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		return data, 0, err
	}
	return data, 0, nil
}

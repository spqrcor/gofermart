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

type OrderClient struct {
	logger               *zap.Logger
	accrualSystemAddress string
}

func NewOrderClient(logger *zap.Logger, accrualSystemAddress string) *OrderClient {
	return &OrderClient{logger: logger, accrualSystemAddress: accrualSystemAddress}
}

func (c OrderClient) CheckOrder(OrderNum string) (services.OrderFromAccrual, int, error) {
	data := services.OrderFromAccrual{}
	resp, err := http.Get(c.accrualSystemAddress + "/api/orders/" + OrderNum)
	if err != nil {
		return data, 0, err
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		seconds, err := strconv.Atoi(retryAfter)
		if err != nil {
			seconds = 0
		}
		return data, seconds, errors.New("Error " + resp.Status)
	} else if resp.StatusCode != http.StatusOK {
		return data, 0, errors.New("Error " + resp.Status)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		return data, 0, err
	}
	return data, 0, nil
}

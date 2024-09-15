package client

import (
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/services"
	"io"
	"net/http"
)

func CheckOrder(OrderNum string) (services.OrderFromAccrual, error) {
	data := services.OrderFromAccrual{}
	resp, err := http.Get(config.Cfg.AccrualSystemAddress + "/api/orders/" + OrderNum)
	if err != nil {
		return data, nil
	}
	if resp.StatusCode != http.StatusOK {
		return data, errors.New("Error " + resp.Status)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.Unmarshal(bodyBytes, &data); err != nil {
		return data, err
	}
	return data, nil
}

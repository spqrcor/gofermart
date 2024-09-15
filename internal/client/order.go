package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/config"
	"github.com/spqrcor/gofermart/internal/services"
	"io"
	"net/http"
)

func SendOrder(OrderNum string) error {
	url := config.Cfg.AccrualSystemAddress + "/api/orders"
	resp, err := http.Post(url, "application/json", bytes.NewReader([]byte(`{"order":`+OrderNum+`,"goods":[{"description":"Чайник Bork","price":7000}]}`)))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("Error " + resp.Status)
	}
	return nil
}

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

func SetReward() error {
	url := config.Cfg.AccrualSystemAddress + "/api/goods"
	resp, err := http.Post(url, "application/json", bytes.NewReader([]byte(`{"match":"Bork","reward":10,"reward_type":"%"}`)))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Error " + resp.Status)
	}
	return nil
}

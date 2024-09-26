package handlers

import (
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/services"
	"net/http"
)

func GetOrdersHandler(o services.Order) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		orders, err := o.GetAll(req.Context())
		if errors.Is(err, services.ErrOrdersNotFound) {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(res)
		if err := enc.Encode(orders); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

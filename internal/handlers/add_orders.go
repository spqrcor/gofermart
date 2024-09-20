package handlers

import (
	"errors"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/utils"
	"net/http"
)

func AddOrdersHandler(o *services.OrderService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		orderNum, err := utils.FromPostPlain(req)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = o.Add(req.Context(), orderNum)
		if errors.Is(err, services.ErrOrderInvalidFormat) {
			http.Error(res, err.Error(), http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, services.ErrOrderAnotherUserExists) {
			res.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, services.ErrOrderUserExists) {
			res.WriteHeader(http.StatusOK)
			return
		}
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusAccepted)
	}
}

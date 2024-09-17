package handlers

import (
	"errors"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/utils"
	"net/http"
)

func AddWithdrawalHandler(w *services.WithdrawalService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input services.InputWithdrawal
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := w.Add(req.Context(), input)
		if errors.Is(err, services.ErrOrderInvalidFormat) {
			http.Error(res, err.Error(), http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, services.ErrBalance) {
			http.Error(res, err.Error(), http.StatusPaymentRequired)
			return
		} else if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

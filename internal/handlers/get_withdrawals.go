package handlers

import (
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/services"
	"net/http"
)

func GetWithdrawalsHandler(w *services.WithdrawalService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		withdrawals, err := w.GetAll(req.Context())
		if errors.Is(err, services.ErrWithdrawNotFound) {
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
		if err := enc.Encode(withdrawals); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

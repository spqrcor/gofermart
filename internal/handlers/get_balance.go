package handlers

import (
	"encoding/json"
	"github.com/spqrcor/gofermart/internal/services"
	"net/http"
)

func GetBalanceHandler(w services.WithdrawalRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		balance, err := w.GetBalance(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(res)
		if err := enc.Encode(balance); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

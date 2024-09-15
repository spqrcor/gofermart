package handlers

import (
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/utils"
	"net/http"
)

func RegisterHandler(u services.UserRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input services.InputDataUser
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		UserID, err := u.Add(req.Context(), input)
		if errors.Is(err, services.ErrValidation) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		} else if errors.Is(err, services.ErrLoginExists) {
			res.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		authenticate.SetCookie(res, UserID)
		res.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(u services.UserRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var input services.InputDataUser
		if err := utils.FromPostJSON(req, &input); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		UserID, err := u.Login(req.Context(), input)
		if errors.Is(err, services.ErrLogin) {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		authenticate.SetCookie(res, UserID)
		res.WriteHeader(http.StatusOK)
	}
}

func AddOrdersHandler(o services.OrderRepository) http.HandlerFunc {
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
		} else if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusAccepted)
	}
}

func GetOrdersHandler(o services.OrderRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		orders, err := o.GetAll(req.Context())
		if errors.Is(err, services.ErrOrdersNotFound) {
			res.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
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

func GetWithdrawalsHandler(w services.WithdrawalRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		withdrawals, err := w.GetAll(req.Context())
		if errors.Is(err, services.ErrWithdrawNotFound) {
			res.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
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

func AddWithdrawalHandler(w services.WithdrawalRepository) http.HandlerFunc {
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

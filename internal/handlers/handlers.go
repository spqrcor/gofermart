package handlers

import (
	"encoding/json"
	"errors"
	"github.com/spqrcor/gofermart/internal/actions"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/utils"
	"net/http"
)

func RegisterHandler(res http.ResponseWriter, req *http.Request) {
	var input actions.InputDataUser
	if err := utils.FromPostJSON(req, &input); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, err := actions.Register(req.Context(), input)
	if errors.Is(err, actions.ErrValidation) {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	} else if errors.Is(err, actions.ErrLoginExists) {
		res.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	authenticate.SetCookie(res, UserID)
	res.WriteHeader(http.StatusOK)
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	var input actions.InputDataUser
	if err := utils.FromPostJSON(req, &input); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, err := actions.Login(req.Context(), input)
	if errors.Is(err, actions.ErrLogin) {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	authenticate.SetCookie(res, UserID)
	res.WriteHeader(http.StatusOK)
}

func AddOrdersHandler(res http.ResponseWriter, req *http.Request) {
	orderNum, err := utils.FromPostPlain(req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = actions.AddOrder(req.Context(), orderNum)
	if errors.Is(err, actions.ErrOrderInvalidFormat) {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	} else if errors.Is(err, actions.ErrOrderAnotherUserExists) {
		res.WriteHeader(http.StatusConflict)
		return
	} else if errors.Is(err, actions.ErrOrderUserExists) {
		res.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func GetOrdersHandler(res http.ResponseWriter, req *http.Request) {
	orders, err := actions.GetOrders(req.Context())
	if errors.Is(err, actions.ErrOrdersNotFound) {
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

func GetBalanceHandler(res http.ResponseWriter, req *http.Request) {
	balanceInfo, err := actions.GetBalanceInfo(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	if err := enc.Encode(balanceInfo); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetWithdrawalsHandler(res http.ResponseWriter, req *http.Request) {
	withdraws, err := actions.GetWithdraws(req.Context())
	if errors.Is(err, actions.ErrWithdrawNotFound) {
		res.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	if err := enc.Encode(withdraws); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func BalanceWithdrawHandler(res http.ResponseWriter, req *http.Request) {
	var input actions.InputWithdraw
	if err := utils.FromPostJSON(req, &input); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := actions.AddWithdraw(req.Context(), input)
	if errors.Is(err, actions.ErrOrderInvalidFormat) {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	} else if errors.Is(err, actions.ErrBalance) {
		http.Error(res, err.Error(), http.StatusPaymentRequired)
		return
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

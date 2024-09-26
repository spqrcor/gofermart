package handlers

import (
	"errors"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/utils"
	"net/http"
)

func RegisterHandler(u services.User, a authenticate.Auth) http.HandlerFunc {
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
		}
		if errors.Is(err, services.ErrLoginExists) {
			res.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.SetCookie(res, UserID)
		res.WriteHeader(http.StatusOK)
	}
}

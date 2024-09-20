package handlers

import (
	"errors"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/utils"
	"net/http"
)

func LoginHandler(u *services.UserService, a *authenticate.Authenticate) http.HandlerFunc {
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
		}
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.SetCookie(res, UserID)
		res.WriteHeader(http.StatusOK)
	}
}

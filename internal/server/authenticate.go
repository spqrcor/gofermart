package server

import (
	"context"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"go.uber.org/zap"
	"net/http"
)

func authenticateMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				http.Error(rw, err.Error(), http.StatusUnauthorized)
				return
			} else {
				UserID, err := authenticate.GetUserIDFromCookie(cookie.Value)
				if err != nil {
					logger.Error(err.Error())
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
				ctx := context.WithValue(r.Context(), authenticate.ContextUserID, UserID)
				next.ServeHTTP(rw, r.WithContext(ctx))
			}
		})
	}
}

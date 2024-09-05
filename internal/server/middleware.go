package server

import (
	"compress/gzip"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"slices"
	"time"
)

func getBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Log.Error(err.Error())
			} else {
				r.Body = gz
			}
			if err = gz.Close(); err != nil {
				logger.Log.Error(err.Error())
			}
		}
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(start)

		logger.Log.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", ww.Status()),
			zap.Int("content-length", ww.BytesWritten()),
			zap.String("duration", duration.String()),
		)
	})
}

func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			idx := slices.IndexFunc(publicRoutes, func(c string) bool { return c == r.RequestURI })
			if idx == -1 {
				http.Error(rw, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(rw, r)
		} else {
			UserID, err := authenticate.GetUserIDFromCookie(cookie.Value)
			if err != nil {
				logger.Log.Error(err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), authenticate.ContextUserID, UserID)
			next.ServeHTTP(rw, r.WithContext(ctx))
		}
	})
}

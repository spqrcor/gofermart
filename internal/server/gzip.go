package server

import (
	"compress/gzip"
	"go.uber.org/zap"
	"net/http"
)

func getBodyMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.Header.Get(`Content-Encoding`) != `gzip` {
				next.ServeHTTP(rw, r)
				return
			}
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Error(err.Error())
			} else {
				r.Body = gz
			}
			if err = gz.Close(); err != nil {
				logger.Error(err.Error())
			}
			next.ServeHTTP(rw, r)
		})
	}
}

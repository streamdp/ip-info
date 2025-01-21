package rest

import (
	"log"
	"net/http"

	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/server"
)

func rateLimiterMW(limiter server.Limiter, l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := limiter.Limit(httpClientIp(r)); err != nil {
			if err = writeJsonResponse(w, getHttpStatus(err), domain.NewResponse(err, nil)); err != nil {
				l.Println(err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

package rest

import (
	"log"
	"net/http"

	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/ratelimiter"
)

func rateLimiterMW(limiter ratelimiter.Limiter, l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := limiter.Limit(httpClientIp(r)); err != nil {
			w.WriteHeader(getHttpStatus(err))
			if _, err = w.Write(domain.NewResponse(err, nil).Bytes()); err != nil {
				l.Println(err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

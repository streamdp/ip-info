package rest

import (
	"errors"
	"fmt"
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

var errWrongContentType = errors.New("content type not implemented")

func contentTypeRestrictionMW(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c := r.Header.Get(contentTypeHeader); c != "" && c != jsonContentType {
			err := fmt.Errorf("%w: %s", errWrongContentType, c)
			if err = writeJsonResponse(w, getHttpStatus(err), domain.NewResponse(err, nil)); err != nil {
				l.Println(err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})

}

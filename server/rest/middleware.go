package rest

import (
	"errors"
	"log"
	"net/http"
	"slices"

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

func contentTypeRestrictionMW(l *log.Logger, f http.HandlerFunc, allowedTypes ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isAllowedContentType(r.Header.Get(contentTypeHeader), allowedTypes) {
			f.ServeHTTP(w, r)

			return
		}

		if err := writeJsonResponse(
			w, getHttpStatus(errWrongContentType), domain.NewResponse(errWrongContentType, nil),
		); err != nil {
			l.Println(err)
		}
	}
}

func isAllowedContentType(c string, allowedTypes []string) bool {
	if c == "" {
		return true
	}

	return slices.Contains(allowedTypes, c)
}
